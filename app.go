package main

import (
	"brigade-service/api"
	"brigade-service/cluster/user"
	"brigade-service/config"
	dbbrigade "brigade-service/database/brigade"
	"brigade-service/service/brigade"
	"context"
	"fmt"
	"io/fs"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	serviceName = "brigade-service"
	dbTimeout   = 15 * time.Second
)

type App struct {
	/* main */
	mainCtx  context.Context
	log      golog.Logger
	settings config.Settings

	/* http */
	server goserver.Server

	/* db */
	postgres     *sqlx.DB
	kafka        gokafka.Kafka
	taskConsumer gokafka.Consumer

	/* services */
	brigadeService *brigade.Service
}

func NewApp(mainCtx context.Context, log golog.Logger, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases(fs fs.FS, path string) (err error) {
	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}

	err = db.Migrate(fs, a.log, a.postgres, path, "postgres")
	if err != nil {
		return fmt.Errorf("migrate postgres: %w", err)
	}

	a.kafka = gokafka.NewKafka(a.settings.Databases.Kafka.Brokers)

	a.taskConsumer, err = a.kafka.Consumer(a.log.WithTags("taskConsumer"), func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.mainCtx)
	}, gokafka.WithTopic(a.settings.Databases.Kafka.Topics.Tasks), gokafka.WithConsumerGroup(serviceName))
	if err != nil {
		return fmt.Errorf("init task consumer: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	brigadeRepository := dbbrigade.NewRepository(a.postgres)

	httpClient := gohttp.NewClient(gohttp.WithTimeout(1 * time.Minute))

	userClient := user.NewClient(httpClient, a.settings.Cluster.UserService)

	a.brigadeService = brigade.NewService(brigadeRepository, userClient)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.mainCtx, a.log, a.settings)
	sb.AddDebug()
	sb.AddBrigades(a.brigadeService)

	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
	a.taskConsumer.Subscribe(a.brigadeService.SubscriberOnTaskEvent(a.mainCtx, a.log.WithTags("taskSubscriber")))
}

func (a *App) Stop(ctx context.Context) {
	consumerCtx, cancelConsumerCtx := context.WithTimeout(ctx, dbTimeout)
	defer cancelConsumerCtx()

	err := a.taskConsumer.Close(consumerCtx)
	if err != nil {
		a.log.Errorf("failed to close task consumer: %v", err)
	}

	a.server.Stop()

	err = a.postgres.Close()
	if err != nil {
		a.log.Errorf("failed to close postgres connection: %v", err)
	}
}
