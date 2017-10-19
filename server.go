package main

import (
	"fmt"
	"github.com/galdor/go-cmdline"
	"github.com/gin-gonic/gin"
	"os"
)

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	cmdline := cmdline.New()
	cmdline.AddArgument("configFile", "the config file")
	cmdline.Parse(os.Args)

	config, err := loadConfiguration(cmdline.ArgumentValue("configFile"))
	if err != nil {
		die("Could not load config: %s", err.Error())
	}

	_db, err = sqlConnect(config)
	if err != nil {
		die("Failure to connect to database: %s\n", err.Error())
	}
	defer _db.Close()

	router := gin.Default()

	router.GET("/users/", handleUsers)
	router.GET("/users/:name/", handleUser)
	router.GET("/users/:name/status", handleUserStatus)
	router.PUT("/users/:name/disable", handleUserDisable)
	router.PUT("/users/:name/reactivate", handleUserReActivate)

	router.Run(config.BindAddress())
}
