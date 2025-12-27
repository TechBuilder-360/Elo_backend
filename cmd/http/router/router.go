package router

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/graph/generated"
	resolver "github.com/Toflex/directory_v2/graph/resolvers"
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func InitializeRoutes(engine *gin.Engine) {
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
		// Notify bug tracker...
		log.Error("An error occurred: %s", err)
		return gqlerror.Errorf("Internal server error!")
	})

	// GraphQL endpoint
	engine.POST("/api", func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	basicAuth := gin.BasicAuth(gin.Accounts{
		configuration.Instance.BasicUsername: configuration.Instance.BasicPassword,
	})

	// GraphQL Playground
	engine.GET("/api", basicAuth, func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})

}
