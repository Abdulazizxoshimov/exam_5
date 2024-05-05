package api

import (
	// "github.com/casbin/casbin/v2"
	"time"

	v1 "exam_5/admin_api_gateway/api/handlers/v1"
	"exam_5/admin_api_gateway/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	redisrepo "exam_5/admin_api_gateway/internal/infrastructure/repository/redis"

	grpcClients "exam_5/admin_api_gateway/internal/infrastructure/grpc_service_client"
	"exam_5/admin_api_gateway/internal/pkg/config"
	"exam_5/admin_api_gateway/internal/usecase/event"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer  event.BrokerProducer
	Cache          redisrepo.Cache
	

}

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		BrokerProducer:  option.BrokerProducer,
		Redis: option.Cache,
	})

	


	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))
	
	router.Use(middleware.Tracing)

	api := router.Group("/v1")
	
	
// users
api.POST("/users/create", HandlerV1.CreateUser)
api.GET("/users/get/:id", HandlerV1.GetUser)
 api.GET("/users/list", HandlerV1.ListUsers)
api.PUT("/users/update/:id", HandlerV1.UpdateUser)
api.DELETE("/users/delete/:id", HandlerV1.DeleteUser)



//jobs
api.POST("/jobs/create", HandlerV1.Create)
api.GET("/jobs/get/:id", HandlerV1.Get)
api.GET("/jobs/list", HandlerV1.List)
api.PUT("/jobs/update/:id", HandlerV1.Update)
api.DELETE("/jobs/delete/:id", HandlerV1.Delete)
api.GET("/jobs/listwithowner", HandlerV1.ListWithOwner)
api.GET("/jobs/listbyclientid/:id", HandlerV1.ListByClientId)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}


