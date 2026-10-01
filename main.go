package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("env could not load", err)
	}
	loadJwtSecret()
	database := Connect()

	defer database.Close()

	r := gin.Default()
	protected := r.Group("/")
	protected.Use(authMilddleware())
	{
		protected.POST("/register", Reg(database))
		protected.GET("/login", login(database))
		protected.PATCH("/updatePost", updatePost(database))
		protected.DELETE("/deletePost", deletePost(database))
	}
	r.Run()
}
