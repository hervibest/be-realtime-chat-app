package main

import (
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/discovery/consul"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/room-svc/internal/adapter"
	"be-realtime-chat-app/services/room-svc/internal/config"
	grpcHandler "be-realtime-chat-app/services/room-svc/internal/delivery/grpc/handler"
	"be-realtime-chat-app/services/room-svc/internal/delivery/http/controller"
	"be-realtime-chat-app/services/room-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/room-svc/internal/delivery/http/route"
	"be-realtime-chat-app/services/room-svc/internal/repository"
	"be-realtime-chat-app/services/room-svc/internal/usecase"
	"context"
	"fmt"
	"net"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	db := config.NewPostgresDatabase()
	defer db.Close()

	redisClient := config.NewRedisClient()
	natsConn := config.NewNATsConn(logger)
	cacheAdapter := adapter.NewCacheAdapter(redisClient)
	messaginAdapter := adapter.NewMessagingAdapter(natsConn)

	customValidator := helper.NewCustomValidator()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.RoomSvcName)
	if err != nil {
		logger.Error("Failed to create consul registry for service" + err.Error())
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.RoomSvcName + "-grpc")
	grpcPortInt, _ := strconv.Atoi(serverConfig.RoomGRPCPort)

	err = registry.RegisterService(ctx, serverConfig.RoomSvcName+"-grpc", GRPCserviceID, serverConfig.RoomGRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logger.Error("Failed to register chat service to consul", zap.Error(err))
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

	go consul.StartHealthCheckLoop(ctx, registry, GRPCserviceID, serverConfig.RoomSvcName+"-grpc", logger)

	roomRepo := repository.NewRoomRepository(logger)

	roomUC := usecase.NewRoomUseCase(db, roomRepo, messaginAdapter, cacheAdapter, customValidator, logger)

	roomController := controller.NewRoomController(roomUC, logger)

	go func() {
		grpcServer = grpc.NewServer()
		logger.Info("Initiate grpc server stage 1")
		reflection.Register(grpcServer)
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverConfig.RoomGRPCAddr, serverConfig.RoomGRPCPort))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to listen: %v", err))
			return
		}
		logger.Info("Closing")

		defer l.Close()

		logger.Info("Initiate grpc server stage 2")
		grpcHandler.NewRoomGRPCHandler(grpcServer, roomUC)

		if err := grpcServer.Serve(l); err != nil {
			logger.Error(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	userMiddleware := middleware.NewUserAuth(userAdapter, logger)

	roomRoute := route.NewRoomRoute(app, roomController, userMiddleware)
	roomRoute.RegisterRoutes()

	logger.Info("Initiate server stage 2")

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Initiate server stage 2")

		serverErrors <- app.Listen(fmt.Sprintf("%s:%s", serverConfig.RoomHTTPAddr, serverConfig.RoomHTTPPort))
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
		log.Error("Failed to start web server", err)
	}
}
