package app

import (
	"context"
	"fmt"
	"gobanking/internal/bank/config"
	http_handler "gobanking/internal/bank/handler/http"
	"gobanking/internal/bank/repository"
	"gobanking/internal/bank/usecase"
	"gobanking/pkg/logger"
	"gobanking/pkg/middlewares"
	natsPkg "gobanking/pkg/nats"
	"gobanking/pkg/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const (
	waitShotDownDuration = 3 * time.Second
)

type app struct {
	log               logger.Logger
	cfg               *config.Config
	doneCh            chan struct{}
	middlewareManager middlewares.MiddlewareManager
	validator         *utils.CustomValidator
	mysqlConn         *gorm.DB
	echo              *echo.Echo
	natsClient        *nats.Conn
}

func NewAppBank(log logger.Logger, cfg *config.Config) *app {
	return &app{
		log:       log,
		cfg:       cfg,
		doneCh:    make(chan struct{}),
		validator: utils.NewFieldError(validator.New()),
		echo:      echo.New(),
	}
}

func (a *app) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// connect mysql
	if err := a.connectMysql(ctx); err != nil {
		return err
	}
	sqlDB, _ := a.mysqlConn.DB()
	defer sqlDB.Close()

	// migrate mysql
	if err := a.migrateMysql(ctx); err != nil {
		return err
	}

	// seeding currencies
	if err := a.seedCurrencies(ctx); err != nil {
		return err
	}

	// setup nats
	natsClient, err := natsPkg.NewNatsConnect(a.cfg.Nats, a.log)
	if err != nil {
		return err
	}
	a.natsClient = natsClient
	defer a.natsClient.Drain()
	defer a.natsClient.Close()

	a.middlewareManager = middlewares.NewMiddlewareManager(a.log,
		&middlewares.MiddlewareConfig{
			HttpClientDebug:     a.cfg.Http.HttpClientDebug,
			DebugErrorsResponse: a.cfg.Http.DebugErrorsResponse,
		},
		nil,
		a.natsClient,
	)

	// setup app
	appRepo := repository.NewRepository(a.mysqlConn)
	appUsecase := usecase.NewUsecase(appRepo, a.log)
	http_handler.NewHandler(appUsecase, a.log, a.cfg, a.middlewareManager, a.validator, a.echo.Group(""))

	go func() {
		if err := a.runHttpServer(); err != nil {
			a.log.Errorf("(runHttpServer) err: %v", err)
			cancel()
		}
	}()
	a.log.Infof("%s is listening on PORT: %v", GetMicroserviceName(a.cfg), a.cfg.Http.Port)

	<-ctx.Done()
	a.waitShutDown(waitShotDownDuration)

	if err := a.echo.Shutdown(ctx); err != nil {
		a.log.Warnf("(Shutdown) err: %v", err)
	}

	<-a.doneCh
	a.log.Infof("%s app exited properly", GetMicroserviceName(a.cfg))
	return nil
}

func (a *app) waitShutDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		a.doneCh <- struct{}{}
	}()
}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
