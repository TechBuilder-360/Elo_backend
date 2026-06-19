package router

import (
	"context"
	"net/http"

	stderrors "errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/graph/generated"
	resolver "github.com/Toflex/directory_v2/graph/resolvers"
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/middlewares"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/errors"
	apperrors "github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/providers/dojah"
	"github.com/Toflex/directory_v2/pkg/verification"
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

func authUserDirective(a authentication.IService) func(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	return func(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
		logger := log.LoggerInContext(ctx)
		opCtx := graphql.GetRequestContext(ctx)
		authHeader := opCtx.Headers.Get("Authorization")

		token, err := middlewares.ExtractBearerToken(authHeader)
		if err != nil {
			logger.WithError(err).Error("failed to extract token from Authorization header")
			return nil, gqlerror.Errorf("unauthorized")
		}

		userID, valid := a.VerifyJWT(ctx, token)
		if !valid || userID == "" {
			logger.Error("unable to validate jwt token")
			return nil, gqlerror.Errorf("unauthorized")
		}

		usr, err := database.DBInstance().User.Query().Where(user.IDEQ(userID)).First(ctx)
		if err != nil {
			logger.WithError(err).WithField("user_id", userID).Error("failed to fetch user")
			return nil, gqlerror.Errorf("unauthorized")
		}

		ctx = context.WithValue(ctx, middlewares.UserContextKey, usr)
		return next(ctx)
	}
}

func initalizeGQLRoute(engine *gin.Engine, basicAuth gin.HandlerFunc) {
	resolverStruct := resolver.Resolver{
		AuthenticationService: do.MustInvoke[authentication.IService](runtime.Injector),
		VerificationService:   verification.NewService(),
	}

	c := generated.Config{Resolvers: &resolverStruct}

	c.Directives.AuthUser = authUserDirective(resolverStruct.AuthenticationService)

	gqlHandler := handler.New(generated.NewExecutableSchema(c))
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
		log.WithError(err).Error("GraphQL panic recovered")
		return gqlerror.Errorf("internal server error")
	})

	gqlHandler.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		var gqlErr *gqlerror.Error
		if stderrors.As(err, &gqlErr) {
			log.WithError(gqlErr).Error("gql error presenter")
			return &gqlerror.Error{
				Message: gqlErr.Message,
				Extensions: map[string]any{
					"code": string(errors.ErrFailed),
				},
			}
		}

		var safe *apperrors.SafeError
		if stderrors.As(err, &safe) {
			log.WithError(safe).Error("gql error presenter")
			return &gqlerror.Error{
				Message: safe.Message,
				Extensions: map[string]any{
					"code":  string(safe.Code),
					"stack": safe.Stack,
				},
			}
		}

		log.WithError(err).Error("gql error presenter")
		return graphql.DefaultErrorPresenter(ctx, err)
	})

	// GraphQL endpoint
	engine.POST("/api", func(c *gin.Context) {
		ctx := log.SetLoggerInContext(c.Request.Context())
		gqlHandler.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	})

	// GraphQL Playground
	engine.GET("/api", basicAuth, func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})
}

func InitializeRoutes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "ELO 👋🏾",
		})
	})

	engine.GET("/health", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "Server is up 🚀",
		})
	})

	basicAuth := gin.BasicAuth(gin.Accounts{
		configuration.Instance.BasicUsername: configuration.Instance.BasicPassword,
	})

	initializeRestAPI(engine)
	initalizeGQLRoute(engine, basicAuth)
	initializeAsynqServer(engine, basicAuth)
}

func initializeAsynqServer(engine *gin.Engine, basicAuth gin.HandlerFunc) {
	monitorConf := &asynqMonitorConfig{}
	configuration.Load(monitorConf)

	parsedOpt, err := asynq.ParseRedisURI(monitorConf.RedisURL)
	if err != nil {
		log.WithError(err).Error("failed to parse REDIS_URL for asynq monitor")
	}

	monitorHandler := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring",
		RedisConnOpt: parsedOpt,
	})

	engine.Any(monitorHandler.RootPath(), basicAuth, gin.WrapH(monitorHandler))
	engine.Any(monitorHandler.RootPath()+"/*path", basicAuth, gin.WrapH(monitorHandler))
}

func initializeRestAPI(engine *gin.Engine) {
	var (
		// verificationController = verification.NewVerificationController(verification.NewService())
		dojahController = dojah.New(runtime.Injector)
	)

	// ****************************************
	// ********* Verification Route ***********
	// ****************************************
	// verificationController.RegisterRoutes(engine)

	// ****************************************
	// ********* Dojah Route ***********
	// ****************************************
	dojahController.RegisterRoutes(engine)

}
