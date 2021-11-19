package main

import (
	"fmt"
	"offersapp/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}
	router := gin.Default()
	router.Use(dbMiddleWare(*conn))
	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UserRegister)
	}
	router.Run(":3000")

}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgres://posgres:@localhost:5432/offersapp")

	if err != nil {
		fmt.Println("Error connecting to DB: ", err)
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleWare(conn pgx.Conn) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}

}
