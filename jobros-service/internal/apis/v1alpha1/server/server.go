package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "8080"
)

func SetupRouter(host string, port string) *gin.Engine {
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		// default port
		port = defaultPort
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	//v1 := router.Group("/api/v1")
	//{
	//	users := v1.Group("/users")
	//	{
	//		users.POST("/", CreateUser)
	//		users.GET("/", GetUsers)
	//		users.GET("/:id", GetUser)
	//		users.PUT("/:id", UpdateUser)
	//		users.DELETE("/:id", DeleteUser)
	//	}
	//
	//	profiles := v1.Group("/profiles")
	//	{
	//		profiles.POST("/", CreateProfile)
	//		profiles.GET("/", GetProfiles)
	//		profiles.GET("/:id", GetProfile)
	//		profiles.PUT("/:id", UpdateProfile)
	//		profiles.DELETE("/:id", DeleteProfile)
	//	}
	//
	//	services := v1.Group("/services")
	//	{
	//		services.POST("/", CreateService)
	//		services.GET("/", GetServices)
	//		services.GET("/:id", GetService)
	//		services.PUT("/:id", UpdateService)
	//		services.DELETE("/:id", DeleteService)
	//	}
	//}

	return router
}

// StartServer starts the server
func StartServer(host string, port string) error {
	router := SetupRouter(host, port)
	return router.Run(fmt.Sprintf("%s:%s", host, port))
}
