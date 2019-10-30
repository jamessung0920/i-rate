package main

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	// "github.com/jinzhu/gorm"
	
	"net/http"
	
	// "app/math"
	"app/common"
	"app/currency"
	"app/webhook"
	"app/database"
	"app/database/migrations"
	// "app/database/models"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = common.NewLogger()
}

func main() {
	Log.Info("start golang...")

	//connect database
	db, err := database.ConnectionDB()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	//migration
	migrations.Migrate()

	//routing
	router := gin.Default()

	currency.AddRoute(router)
	webhook.AddRoute(router)

	//test route
	router.GET("test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test2",
		})
	})

	router.Run()
}