package main

import (
	"api-golang/database"
	"api-golang/helpers"
	"api-golang/router"
)

func init() {
	helpers.LoadEnvVariables()
	database.Connection()
	database.SyncDatabase()
}

func main() {
	router.Router().Run()
}
