package database

import (
	"database/sql"
	"math"
	"slender/internal/data"
	"slender/internal/logger"
	"slender/internal/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

// AddFolder adds new folder.
//
// Causes panic when a transaction error occurs.
func AddFolder(body *model.FolderPostBody) (int64, error) {
	log := logger.New("name", body.Name)

	tx := db.MustBegin()

	res, err := tx.NamedExec("insert into folders(name, description, large, privacy, weight, created_time, modified_time) values(:name, :description, :large, :privacy, :weight, datetime('now', 'localtime'), datetime('now', 'localtime'))", body)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("add folder error", err)
	}

	folderID, err := res.LastInsertId()
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("read the created folder id error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	return folderID, nil
}

// GetFolderList gets folder list.
func GetFolderList(cond *model.FolderListCondition, body *model.FolderListData) error {
	body.PageSize = cond.Size
	body.PageNo = cond.Page

	filter, params := getFolderFilterCondition(cond)
	if filter != "" {
		filter = " where " + filter
	}

	gStmt, err := db.PrepareNamed("select count(*) from folders f" + filter)
	if err != nil {
		return logger.Err("prepared folder list count query statement error", err)
	}
	defer gStmt.Close()
	err = gStmt.Get(&body.Total, params)
	if err != nil {
		return logger.Err("get folder list count error", err)
	} else if body.Total == 0 {
		return nil
	}

	o := getFolderListOrder(cond.Order)

	qStmt, err := db.PrepareNamed("select f.id, f.name, f.description, f.large, f.privacy, f.weight, f.created_time, f.modified_time, count(b.id) as bookmark_total from folders f left outer join bookmarks b on f.id = b.folder_id " + filter + " group by f.id order by " + o + " limit :offset,:size")
	if err != nil {
		return logger.Err("prepared folder list query statement error", err)
	}
	defer qStmt.Close()

	params["offset"] = cond.Size * (cond.Page - 1)
	params["size"] = cond.Size

	err = qStmt.Select(&body.List, params)
	if err != nil {
		return logger.Err("get folder list error", err)
	}

	return nil
}

// GetHomeFolderList gets folder list used by the homepage.
func GetHomeFolderList(list *[]model.HomeFolderListItem) error {
	err := db.Select(list, "select id, name, description, large from folders order by large desc, weight desc, id")
	if err != nil {
		return logger.Err("get folder list used by the homepage error", err)
	}

	return nil
}

// GetBookmarkFolderInfoList gets folder list used by the bookmark.
func GetBookmarkFolderList(list *[]model.BookmarkFolderInfo) error {
	err := db.Select(list, "select id, name, privacy from folders order by large desc, weight desc, id")
	if err != nil {
		return logger.Err("get folder list used by the bookmark error", err)
	}

	return nil
}

// GetFolder gets folder item data.
func GetFolder(folderID int64, body *model.FolderBaseData) error {
	log := logger.New("folder_id", folderID)

	err := db.Get(body, "select * from folders where id = ?", folderID)
	if err == sql.ErrNoRows {
		log.Info("folder does not exist")
		return model.ErrNotExist
	} else if err != nil {
		return log.Err("get folder base info error", err)
	}

	return nil
}

// UpdateFolder updates folder.
func UpdateFolder(folderID int64, body *model.FolderPatchBody) error {
	log := logger.New("folder_id", folderID)

	cond := make([]string, 0)
	m := data.StructToMap(*body)

	for k, v := range m {
		if k == "name" {
			//? Not allowed to set the name as empty.
			switch s := v.(type) {
			case string, model.MyString:
				if s == "" {
					delete(m, "name")
				} else {
					cond = append(cond, k+" = :"+k)
				}
			default:
				delete(m, "name")
			}
		} else {
			cond = append(cond, k+" = :"+k)
		}
	}

	if len(cond) > 0 {
		//* First check if the folder exists
		state, qErr := isFolderExists(folderID)
		if qErr == nil {
			if !state {
				log.Info("folder does not exist")
				return model.ErrNotExist
			}
		} else {
			return qErr
		}

		m["id"] = folderID
		_, err := db.NamedExec("update folders set "+strings.Join(cond, ", ")+", modified_time = datetime('now', 'localtime') where id = :id", m)
		if err != nil {
			return log.Err("update folder error", err)
		}
	} else {
		log.Info("no data that needs to be processed")
		return model.ErrDoNothing
	}

	return nil
}

// DeleteFolder deletes folder.
func DeleteFolder(folderID int64) error {
	log := logger.New("folder_id", folderID)

	_, err := db.Exec("delete from folders where id = ?", folderID)
	if err != nil {
		return log.Err("delete folder error", err)
	}

	log.Info("delete folder")

	return nil
}

// FolderBatchHandler handles folder in batches.
//
// Causes panic when a transaction error occurs.
func FolderBatchHandler(body *model.BatchPatchBody) error {
	if len(body.DataSet) == 0 {
		return nil
	}

	tx := db.MustBegin()

	if body.Action == "delete" {
		query, args, err := sqlx.In("delete from folders where id in (?)", body.DataSet)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("in for delete folders statement error", err)
		}
		_, err = tx.Exec(query, args...)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("delete folders error", err)
		}
	} else if body.Action == "setLarge" || body.Action == "setPrivacy" {
		b, ok := body.Payload.(bool)
		if ok {
			column := "large"
			if body.Action == "setPrivacy" {
				column = "privacy"
			}

			log := logger.New("column", column)

			query, args, err := sqlx.In("update folders set modified_time = datetime('now', 'localtime'), "+column+" = ? where id in (?)", b, body.DataSet)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update folders statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update folders error", err)
			}
		} else {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return model.ErrQueryParamMissing
		}
	} else if body.Action == "setWeight" || body.Action == "incWeight" {
		i, ok := body.Payload.(float64)
		if ok {
			log := logger.New("column", "weight")

			var params []interface{}
			str := ""
			if body.Action == "setWeight" {
				str = "update folders set modified_time = datetime('now', 'localtime'), weight = ? where id in (?)"
				params = []interface{}{i, body.DataSet}
			} else {
				str = "update folders set modified_time = datetime('now', 'localtime'), weight = weight " + data.Int16ToStringWithSign(int16(math.Round(i))) + " where id in (?)"
				params = []interface{}{body.DataSet}
			}

			query, args, err := sqlx.In(str, params...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update folders statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update folders error", err)
			}
		} else {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return model.ErrQueryParamMissing
		}
	} else {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return model.ErrDoNothing
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	return nil
}

