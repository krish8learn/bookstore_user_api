package app

import (
	"github.com/gin-gonic/gin"
)

//2nd layer
var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	router.Run(":8085")
}
