package main

import (
	"avito-trainee/external/database"
	"avito-trainee/external/httpserver"
)

func main() {
	db := database.InitDB()
	httpserver.InitAndStart(db)
}
