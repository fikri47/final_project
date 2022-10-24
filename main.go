package main

import (
	"final_project/database"
	"final_project/routers"
)

func main() {
	database.StartDB()
	r := routers.StartApp()
	r.Run(":8080")
}
