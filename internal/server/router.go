package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"github.com/openmfp/golang-commons/logger"
	"github.com/rs/cors"

	"github.com/openmfp/extension-content-operator/internal/config"
	"github.com/openmfp/extension-content-operator/pkg/validation"
)

func CreateRouter(
	appConfig config.Config,
	log *logger.Logger,
	validator validation.ExtensionConfiguration,
) *chi.Mux {
	router := chi.NewRouter()

	if appConfig.IsLocal {
		rl := requestLogger{
			log: log,
		}

		router.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{headers.Accept, headers.Authorization, headers.ContentType, headers.XCSRFToken},
			Debug:            false,
		}).Handler)
		router.Use(rl.Handler)
	}

	vh := validationHandler{
		validator: validator,
		log:       log,
	}

	router.MethodFunc(http.MethodPost, "/validate", vh.Handler)

	return router
}
