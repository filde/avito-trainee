package main

import (
	"avito-trainee/common/metrics"
	"avito-trainee/external/database"
	"avito-trainee/external/httpserver"
)

func main() {
	ms := metrics.Init()
	db := database.InitDB()
	httpserver.InitAndStart(db, ms)
}