// getFolderFilterCondition returns folder list filter condition.
func getFolderFilterCondition(cond *model.FolderListCondition) (string, map[string]interface{}) {
	condList := make([]string, 0)
	params := map[string]interface{}{}

	if cond.Privacy != nil {
		condList = append(condList, ("f.privacy = :privacy"))
		params["privacy"] = *cond.Privacy
	}

	if cond.Name != "" {
		condList = append(condList, "(instr(f.name, :name) > 0)")
		params["name"] = cond.Name
	}

	if cond.Des != "" {
		condList = append(condList, "(instr(f.description, :description) > 0)")
		params["description"] = cond.Des
	}

	return strings.Join(condList, " and "), params
}

// getFolderListOrder returns folder list order condition.
//
// Optional order values: created-time | modified-time | bookmark-total | weight.
func getFolderListOrder(order string) string {
	var o string

	switch order {
	case "created-time":
		o = "f.created_time desc"
	case "modified-time":
		o = "f.modified_time desc"
	case "bookmark-total":
		o = "bookmark_total desc"
	case "weight":
		fallthrough
	default:
		o = "f.large desc, f.weight desc"
	}

	return o
}

// isFolderExists returns true when the folder exists.
func isFolderExists(folderID int64) (bool, error) {
	log := logger.New("folder_id", folderID)

	var readID model.MyInt64
	err := db.Get(&readID, "select id from folders where id = ?", folderID)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	} else {
		return false, log.Err("an error occurred while checking whether the folder exists", err)
	}
}
