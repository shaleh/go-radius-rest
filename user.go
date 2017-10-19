package main

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
)

type RadiusUser struct {
	db       *sql.DB
	Name     string `json:"name"`
	Disabled bool   `json:"disabled"`
}

var UnknownUser = errors.New("Unknown user")

func FindRadiusUser(db *sql.DB, name string) (user *RadiusUser, err error) {
	stmt, err := db.Prepare("select id from radcheck where username=?")
	if err != nil {
		return
	}
	var id int
	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = UnknownUser
		}
		return
	}

	return NewRadiusUser(db, name)
}

func NewRadiusUser(db *sql.DB, name string) (user *RadiusUser, err error) {
	user = &RadiusUser{db: db, Name: name}
	user.Disabled, err = user.IsDisabled()
	if err != nil {
		return
	}

	return
}

func (user *RadiusUser) IsDisabled() (bool, error) {
	stmt, err := user.db.Prepare("SELECT 1 FROM radusergroup WHERE username=? and groupname='disabled'")
	if err != nil {
		return false, err
	}

	var value int
	err = stmt.QueryRow(user.Name).Scan(&value)
	switch err {
	case nil:
		user.Disabled = true
	case sql.ErrNoRows:
		err = nil // now rows, user is not disabled
		user.Disabled = false
	}

	return user.Disabled, err
}

func (user *RadiusUser) Disable() error {
	stmt, err := user.db.Prepare("INSERT INTO radusergroup (username, groupname, priority) values(?,'disabled', 1)")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(user.Name)
	if err == nil {
		_, err = result.RowsAffected()
	}
	return err
}

func (user *RadiusUser) ReActivate() error {
	stmt, err := user.db.Prepare("DELETE FROM radusergroup WHERE username=? and groupname='disabled'")
	if err == nil {
		_, err = stmt.Exec(user.Name)
	}
	return err
}

func getUsers(db *sql.DB) (users []*RadiusUser, err error) {
	stmt, err := db.Prepare("SELECT DISTINCT username FROM radcheck")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return
		}

		var user *RadiusUser
		user, err = NewRadiusUser(db, username)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

func getUser(c *gin.Context) (user *RadiusUser, err error) {
	db := getDB()
	name := c.Param("name")
	user, err = FindRadiusUser(db, name)
	return
}
