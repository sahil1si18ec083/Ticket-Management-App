package main

import (
	"fmt"
	"os"
	"ticket-app-gin-golang/bootstrap"
)

func main() {
	bootstrap.LoadEnv()
	db := bootstrap.ConnectDatabase()
	r := bootstrap.InitServer(db)
	port := os.Getenv("PORT")
	fmt.Println("Server starting on port:", port)
	r.Run(":" + port)
}
