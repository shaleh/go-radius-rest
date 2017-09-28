package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

func handleFailure(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	switch err {
	case UnknownUser:
		code = http.StatusNotFound
	case sql.ErrNoRows:
		panic("unhandled empty rows")
	}
	ctx.String(code, err.Error())
}

func handleUser(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		handleFailure(c, err)
		return
	}

	json.NewEncoder(c.Writer).Encode(user)
}

func handleUserDisable(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		handleFailure(c, err)
		return
	}
	err = user.Disable()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}

	if e, ok := err.(*mysql.MySQLError); ok {
		if e.Number == 1062 {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
			return
		}
	}

	handleFailure(c, err)
}

func handleUserReActivate(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		handleFailure(c, err)
		return
	}
	err = user.ReActivate()
	if err != nil {
		handleFailure(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func handleUserStatus(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		handleFailure(c, err)
		return
	}
	disabled, err := user.IsDisabled()
	if err != nil {
		handleFailure(c, err)
		return
	}

	json.NewEncoder(c.Writer).Encode(map[string]bool{"disabled": disabled})
}

func handleUsers(c *gin.Context) {
	db := getDB()
	users, err := getUsers(db)
	if err != nil {
		handleFailure(c, err)
		return
	}

	json.NewEncoder(c.Writer).Encode(users)
}
