package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/dealense7/go-rate-app/internal/handlers"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/repositories"
	"github.com/dealense7/go-rate-app/internal/services"
	"github.com/dealense7/go-rate-app/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			// so we get a *zap.Logger
			zap.NewDevelopment,

			utils.NewDB,
			NewGinEngine,
			handlers.NewWebHandler,

			fx.Annotate(
				services.NewGasService,
				fx.As(new(interfaces.GasService)),
			),

			fx.Annotate(
				services.NewStoreService,
				fx.As(new(interfaces.StoreService)),
			),

			fx.Annotate(
				repositories.NewMySQLGasRepository,
				fx.As(new(interfaces.GasRepository)),
			),

			fx.Annotate(
				repositories.NewMySQLStoreRepository,
				fx.As(new(interfaces.StoreRepository)),
			),
		),
		fx.Invoke(
			RegisterRoutes,
			NewHTTPServer, // ← now Fx will call this, registering your OnStart hook
		),
	).Run()
}

// NewGinEngine constructs the Gin router with middleware, templates & static files.
func NewGinEngine() *gin.Engine {
	router := gin.Default()

	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))
	router.SetHTMLTemplate(tmpl)

	router.Static("/static", "./static")
	return router

}

// NewHTTPServer wires up an http.Server and registers its lifecycle hooks.
func NewHTTPServer(lc fx.Lifecycle, router *gin.Engine, logger *zap.Logger) {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting Gin server", zap.String("addr", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down Gin server")
			return srv.Shutdown(ctx)
		},
	})
}

// RegisterRoutes binds your HTTP endpoints to the WebHandler.
func RegisterRoutes(router *gin.Engine, h *handlers.WebHandler) {
	router.GET("/", h.GetProducts)
	router.GET("/items", h.GetProductList)
	router.GET("/prices/:id", h.GetProductPrices)

}
