package main

import (
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/discovery/consul"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/user-svc/internal/adapter"
	"be-realtime-chat-app/services/user-svc/internal/config"
	grpcHandler "be-realtime-chat-app/services/user-svc/internal/delivery/grpc/handler"
	"be-realtime-chat-app/services/user-svc/internal/delivery/http/controller"
	"be-realtime-chat-app/services/user-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/user-svc/internal/delivery/http/route"
	"be-realtime-chat-app/services/user-svc/internal/repository"
	"be-realtime-chat-app/services/user-svc/internal/usecase"
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
	db := config.NewPostgresDatabase()
	defer db.Close()

	redisClient := config.NewRedisClient()
	jwtAdapter := adapter.NewJWTAdapter()
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
		logger.Error("Failed to register user service to consul", zap.Error(err))
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

	userRepo := repository.NewUserRepository(logger)

	userUC := usecase.NewUserUseCase(db, userRepo, jwtAdapter, cacheAdapter, customValidator, logger)

	userController := controller.NewUserController(userUC, logger)

	go func() {
		grpcServer = grpc.NewServer()
		logger.Info("Initiate grpc server stage 1")
		reflection.Register(grpcServer)
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverConfig.UserGRPCAddr, serverConfig.UserGRPCPort))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to listen: %v", err))
			return
		}
		logger.Info("Closing")

		defer l.Close()

		logger.Info("Initiate grpc server stage 2")
		grpcHandler.NewUserGRPCHandler(grpcServer, userUC)

		if err := grpcServer.Serve(l); err != nil {
			logger.Error(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	userMiddleware := middleware.NewUserAuth(userUC, customValidator, logger)

	userRoute := route.NewUserRoute(app, userController, userMiddleware)
	userRoute.RegisterRoutes()

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
