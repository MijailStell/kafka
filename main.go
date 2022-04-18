package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"compartamos/poc/kafka/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	gindump "github.com/tpkeeper/gin-dump"
)

const VIDEOS_ROOT = "/movies"

func setupLogOutPut() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func setupConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	mode, ok := viper.Get("SERVICE.MODE").(string)

	if ok {
		gin.SetMode(mode)
	}
	viper.SetDefault("SERVICE.PORT", "5000")
}

func main() {
	setupLogOutPut()
	setupConfig()
	server := gin.New()

	server.Use(gin.Recovery(),
		middleware.Logger(),
		gindump.Dump())

	// JWT Authorization Middleware applies to "/api" only.
	apiRoutes := server.Group("/api")
	{
		apiRoutes.GET(VIDEOS_ROOT, func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"title":       "Pulp Fiction",
				"description": "The lives of two mob hitmen, a boxer, a gangster's wife, and a pair of diner bandits intertwine in four tales of violence and redemption.",
				"genre":       "Crime, Drama",
				"releaseDate": "2022-04-18T16:09:03.605Z",
			})
		})
	}

	servicePort, ok := viper.Get("SERVICE.PORT").(string)
	if os.Getenv("ASPNETCORE_PORT") != "" {
		servicePort = os.Getenv("ASPNETCORE_PORT")
		fmt.Println(servicePort)
	}

	// if type assert is not valid it will throw an error
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	if servicePort == "" {
		servicePort = "5000"
	}

	server.Run(":" + servicePort)
}
