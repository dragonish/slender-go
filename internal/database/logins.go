package database

import (
	"slender/internal/global"
	"slender/internal/logger"
	"slender/internal/model"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// AddLogin records login log.
func AddLogin(loginID string, loginTime time.Time, ip, ua string, isAdmin bool, maxAge uint16) error {
	log := logger.New("login_id", loginID, "login_time", loginTime, "ip", ip, "user_agent", ua, "is_admin", isAdmin, "max_age", maxAge)

	_, err := db.Exec("insert into logins(login_id, login_time, ip, ua, is_admin, max_age, active) values(?, ?, ?, ?, ?, ?, ?)", loginID, loginTime, ip, ua, isAdmin, maxAge, true)
	if err == nil {
		log.Info("logged in")
	} else {
		return log.Err("error recording login log", err)
	}

	global.Flags.LoginIDs = append(global.Flags.LoginIDs, loginID)
	return nil
}

// Logout records logout log.
func Logout(LoginID string) error {
	log := logger.New("login_id", LoginID)

	_, err := db.Exec("update logins set active=false where login_id=?", LoginID)
	if err == nil {
		log.Info("logged out")
	} else {
		return log.Err("error recording logout log", err)
	}

	global.Flags.LoginIDs = slices.DeleteFunc(global.Flags.LoginIDs, func(id string) bool { return id == LoginID })
	return nil
}

// GetLoginList gets login list.
func GetLoginList(cond *model.LoginListCondition, body *model.LoginListData) error {
	body.PageSize = cond.Size
	body.PageNo = cond.Page

	filter, params := getLoginFilterCondition(cond)
	if filter != "" {
		filter = " where " + filter
	}

	gStmt, err := db.PrepareNamed("select count(*) from logins l" + filter)
	if err != nil {
		return logger.Err("prepared login list count query statement error", err)
	}
	defer gStmt.Close()
	err = gStmt.Get(&body.Total, params)
	if err != nil {
		return logger.Err("get login list count error", err)
	} else if body.Total == 0 {
		return nil
	}

	o := getLoginListOrder(cond.Order)

	qStmt, err := db.PrepareNamed("select l.login_id, l.login_time, l.ip, l.ua, l.is_admin, l.max_age, l.active from logins l " + filter + " order by " + o + " limit :offset,:size")
	if err != nil {
		return logger.Err("prepared login list query statement error", err)
	}
	defer qStmt.Close()

	params["offset"] = cond.Size * (cond.Page - 1)
	params["size"] = cond.Size

	err = qStmt.Select(&body.List, params)
	if err != nil {
		return logger.Err("get login list error", err)
	}

	return nil
}

// ClearLogins clears login log.
func ClearLogins() error {
	_, err := db.Exec("delete from logins")
	if err != nil {
		return logger.Err("clear logins error", err)
	}

	logger.Info("clear logins")

	return nil
}

// GetLoginedList returns logined list.
func GetLoginedList(list *[]string) error {
	err := db.Select(list, "select login_id from logins where active = true")
	if err != nil {
		return logger.Err("get logined list error", err)
	}

	return nil
}

// ResetLoginedList resets login list.
func ResetLoginedList() error {
	_, err := db.Exec("update logins set active = null where active = true")
	if err != nil {
		return logger.Err("reset logined list error", err)
	}

	return nil
}

// LogoutAllLogins logouts all logins.
//
// Transaction operation.
func LogoutAllLogins(tx *sqlx.Tx) error {
	_, err := tx.Exec("update logins set active = false where active = true")
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return logger.Err("logout all logins error", err)
	}

	logger.Info("logout all logins")

	return nil
}

// getLoginFilterCondition returns login list filter condition.
func getLoginFilterCondition(cond *model.LoginListCondition) (string, map[string]interface{}) {
	condList := make([]string, 0)
	params := map[string]interface{}{}

	if cond.Admin != nil {
		condList = append(condList, "(l.is_admin = :is_admin)")
		params["is_admin"] = *cond.Admin
	}

	if cond.Active != nil {
		condList = append(condList, "(l.active = :active)")
		params["active"] = *cond.Active
	}

	if cond.IP != "" {
		condList = append(condList, `l.ip like :ip escape '\'`)
		params["ip"] = cond.IP.LikeMatchingString()
	}

	if cond.UA != "" {
		condList = append(condList, `l.ua like :ua escape '\'`)
		params["ua"] = cond.UA.LikeMatchingString()
	}

	return strings.Join(condList, " and "), params
}

func getLoginListOrder(order string) string {
	var o string

	switch order {
	case "login_time":
		fallthrough
	default:
		o = "l.login_time desc"
	}

	return o
}
