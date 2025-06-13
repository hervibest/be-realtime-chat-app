package main

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/config"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/delivery/http/controller"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/delivery/http/middleware"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/delivery/http/route"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/discovery/consul"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"fmt"
	"os/signal"
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
	db := config.NewPostgresDatabase()
	defer db.Close()

	natsConn := config.NewNATsConn(logger)
	messagingAdapter := adapter.NewMessagingAdapter(natsConn)

	customValidator := helper.NewCustomValidator()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.ChatSvcName)
	if err != nil {
		logger.Error("Failed to create consul registry for service" + err.Error())
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create user adapter", zap.Error(err))
	}

	roomAdapter, err := adapter.NewRoomAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create room adapter", zap.Error(err))
	}

	queryAdapter, err := adapter.NewQueryAdapter(ctx, registry, logger)
	if err != nil {
		logger.Error("Failed to create room adapter", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		logger.Info("Context canceled. Deregistering services...")

		logger.Info("Shutting down servers...")
		if err := app.Shutdown(); err != nil {
		}
		if grpcServer != nil {
			grpcServer.GracefulStop()
		}
		logger.Info("Successfully shutdown...")
	}()

	chatUC := usecase.NewChatUseCase(messagingAdapter, roomAdapter, queryAdapter, customValidator, logger)

	chatController := controller.NewChatController(chatUC, logger)

	userMiddleware := middleware.NewUserAuth(userAdapter, logger)

	chatRoute := route.NewRoomRoute(app, chatController, userMiddleware)
	chatRoute.RegisterRoutes()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Listen(fmt.Sprintf("%s:%s", serverConfig.ChatHTTPAddr, serverConfig.ChatHTTPPort))
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
