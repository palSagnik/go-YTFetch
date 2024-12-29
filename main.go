package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/backend/config"
	"github.com/palSagnik/go-YTFetch.git/backend/database"
	"github.com/palSagnik/go-YTFetch.git/backend/routes"
)

func main() {


	// loop till database is initialised
	for {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			fmt.Println("waiting for 30 seconds before trying again")
			time.Sleep(time.Second * 30)
			continue
		}
		break
	}

	err := database.MigrateUp()
	if err != nil {
		panic(err)
	}

	// setting up web-app
	corsConfig := cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }
	
	r := gin.Default()
	r.Use(cors.New(corsConfig))
    routes.SetUpRoutes(r)
    
	port := fmt.Sprintf(":%d", config.APP_PORT)
    r.Run(port)
}