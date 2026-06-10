package router

import (
	"context"
	"fmt"
	"net/http"

	stderrors "errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/graph/generated"
	resolver "github.com/Toflex/directory_v2/graph/resolvers"
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/middlewares"
	"github.com/Toflex/directory_v2/pkg/configuration"
	apperrors "github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/samber/do/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type asynqMonitorConfig struct {
	RedisURL      string `env:"REDIS_URL"`
	RedisDB       int    `env:"REDIS_DB"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func InitializeRoutes(engine *gin.Engine) {

	engine.GET("/health", middlewares.Logger(), func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "Server is up 🚀",
		})
	})

	resolverStruct := resolver.Resolver{
		AuthenticationService: do.MustInvoke[authentication.IService](runtime.Injector),
	}

	gqlHandler := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolverStruct}))
	playgroundHandler := playground.Handler("API GraphQL playground", "/api")

	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})

	gqlHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	gqlHandler.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Error("GraphQL panic recovered: %v", err)
		return gqlerror.Errorf("internal server error")
	})

	gqlHandler.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		var gqlErr *gqlerror.Error
		if stderrors.As(err, &gqlErr) {
			return gqlErr
		}

		var safe *apperrors.SafeError
		if stderrors.As(err, &safe) {
			return &gqlerror.Error{
				Message: safe.Message,
				Extensions: map[string]any{
					"code": string(safe.Code),
				},
			}
		}
		return graphql.DefaultErrorPresenter(ctx, err)
	})

	// GraphQL endpoint
	engine.POST("/api", middlewares.Logger(), func(c *gin.Context) {
		ctx := log.SetLoggerInContext(c.Request.Context())
		gqlHandler.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	})

	basicAuth := gin.BasicAuth(gin.Accounts{
		configuration.Instance.BasicUsername: configuration.Instance.BasicPassword,
	})

	// GraphQL Playground
	engine.GET("/api", basicAuth, func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})

	monitorConf := &asynqMonitorConfig{}
	configuration.Load(monitorConf)

	monitorHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/monitoring",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     monitorConf.RedisURL,
			Password: monitorConf.RedisPassword,
			DB:       monitorConf.RedisDB,
		},
	})

	engine.Any(monitorHandler.RootPath(), basicAuth, gin.WrapH(monitorHandler))
	engine.Any(monitorHandler.RootPath()+"/*path", basicAuth, gin.WrapH(monitorHandler))

}
