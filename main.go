package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"ia-boilerplate/src/db"
	"ia-boilerplate/src/handlers"
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/src/logger"
	middlewares2 "ia-boilerplate/src/middlewares"
	"net/http"
	"time"
)

func main() {
	if err := logger.Init(); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Log.Sync()

	// Configura salida de logs de Gin a stdout (opcional)
	gin.DefaultWriter = zap.NewStdLog(logger.Log).Writer()
	gin.SetMode(gin.ReleaseMode)

	// Configura router con middleware zap
	router := gin.New()
	router.Use(logger.GinZapLogger(), gin.Recovery())

	router.SetFuncMap(infrastructure.FuncMapTemplates)
	router.LoadHTMLGlob("templates/*.*")

	// Log: inicialización de base de datos
	logger.Info("Initializing database")
	if err := db.InitDatabase(); err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		panic(err)
	}
	logger.Info("Database initialized", zap.Time("at", time.Now()))

	c := cron.New()
	_, err := c.AddFunc("0 1 * * *", func() {
		logger.Info("Scheduled task executed", zap.Time("at", time.Now()))
	})
	if err != nil {
		logger.Error("Error setting up cron job", zap.Error(err))
	}
	c.Start()
	defer c.Stop()

	SetupRoutes(router)

	// Arranca el servidor
	logger.Info("Server is running", zap.String("address", "http://localhost:8080"))
	if err := router.Run(":8080"); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
		panic(err)
	}
}

func SetupRoutes(router *gin.Engine) {
	router.Use(middlewares2.CorsMiddleware())
	router.Use(middlewares2.Handler)
	r := router.Group("/")
	r.POST("/login", handlers.Login)
	r.POST("/access-token/refresh", handlers.AccessTokenByRefreshToken)

	api := r.Group("/api")

	api.Use(middlewares2.JWTAuthMiddleware())

	device := api.Group("/device")
	device.Use(middlewares2.DeviceInfoInterceptor())
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
		userRoutes.GET("", handlers.GetUsers)
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.POST("", handlers.CreateUser)
		userRoutes.PUT("/:id", handlers.UpdateUser)
		userRoutes.DELETE("/:id", handlers.DeleteUser)

		roleRoutes := userRoutes.Group("/roles")
		{
			roleRoutes.GET("", handlers.GetRoles)
			roleRoutes.GET("/:id", handlers.GetRole)
			roleRoutes.POST("", handlers.CreateRole)
			roleRoutes.PUT("/:id", handlers.UpdateRole)
			roleRoutes.DELETE("/:id", handlers.DeleteRole)
		}

		deviceRoutes := userRoutes.Group("/devices")
		{
			deviceRoutes.GET("/user-id/:userId", handlers.GetDevicesByUser)
			deviceRoutes.GET("/:id", handlers.GetDevice)
			deviceRoutes.POST("", handlers.CreateDevice)
			deviceRoutes.PUT("/:id", handlers.UpdateDevice)
			deviceRoutes.DELETE("/:id", handlers.DeleteDevice)
		}

	}

	medicineRoutes := api.Group("/medicines")
	{
		medicineRoutes.GET("/:id", handlers.GetMedicine)
		medicineRoutes.POST("", handlers.CreateMedicine)
		medicineRoutes.PUT("/:id", handlers.UpdateMedicine)
		medicineRoutes.DELETE("/:id", handlers.DeleteMedicine)
		medicineRoutes.GET("/search-paginated", handlers.SearchMedicinesPaginated)
		medicineRoutes.GET("/search-by-property", handlers.SearchMedicineCoincidencesByProperty)

	}

	icdcieRoutes := api.Group("/icd-cie")
	{
		icdcieRoutes.GET("", handlers.GetICDCies)
		icdcieRoutes.GET("/:id", handlers.GetICDCie)
		icdcieRoutes.POST("", handlers.CreateICDCie)
		icdcieRoutes.PUT("/:id", handlers.UpdateICDCie)
		icdcieRoutes.DELETE("/:id", handlers.DeleteICDCie)
		icdcieRoutes.GET("/search-paginated", handlers.SearchICDCiePaginated)
		icdcieRoutes.GET("/search-by-property", handlers.SearchIcdCoincidencesByProperty)

	}

	clientRoutes := api.Group("/clients")
	{
		clientRoutes.GET("", handlers.GetClients)
		clientRoutes.GET("/:id", handlers.GetClient)
		clientRoutes.POST("", handlers.CreateClient)
		clientRoutes.PUT("/:id", handlers.UpdateClient)
		clientRoutes.DELETE("/:id", handlers.DeleteClient)
		clientRoutes.GET("/search", handlers.SearchClientsPaginated)
		clientRoutes.GET("/search-by-property", handlers.SearchClientCoincidencesByProperty)
	}

}
