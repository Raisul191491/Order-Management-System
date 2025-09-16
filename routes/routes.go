package routes

import (
	"oms/config"
	"oms/connection"
	"oms/handler"
	"oms/middleware"
	"oms/repository"
	"oms/service"

	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	masterDB, replicaDB := connection.InitDB(config.Conf)
	//redis := connection.GetRedis(config.Conf)

	cityRepository := repository.NewCityRepository(masterDB, replicaDB)
	storeRepository := repository.NewStoreRepository(masterDB, replicaDB)
	zoneRepository := repository.NewZoneRepository(masterDB, replicaDB)
	itemTypeRepository := repository.NewItemTypeRepository(masterDB, replicaDB)
	deliveryTypeRepository := repository.NewDeliveryTypeRepository(masterDB, replicaDB)
	userRepository := repository.NewUserRepository(masterDB, replicaDB)
	userSessionRepository := repository.NewUserSessionRepository(masterDB, replicaDB)
	orderRepository := repository.NewOrderRepository(masterDB, replicaDB)

	cityService := service.NewCityService(cityRepository)
	storeService := service.NewStoreService(storeRepository)
	zoneService := service.NewZoneService(zoneRepository, cityRepository)
	itemTypeService := service.NewItemTypeService(itemTypeRepository)
	deliveryTypeService := service.NewDeliveryTypeService(deliveryTypeRepository)
	userService := service.NewUserService(userRepository)
	userSessionService := service.NewUserSessionService(userSessionRepository, config.Conf)
	authService := service.NewAuthService(userRepository, userSessionService)
	orderService := service.NewOrderService(orderRepository, storeService, cityService)

	cityHandler := handler.NewCityHandler(cityService)
	storeHandler := handler.NewStoreHandler(storeService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	itemTypeHandler := handler.NewItemTypeHandler(itemTypeService)
	deliveryTypeHandler := handler.NewDeliveryTypeHandler(deliveryTypeService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)

	omsRoutes := e.Group("/oms")

	omsRoutes.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	cityRoutes := omsRoutes.Group("/cities").Use(middleware.Auth(userSessionService))
	{
		cityRoutes.POST("", cityHandler.CreateCity)
		cityRoutes.GET("", cityHandler.GetAllCities)
		cityRoutes.GET("/:id", cityHandler.GetCityByID)
		cityRoutes.PUT("", cityHandler.UpdateCity)
		cityRoutes.DELETE("/:id", cityHandler.DeleteCity)
		cityRoutes.GET("/name/:name", cityHandler.GetCityByName)
	}

	storeRoutes := omsRoutes.Group("/stores").Use(middleware.Auth(userSessionService))
	{
		storeRoutes.POST("", storeHandler.CreateStore)
		storeRoutes.GET("", storeHandler.GetAllStores)
		storeRoutes.GET("/:id", storeHandler.GetStoreByID)
		storeRoutes.PUT("", storeHandler.UpdateStore)
		storeRoutes.DELETE("/:id", storeHandler.DeleteStore)
	}

	zoneRoutes := omsRoutes.Group("/zones").Use(middleware.Auth(userSessionService))
	{
		zoneRoutes.POST("", zoneHandler.CreateZone)
		zoneRoutes.GET("", zoneHandler.GetAllZones)
		zoneRoutes.GET("/:id", zoneHandler.GetZoneByID)
		zoneRoutes.PUT("", zoneHandler.UpdateZone)
		zoneRoutes.DELETE("/:id", zoneHandler.DeleteZone)
	}

	itemTypeRoutes := omsRoutes.Group("/item-types").Use(middleware.Auth(userSessionService))
	{
		itemTypeRoutes.POST("", itemTypeHandler.CreateItemType)
		itemTypeRoutes.GET("", itemTypeHandler.GetAllItemTypes)
		itemTypeRoutes.GET("/:id", itemTypeHandler.GetItemTypeByID)
		itemTypeRoutes.PUT("", itemTypeHandler.UpdateItemType)
		itemTypeRoutes.DELETE("/:id", itemTypeHandler.DeleteItemType)
	}

	deliveryTypeRoutes := omsRoutes.Group("/delivery-types").Use(middleware.Auth(userSessionService))
	{
		deliveryTypeRoutes.POST("", deliveryTypeHandler.CreateDeliveryType)
		deliveryTypeRoutes.GET("", deliveryTypeHandler.GetAllDeliveryTypes)
		deliveryTypeRoutes.GET("/:id", deliveryTypeHandler.GetDeliveryTypeByID)
		deliveryTypeRoutes.PUT("", deliveryTypeHandler.UpdateDeliveryType)
		deliveryTypeRoutes.DELETE("/:id", deliveryTypeHandler.DeleteDeliveryType)
	}

	userRoutes := omsRoutes.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.PUT("/email", userHandler.UpdateUserEmail)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	orderRoutes := omsRoutes.Group("/orders").Use(middleware.Auth(userSessionService))
	{
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.GET("/:consignment_id", orderHandler.GetOrderByConsignmentID)
		orderRoutes.GET("/all", orderHandler.ListAllOrders)
		orderRoutes.PUT("", orderHandler.UpdateOrder)
		orderRoutes.DELETE("/:id", orderHandler.DeleteOrder)
		orderRoutes.POST("/:consignment_id/cancel", orderHandler.CancelOrder)
	}

	loginRoutes := omsRoutes.Group("/auth")
	{
		loginRoutes.POST("/login", authHandler.Login)
	}

	logoutRoutes := omsRoutes.Group("/auth").Use(middleware.Auth(userSessionService))
	{
		logoutRoutes.POST("/logout", authHandler.Logout)
	}
}
