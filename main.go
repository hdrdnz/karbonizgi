package main

import (
	router "carbonfootprint/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.Load(r)
	r.Run(":8000")

}
