package main

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-query-svc/internal/config"
	grpcHandler "be-realtime-chat-app/services/chat-query-svc/internal/delivery/grpc/handler"
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/controller"
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/route"
	"be-realtime-chat-app/services/chat-query-svc/internal/repository"
	"be-realtime-chat-app/services/chat-query-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/discovery/consul"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"fmt"
	"net"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServer *grpc.Server
	app        *fiber.App
)

func webServer(ctx context.Context) error {
	serverConfig := config.NewServerConfig()
	app = config.NewApp()

	logger, _ := logs.NewLogger()

	cqlClient, err := config.NewCQL()
	elasticsearchClient, err := config.NewElasticsearch()

	customValidator := helper.NewCustomValidator()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.QuerySvcName)
	if err != nil {
		logger.Error("Failed to create consul registry for service" + err.Error())
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.QuerySvcName + "-grpc")
	grpcPortInt, _ := strconv.Atoi(serverConfig.QueryGRPCPort)

	err = registry.RegisterService(ctx, serverConfig.QuerySvcName+"-grpc", GRPCserviceID, serverConfig.QueryGRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logger.Error("Failed to register realtime service to consul", zap.Error(err))
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create user adapter", zap.Error(err))
	}

	roomAdapter, err := adapter.NewRoomAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create room adapter", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		logger.Info("Context canceled. Deregistering services...")
		registry.DeregisterService(context.Background(), GRPCserviceID)

		logger.Info("Shutting down servers...")
		if err := app.Shutdown(); err != nil {
		}
		if grpcServer != nil {
			grpcServer.GracefulStop()
		}
		logger.Info("Successfully shutdown...")
	}()

	go consul.StartHealthCheckLoop(ctx, registry, GRPCserviceID, serverConfig.QuerySvcName+"-grpc", logger)

	messageCQLRepo := repository.NewMessageCQLRepository(cqlClient)
	messageElasticRepo := repository.NewMessageElasticRepo(elasticsearchClient)

	queryUC := usecase.NewQueryUseCase(messageCQLRepo, messageElasticRepo, roomAdapter, customValidator, logger)

	queryController := controller.NewQueryController(queryUC, logger)

	userMiddleware := middleware.NewUserAuth(userAdapter, logger)

	roomRoute := route.NewRoomRoute(app, queryController, userMiddleware)
	roomRoute.RegisterRoutes()

	go func() {
		grpcServer = grpc.NewServer()
		reflection.Register(grpcServer)
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverConfig.QueryGRPCAddr, serverConfig.QueryGRPCPort))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to listen: %v", err))
			return
		}

		defer l.Close()

		grpcHandler.NewQueryGRPCHandler(grpcServer, queryUC)

		if err := grpcServer.Serve(l); err != nil {
			logger.Error(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Listen(fmt.Sprintf("%s:%s", serverConfig.QueryHTTPAddr, serverConfig.QueryHTTPPort))
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-serverErrors:
		return err
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := webServer(ctx); err != nil {
		fmt.Printf("Error starting web server: %v\n", err)
		return
	}
}
