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

	cityService := service.NewCityService(cityRepository)

	cityHandler := handler.NewCityHandler(cityService)

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
}
