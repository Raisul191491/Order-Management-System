package routes

import (
	"fmt"
	"oms/config"
	"oms/connection"

	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	masterDB, replicaDB := connection.InitDB(config.Conf)
	redis := connection.GetRedis(config.Conf)

	fmt.Println(masterDB, replicaDB, redis)

}
