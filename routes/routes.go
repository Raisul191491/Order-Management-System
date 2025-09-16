package routes

import (
	"oms/config"
	"oms/connection"
	"oms/handler"
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

	cityService := service.NewCityService(cityRepository)
	storeService := service.NewStoreService(storeRepository)
	zoneService := service.NewZoneService(zoneRepository, cityRepository)
	itemTypeService := service.NewItemTypeService(itemTypeRepository)
	deliveryTypeService := service.NewDeliveryTypeService(deliveryTypeRepository)
	userService := service.NewUserService(userRepository)

	cityHandler := handler.NewCityHandler(cityService)
	storeHandler := handler.NewStoreHandler(storeService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	itemTypeHandler := handler.NewItemTypeHandler(itemTypeService)
	deliveryTypeHandler := handler.NewDeliveryTypeHandler(deliveryTypeService)
	userHandler := handler.NewUserHandler(userService)

	omsRoutes := e.Group("/oms")

	omsRoutes.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	cityRoutes := omsRoutes.Group("/cities")
	{
		cityRoutes.POST("", cityHandler.CreateCity)
		cityRoutes.GET("", cityHandler.GetAllCities)
		cityRoutes.GET("/:id", cityHandler.GetCityByID)
		cityRoutes.PUT("", cityHandler.UpdateCity)
		cityRoutes.DELETE("/:id", cityHandler.DeleteCity)
		cityRoutes.GET("/name/:name", cityHandler.GetCityByName)
	}

	storeRoutes := omsRoutes.Group("/stores")
	{
		storeRoutes.POST("", storeHandler.CreateStore)
		storeRoutes.GET("", storeHandler.GetAllStores)
		storeRoutes.GET("/:id", storeHandler.GetStoreByID)
		storeRoutes.PUT("", storeHandler.UpdateStore)
		storeRoutes.DELETE("/:id", storeHandler.DeleteStore)
	}

	zoneRoutes := omsRoutes.Group("/zones")
	{
		zoneRoutes.POST("", zoneHandler.CreateZone)
		zoneRoutes.GET("", zoneHandler.GetAllZones)
		zoneRoutes.GET("/:id", zoneHandler.GetZoneByID)
		zoneRoutes.PUT("", zoneHandler.UpdateZone)
		zoneRoutes.DELETE("/:id", zoneHandler.DeleteZone)
	}

	itemTypeRoutes := omsRoutes.Group("/item-types")
	{
		itemTypeRoutes.POST("", itemTypeHandler.CreateItemType)
		itemTypeRoutes.GET("", itemTypeHandler.GetAllItemTypes)
		itemTypeRoutes.GET("/:id", itemTypeHandler.GetItemTypeByID)
		itemTypeRoutes.PUT("", itemTypeHandler.UpdateItemType)
		itemTypeRoutes.DELETE("/:id", itemTypeHandler.DeleteItemType)
	}

	deliveryTypeRoutes := omsRoutes.Group("/delivery-types")
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

}
