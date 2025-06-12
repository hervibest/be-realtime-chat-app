package main

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-query-svc/internal/config"
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/chat-query-svc/internal/delivery/http/route"
	"be-realtime-chat-app/services/chat-query-svc/internal/repository"
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/discovery/consul"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	grpcServer *grpc.Server
	app        *fiber.App
)

func webServer(ctx context.Context) error {
	serverConfig := config.NewServerConfig()
	app = config.NewApp()

	logger, _ := logs.NewLogger()

	redisClient := config.NewRedisClient()
	cqlClient, err := config.NewCassandraSession()
	cacheAdapter := adapter.NewCacheAdapter(redisClient)

	customValidator := helper.NewCustomValidator()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.UserSvcName)
	if err != nil {
		logger.Error("Failed to create consul registry for service" + err.Error())
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.UserSvcName + "-grpc")
	grpcPortInt, _ := strconv.Atoi(serverConfig.UserGRPCPort)

	err = registry.RegisterService(ctx, serverConfig.UserSvcName+"-grpc", GRPCserviceID, serverConfig.UserGRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logger.Error("Failed to register realtime service to consul", zap.Error(err))
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create user adapter", zap.Error(err))
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

	go consul.StartHealthCheckLoop(ctx, registry, GRPCserviceID, serverConfig.UserSvcName+"-grpc", logger)

	roomRepo := repository.MessageCQLRepository(logger)

	roomUC := usecase.NewRoomUseCase(db, roomRepo, messaginAdapter, cacheAdapter, customValidator, logger)

	roomController := controller.NewRoomController(roomUC, logger)

	userMiddleware := middleware.NewUserAuth(userAdapter, logger)

	roomRoute := route.NewRoomRoute(app, roomController, userMiddleware)
	roomRoute.RegisterRoutes()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Listen(fmt.Sprintf("%s:%s", serverConfig.UserHTTPAddr, serverConfig.UserHTTPPort))
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
	}
}
