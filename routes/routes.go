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

	cityService := service.NewCityService(cityRepository)
	storeService := service.NewStoreService(storeRepository)
	zoneService := service.NewZoneService(zoneRepository, cityRepository)
	itemTypeService := service.NewItemTypeService(itemTypeRepository)

	cityHandler := handler.NewCityHandler(cityService)
	storeHandler := handler.NewStoreHandler(storeService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	itemTypeHandler := handler.NewItemTypeHandler(itemTypeService)

	omsRoutes := e.Group("/oms")

	omsRoutes.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	cityRoutes := omsRoutes.Group("/cities")
	{
		cityRoutes.POST("", cityHandler.CreateCity)
		cityRoutes.GET("", cityHandler.GetAllCities)
		cityRoutes.GET("/:id", cityHandler.GetCityByID)
		cityRoutes.PUT("/:id", cityHandler.UpdateCity)
		cityRoutes.DELETE("/:id", cityHandler.DeleteCity)
		cityRoutes.GET("/name/:name", cityHandler.GetCityByName)
	}

	storeRoutes := omsRoutes.Group("/stores")
	{
		storeRoutes.POST("", storeHandler.CreateStore)
		storeRoutes.GET("", storeHandler.GetAllStores)
		storeRoutes.GET("/:id", storeHandler.GetStoreByID)
		storeRoutes.PUT("/:id", storeHandler.UpdateStore)
		storeRoutes.DELETE("/:id", storeHandler.DeleteStore)
	}

	zoneRoutes := omsRoutes.Group("/zones")
	{
		zoneRoutes.POST("", zoneHandler.CreateZone)
		zoneRoutes.GET("", zoneHandler.GetAllZones)
		zoneRoutes.GET("/:id", zoneHandler.GetZoneByID)
		zoneRoutes.PUT("/:id", zoneHandler.UpdateZone)
		zoneRoutes.DELETE("/:id", zoneHandler.DeleteZone)
	}

	itemTypeRoutes := omsRoutes.Group("/item_types")
	{
		itemTypeRoutes.POST("", itemTypeHandler.CreateItemType)
		itemTypeRoutes.GET("", itemTypeHandler.GetAllItemTypes)
		itemTypeRoutes.GET("/:id", itemTypeHandler.GetItemTypeByID)
		itemTypeRoutes.PUT("/:id", itemTypeHandler.UpdateItemType)
		itemTypeRoutes.DELETE("/:id", itemTypeHandler.DeleteItemType)
	}

}
