package database

import (
	"database/sql"

	"slender/internal/data"
	"slender/internal/global"
	"slender/internal/logger"
	"slender/internal/model"
)

// Register handles register user action.
//
// Record a set of salt and secret for the specified username.
func Register() error {
	type User struct {
		Salt   model.MyString `db:"salt"`
		Secret model.MyString `db:"secret"`
	}
	var user User

	log := logger.New("username", model.DB_USERNAME)

	err := db.Get(&user, "select u.salt, u.secret from users u where u.username = ?", model.DB_USERNAME)
	if err == sql.ErrNoRows {
		if global.Flags.AccessPassword != "" || global.Flags.AdminPassword != "" {
			global.Flags.Salt = data.SaltGenerator(64)
			global.Flags.Secret = data.SaltGenerator(16)

			_, err := db.Exec("insert into users(username, salt, secret) values(?, ?, ?)", model.DB_USERNAME, global.Flags.Salt, global.Flags.Secret)
			if err != nil {
				return log.Err("register user error", err)
			}
			log.Debug("register user")
		}
	} else if err != nil {
		return log.Err("error checking username exists", err)
	} else {
		if global.Flags.AccessPassword == "" && global.Flags.AdminPassword == "" {
			//? Delete user from database.
			_, err := db.Exec("delete from users where username = ?", model.DB_USERNAME)
			if err != nil {
				return log.Err("delete user error", err)
			}
			log.Debug("deleted user")
		} else {
			global.Flags.Salt = user.Salt.String()
			global.Flags.Secret = user.Secret.String()
		}
	}

	if global.Flags.AccessPassword != "" {
		global.Flags.AccessToken = data.Sha256Generator(global.Flags.AccessPassword, global.Flags.Salt)
	}
	if global.Flags.AdminPassword != "" {
		global.Flags.AdminToken = data.Sha256Generator(global.Flags.AdminPassword, global.Flags.Salt)
	}

	return nil
}

// LogoutAll logouts all users.
//
// Reset the salt and secret for the specified username.
func LogoutAll() error {
	log := logger.New("username", model.DB_USERNAME)

	//? Delete user from database.
	_, err := db.Exec("delete from users where username = ?", model.DB_USERNAME)
	if err != nil {
		return log.Err("delete user error", err)
	}
	log.Debug("deleted user")

	if global.Flags.AccessPassword != "" || global.Flags.AdminPassword != "" {
		global.Flags.Salt = data.SaltGenerator(64)
		global.Flags.Secret = data.SaltGenerator(16)

		_, err := db.Exec("insert into users(username, salt, secret) values(?, ?, ?)", model.DB_USERNAME, global.Flags.Salt, global.Flags.Secret)
		if err != nil {
			return log.Err("register user error", err)
		}
		log.Debug("re-register user")
	}

	if global.Flags.AccessPassword != "" {
		global.Flags.AccessToken = data.Sha256Generator(global.Flags.AccessPassword, global.Flags.Salt)
	}
	if global.Flags.AdminPassword != "" {
		global.Flags.AdminToken = data.Sha256Generator(global.Flags.AdminPassword, global.Flags.Salt)
	}

	return nil
}
