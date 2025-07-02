package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"ia-boilerplate/src/handlers"
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/src/middlewares"
	"ia-boilerplate/src/repository"
	"net/http"
	"time"
)

func main() {

	logger, err := infrastructure.NewLogger()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
		return
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			panic("failed to sync logger: " + err.Error())
		}
	}(logger.Log)

	gin.DefaultWriter = zap.NewStdLog(logger.Log).Writer()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.GinZapLogger(), gin.Recovery())

	auth := infrastructure.NewAuth(logger)
	repo := &repository.Repository{
		Auth:   auth,
		Logger: logger,
	}

	logger.Info("Initializing database")
	if err := repo.InitDatabase(); err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		panic(err)
	}
	logger.Info("Database initialized", zap.Time("at", time.Now()))

	h := handlers.NewHandler(repo, logger, auth)

	c := cron.New()
	_, err = c.AddFunc("0 1 * * *", func() {
		logger.Info("Scheduled task executed", zap.Time("at", time.Now()))
	})
	if err != nil {
		logger.Error("Error setting up cron job", zap.Error(err))
	}
	c.Start()
	logger.Info("Cron scheduler started")
	defer c.Stop()

	SetupRoutes(router, h)
	logger.Info("Routes configured")

	logger.Info("Starting server", zap.String("address", "http://localhost:8080"))
	if err := router.Run(":8080"); err != nil {
		logger.Error("Server failed to start", zap.Error(err))
		panic(err)
	}
}

func SetupRoutes(router *gin.Engine, handler *handlers.Handler) {
	router.Use(middlewares.CorsMiddleware())
	router.Use(middlewares.Handler)
	r := router.Group("/")
	r.POST("/login", handler.Login)
	r.POST("/access-token/refresh", handler.AccessTokenByRefreshToken)
	api := r.Group("/api")

	api.Use(middlewares.JWTAuthMiddleware(handler))

	device := api.Group("/device")
	device.Use(middlewares.DeviceInfoInterceptor())
	device.GET("", func(c *gin.Context) {
		if deviceInfo, exists := c.Get("deviceInfo"); exists {
			c.JSON(http.StatusOK, deviceInfo)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Device info not found"})
		}
	})

	healthCheckAuth := api.Group("/health-check-auth")
	healthCheckAuth.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	userRoutes := api.Group("/users")
	{
		userRoutes.GET("", handler.GetUsers)
		userRoutes.GET("/:id", handler.GetUser)
		userRoutes.POST("", handler.CreateUser)
		userRoutes.PUT("/:id", handler.UpdateUser)
		userRoutes.DELETE("/:id", handler.DeleteUser)

		roleRoutes := userRoutes.Group("/roles")
		{
			roleRoutes.GET("", handler.GetRoles)
			roleRoutes.GET("/:id", handler.GetRole)
			roleRoutes.POST("", handler.CreateRole)
			roleRoutes.PUT("/:id", handler.UpdateRole)
			roleRoutes.DELETE("/:id", handler.DeleteRole)
		}

		deviceRoutes := userRoutes.Group("/devices")
		{
			deviceRoutes.GET("/user-id/:userId", handler.GetDevicesByUser)
			deviceRoutes.GET("/:id", handler.GetDevice)
			deviceRoutes.POST("", handler.CreateDevice)
			deviceRoutes.PUT("/:id", handler.UpdateDevice)
			deviceRoutes.DELETE("/:id", handler.DeleteDevice)
			deviceRoutes.GET("/search-paginated", handler.SearchDeviceDetailsPaginated)
			deviceRoutes.GET("/search-by-property", handler.SearchDeviceCoincidencesByProperty)
		}
	}

	medicineRoutes := api.Group("/medicines")
	{
		medicineRoutes.GET("/:id", handler.GetMedicine)
		medicineRoutes.POST("", handler.CreateMedicine)
		medicineRoutes.DELETE("/:id", handler.DeleteMedicine)
		medicineRoutes.GET("/search-paginated", handler.SearchMedicinesPaginated)
		medicineRoutes.GET("/search-by-property", handler.SearchMedicineCoincidencesByProperty)
	}

	icdcieRoutes := api.Group("/icd-cie")
	{
		icdcieRoutes.GET("", handler.GetICDCies)
		icdcieRoutes.GET("/:id", handler.GetICDCie)
		icdcieRoutes.POST("", handler.CreateICDCie)
		icdcieRoutes.PUT("/:id", handler.UpdateICDCie)
		icdcieRoutes.DELETE("/:id", handler.DeleteICDCie)
		icdcieRoutes.GET("/search-paginated", handler.SearchICDCiePaginated)
		icdcieRoutes.GET("/search-by-property", handler.SearchIcdCoincidencesByProperty)
	}
}
