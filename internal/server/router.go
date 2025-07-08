package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"github.com/platform-mesh/golang-commons/logger"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/openmfp/extension-manager-operator/internal/config"
	"github.com/openmfp/extension-manager-operator/pkg/graph"
	"github.com/openmfp/extension-manager-operator/pkg/resolver"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
)

func CreateRouter(
	mgr mcmanager.Manager,
	appConfig config.ServerConfig,
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
			AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		}).Handler)
		router.Use(rl.Handler)
	}

	vh := NewHttpValidateHandler(log, validator)

	hd := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver.NewResolver(mgr, appConfig.ProviderWorkspaceID),
	}))

	hd.AddTransport(transport.POST{})
	hd.AddTransport(transport.GET{})
	hd.AddTransport(transport.Options{})
	hd.Use(extension.Introspection{})

	if appConfig.IsLocal {
		router.Method(http.MethodGet, "/", playground.Handler("GraphQL playground", "/graphql"))
	}

	router.Method(http.MethodPost, "/graphql", hd)

	router.MethodFunc(http.MethodPost, "/validate", vh.HandlerValidate)
	router.MethodFunc(http.MethodGet, "/healthz", vh.HandlerHealthz)
	router.MethodFunc(http.MethodGet, "/readyz", vh.HandlerHealthz)

	return router
}
