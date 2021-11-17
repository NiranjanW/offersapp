package main

import (
	"context"
	"fmt"
	"offersapp/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

func main() {

	// conn , err := ConnectDB()
	router := gin.Default()
	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UserRegister)
	}
	router.Run(":3000")

}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgres://posgres:@localhost:5432")

	if err != nil {
		fmt.Println("Error connecting to DB: ", err)
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware() gin.HandlerFuncs {

	return func(c *gin.Context) {
		c.Next()
	}

}
