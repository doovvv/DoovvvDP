package main

import (
	"doovvvDP/config"
	"doovvvDP/dal/mysql"
	"doovvvDP/router"
)
func main() {
	config.ConfigInit()
	mysql.Init()
	router.RouterInit()
} 