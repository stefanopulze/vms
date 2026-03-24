package api

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vms-core/internal/api/handler"
	"vms-core/internal/config"
	"vms-core/internal/serial"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func BindApi(cfg config.ServerConfig, router *chi.Mux, port serial.Serial, ih *handler.InverterHandler, sh *handler.StatsHandler) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	if len(cfg.CorsAllowedOrigins) > 0 {
		slog.Debug("Enabling CORS")
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   cfg.CorsAllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}

	router.Route("/api", func(r chi.Router) {
		bindInverterApi(r, ih)
		bindStatsApi(r, sh)
		bindHealthApi(r, port)
	})

	bindStaticFiles(router)
}

func bindStaticFiles(router *chi.Mux) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})
}

func bindStatsApi(router chi.Router, handler *handler.StatsHandler) {
	router.Get("/stats", handler.GetDayStats)
	router.Get("/downsampling", handler.DownsamplingDayStats)
}

func bindHealthApi(router chi.Router, port serial.Serial) {
	sh := handler.NewStatus(port)
	router.Get("/health", sh.Health)
	router.Get("/status", sh.Status)
}

func bindInverterApi(router chi.Router, ih *handler.InverterHandler) {
	router.Route("/inverter", func(r chi.Router) {
		r.Get("/info", ih.AggregateInfo)
		r.Get("/flags", ih.QueryFlags)
		r.Put("/flags", ih.UpdateFlags)
		r.Get("/rating-info", ih.QueryRatingInfo)
		r.Get("/general-status", ih.QueryStatus)
		r.Get("/mode", ih.QueryMode)
		r.Get("/warnings", ih.QueryWarnings)
		r.Get("/time", ih.QueryTime)
		r.Post("/time", ih.UpdateTime)
		r.Get("/query", ih.QueryCommand)
		r.Put("/source-priority", ih.UpdateSourcePriority)
		r.Put("/charger-source-priority", ih.UpdateChargerSourcePriority)

		r.Get("/max-ac-charging-current-values", ih.QueryMaxAcChargingCurrentValues)
		r.Put("/max-ac-charging-current", ih.UpdateMaxAcChargingCurrent)
	})
}
