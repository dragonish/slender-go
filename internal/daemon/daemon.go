package daemon

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"slender/internal/apis"
	"slender/internal/config"
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/icons"
	"slender/internal/logger"
	"slender/internal/model"
	"slender/internal/pages"
)

// StartDaemon starts the daemon.
func StartDaemon() {
	if !global.Flags.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGALRM)
	defer stop()

	dirLog := logger.New("dir", model.DATA_DIR)
	if !data.IsPathExists(model.DATA_DIR + "/") {
		dirLog.Info("create data directory")
		err := os.MkdirAll(model.DATA_DIR, 0777)
		if err != nil {
			dirLog.Fatal("error creating data directory", err)
		}
	} else {
		dirLog.Debug("data directory exists")
	}

	var (
		dbFilename     string
		configFilename string
	)
	if global.Flags.DebugMode {
		dbFilename = "slender_debug"
		configFilename = "config_debug"
	} else {
		dbFilename = "slender"
		configFilename = "config"
	}

	database.Load(dbFilename)
	config.Load(configFilename)

	icons.Build()

	router := gin.New()
	if global.Flags.DebugMode {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())

	if !global.Flags.DebugMode {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	router.Use(getAllowCors())

	//* Register routing
	apis.Apis(router)
	pages.Pages(router)

	srv := &http.Server{
		Addr:              ":" + global.Flags.GetPortStr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err("application startup error", err)
		}
	}()
	logger.Info("application started")

	<-ctx.Done()

	stop()
	logger.Info("application closing", "note", "Press Ctrl+C to exit immediately")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Err("application is forced to close", err)
	}

	logger.Info("application has been closed")
}

// getAllowCors returns allow all origins cors handler.
func getAllowCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
	})
}
