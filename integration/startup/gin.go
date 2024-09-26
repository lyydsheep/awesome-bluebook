package startup

import "github.com/gin-gonic/gin"

func InitGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	return server
}
