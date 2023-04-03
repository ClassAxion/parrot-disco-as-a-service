package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/config"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/router"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/authservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/dashboardservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/deployservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/userservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/worker"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mlvzk/gopgs/pgkv"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

func startServer(ctx context.Context, services *service.Services, port int) {
	r := gin.Default()

	r.HTMLRender = ginview.New(goview.Config{
		Root:         "views",
		Extension:    ".tpl",
		Master:       "layouts/master",
		Funcs:        template.FuncMap{},
		DisableCache: true,
	})

	store := cookie.NewStore([]byte("video-creator-manager"))

	store.Options(sessions.Options{
		MaxAge: 2678400,
		Path:   "/",
	})

	r.Use(sessions.Sessions("session", store))

	handler := router.New(r, services)

	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	server.ListenAndServe()
}

func startWorkers(ctx context.Context, db *database.Database, kv *pgkv.Store, services *service.Services, config *config.Config) {
	workerContext := worker.WorkerContext{
		Config:   config,
		DB:       db,
		Services: services,
		KeyValue: kv,
	}

	workers := []worker.Worker{}

	var wg sync.WaitGroup
	for i, w := range workers {
		wg.Add(1)

		go func(i int, w worker.Worker) {
			defer wg.Done()

			if err := w(ctx, workerContext); err != nil {
				if errors.Is(err, ctx.Err()) {
					return
				}

				panic(fmt.Errorf("failed while running worker(index %d): %w", i, err))
			}
		}(i, w)
	}

	wg.Wait()
}

func main() {
	cfg := config.LoadConfig()

	// if err := database.RunMigrations(context.Background(), cfg.DatabaseUrl); err != nil {
	// 	log.Fatalln("Failed to run migrations", err)
	// }

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db := database.NewPg(cfg.DatabaseUrl)

	dashboardService := dashboardservice.New()
	authService := authservice.New(db)
	userService := userservice.New(db)
	deployService := deployservice.New(db)

	oauthcfg := &oauth2.Config{}

	ts := oauthcfg.TokenSource(ctx, &oauth2.Token{AccessToken: cfg.VultrApiKey})
	vultrClient := govultr.NewClient(oauth2.NewClient(ctx, ts))

	services := service.Services{
		DashboardService: dashboardService,
		AuthService:      authService,
		UserService:      userService,
		DeployService:    deployService,
		Vultr:            vultrClient,
	}

	kv, err := pgkv.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("failed to connect to kv: %w", err))
	}

	defer kv.Close()

	go func() {
		<-ctx.Done()
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}()

	go startWorkers(ctx, db, kv, &services, cfg)

	startServer(ctx, &services, cfg.Port)
}
