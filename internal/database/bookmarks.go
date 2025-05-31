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

// AddBookmark adds new bookmark.
//
// Causes panic when a transaction error occurs.
func AddBookmark(body *model.BookmarkPostBody) (int64, error) {
	log := logger.New("name", body.Name, "url", body.URL)

	tx := db.MustBegin()

	res, err := tx.NamedExec("insert into bookmarks(url, name, description, icon, privacy, weight, created_time, modified_time, folder_id, hide_in_other) values(:url, :name, :description, :icon, :privacy, :weight, datetime('now', 'localtime'), datetime('now', 'localtime'), :folder_id, :hide_in_other)", body)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("add bookmark error", err)
	}

	bookmarkID, err := res.LastInsertId()
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("read the created bookmark id error", err)
	}

	if len(body.Files) > 0 {
		for _, item := range body.Files {
			err := updateFile(tx, item, model.MyInt64(bookmarkID))
			if err != nil {
				return 0, err
			}
		}
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	log.Info("add new bookmark", "bookmark_id", bookmarkID)

	return bookmarkID, nil
}

// GetBookmarkList gets bookmark list.
func GetBookmarkList(cond *model.BookmarkListCondition, body *model.BookmarkListData) error {
	body.PageSize = cond.Size
	body.PageNo = cond.Page

	filter, params := getBookmarkFilterCondition(cond)

	logger.Debug("bookmark filter condition", "filter", filter, "params", params)

	if filter != "" {
		filter = " where " + filter
	}

	gStmt, err := db.PrepareNamed("select count(*) from bookmarks b" + filter)
	if err != nil {
		return logger.Err("prepared bookmark list count query statement error", err)
	}
	defer gStmt.Close()
	err = gStmt.Get(&body.Total, params)
	if err != nil {
		return logger.Err("get bookmark list count error", err)
	} else if body.Total == 0 {
		return nil
	}

	o := getBookmarkListOrder(cond.Order)

	qStmt, err := db.PrepareNamed("select b.id, b.url, b.name, b.description, b.icon, b.privacy, b.weight, b.created_time, b.modified_time, b.visits, b.folder_id, b.hide_in_other, f.name as folder_name from bookmarks b left outer join folders f on b.folder_id = f.id " + filter + " order by " + o + " limit :offset,:size")
	if err != nil {
		return logger.Err("prepared bookmark list query statement error", err)
	}
	defer qStmt.Close()

	params["offset"] = cond.Size * (cond.Page - 1)
	params["size"] = cond.Size

	err = qStmt.Select(&body.List, params)
	if err != nil {
		return logger.Err("get bookmark list error", err)
	}

	return nil
}

// GetHomeBookmarkList gets bookmark list used by the homepage.
func GetHomeBookmarkList(privacy bool, inOtherNetwork bool, list *[]model.HomeBookmarkListItem) error {
	otherCond := ""
	if inOtherNetwork {
		otherCond = " and (b.hide_in_other = false)"
	}

	sqlStr := "select b.id, b.url, b.name, b.description, b.icon, b.folder_id, b.hide_in_other from bookmarks b where (b.privacy = false)" + otherCond + " and ((b.folder_id is null) or (b.folder_id in (select f.id from folders f where f.privacy = false))) order by b.weight desc, b.id"
	if privacy {
		if inOtherNetwork {
			otherCond = " where b.hide_in_other = false"
		}
		sqlStr = "select b.id, b.url, b.name, b.description, b.icon, b.folder_id, b.hide_in_other from bookmarks b" + otherCond + " order by b.weight desc, b.id"
	}

	err := db.Select(list, sqlStr)
	if err != nil {
		return logger.Err("get bookmark list used by the homepage error", err, "privacy", privacy, "in_other_newwork", inOtherNetwork)
	}

	return nil
}

// GetHomeLatestBookmarkList gets latest bookmark list used by the homepage.
func GetHomeLatestBookmarkList(privacy bool, inOtherNetwork bool, size uint8, list *[]model.HomeBookmarkListItem) error {
	otherCond := ""
	if inOtherNetwork {
		otherCond = " and (b.hide_in_other = false)"
	}

	sqlStr := "select b.id, b.url, b.name, b.description, b.icon, b.hide_in_other from bookmarks b where (b.privacy = false)" + otherCond + " and ((b.folder_id is null) or (b.folder_id in (select f.id from folders f where f.large = false and f.privacy = false))) and (b.created_time >= datetime('now', '-15 days')) order by b.created_time desc, b.weight desc, b.id limit ?"
	if privacy {
		if inOtherNetwork {
			otherCond = " (b.hide_in_other = false) and"
		}
		sqlStr = "select b.id, b.url, b.name, b.description, b.icon, b.hide_in_other from bookmarks b where" + otherCond + " ((b.folder_id is null) or (b.folder_id in (select f.id from folders f where f.large = false))) and (b.created_time >= datetime('now', '-15 days')) order by b.created_time desc, b.weight desc, b.id limit ?"
	}

	err := db.Select(list, sqlStr, size)
	if err != nil {
		return logger.Err("get latest bookmark list used by the homepage error", err, "privacy", privacy, "in_other_network", inOtherNetwork)
	}

	return nil
}

// GetHomeHotBookmarkList gets hot bookmark list used by the homepage.
func GetHomeHotBookmarkList(privacy bool, inOtherNetwork bool, size uint8, list *[]model.HomeBookmarkListItem) error {
	otherCond := ""
	if inOtherNetwork {
		otherCond = " and (b.hide_in_other = false)"
	}

	sqlStr := "select b.id, b.url, b.name, b.description, b.icon, b.hide_in_other from bookmarks b where (b.privacy = false)" + otherCond + " and ((b.folder_id is null) or (b.folder_id in (select f.id from folders f where f.large = false and f.privacy = false))) and b.visits > 0 order by b.visits desc, b.weight desc, b.id limit ?"
	if privacy {
		if inOtherNetwork {
			otherCond = " (b.hide_in_other = false) and"
		}
		sqlStr = "select b.id, b.url, b.name, b.description, b.icon, b.hide_in_other from bookmarks b where" + otherCond + " ((b.folder_id is null) or (b.folder_id in (select f.id from folders f where f.large = false))) and b.visits > 0 order by b.visits desc, b.weight desc, b.id limit ?"
	}

	err := db.Select(list, sqlStr, size)
	if err != nil {
		return logger.Err("get hot bookmark list used by the homepage error", err, "privacy", privacy, "in_other_network", inOtherNetwork)
	}

	return nil
}

// GetBookmark gets bookmark item data.
func GetBookmark(bookmarkID int64, body *model.BookmarkResData) error {
	log := logger.New("bookmark_id", bookmarkID)

	err := db.Get(body, "select * from bookmarks where id = ?", bookmarkID)
	if err == sql.ErrNoRows {
		log.Info("bookmark does not exist")
		return model.ErrNotExist
	} else if err != nil {
		return log.Err("get bookmark base info error", err)
	}

	err = queryFiles(bookmarkID, &body.Files)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBookmark updates bookmark.
//
// Auto delete files that are no longer used.
// Causes panic when a transaction error occurs.
func UpdateBookmark(bookmarkID int64, body *model.BookmarkPatchBody) error {
	log := logger.New("bookmark_id", bookmarkID)

	cond := make([]string, 0)
	m := data.StructToMap(*body, "Files")
	var mod bool

	for k, v := range m {
		if k == "url" {
			//? Not allowed to set the URL as empty.
			switch s := v.(type) {
			case string, model.MyString:
				if s == "" {
					delete(m, "url")
				} else {
					cond = append(cond, k+" = :"+k)
					mod = true
				}
			default:
				delete(m, "url")
			}
		} else {
			if k != "visits" {
				mod = true
			}
			cond = append(cond, k+" = :"+k)
		}
	}

	//* First check if the bookmark exists
	state, qErr := isBookmarkExists(bookmarkID)
	if qErr == nil {
		if !state {
			log.Info("bookmark does not exist")
			return model.ErrNotExist
		}
	} else {
		return qErr
	}

	if len(cond) > 0 || body.Files != nil {
		tx := db.MustBegin()

		//* Basic information
		if len(cond) > 0 {
			m["id"] = bookmarkID
			modTimeStr := ""
			if mod {
				modTimeStr = "modified_time = datetime('now', 'localtime'), "
			}

			_, err := tx.NamedExec("update bookmarks set "+modTimeStr+strings.Join(cond, ", ")+" where id = :id", m)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update bookmark error", err)
			}
		}

		//* Files
		if body.Files != nil {
			var files = make([]model.MyInt64, 0)
			err := readFiles(tx, bookmarkID, &files)
			if err != nil {
				return nil
			}

			del := data.Defference(files, body.Files)
			pathList, pErr := getFiles(tx, del)
			if pErr != nil {
				return pErr
			}
			if dErr := deleteFiles(tx, del); dErr != nil {
				return dErr
			}

			add := data.Defference(body.Files, files)
			if aErr := updateFiles(tx, bookmarkID, add); aErr != nil {
				return aErr
			}

			//! Delete associated files.
			for _, item := range pathList {
				path := model.UPLOAD_FILES_PATH + "/" + item.String()
				//? delete real file errors will not cause the transaction to rollback.
				data.DeleteFile(path)
			}
		}

		if cErr := tx.Commit(); cErr != nil {
			panic(cErr)
		}
	} else {
		log.Info("no data that needs to be processed")
		return model.ErrDoNothing
	}

	return nil
}

// IncreaseBookmarkVisits increase bookmark visits.
func IncreaseBookmarkVisits(bookmarkID int64) error {
	log := logger.New("bookmark_id", bookmarkID)
	_, err := db.Exec("update bookmarks set visits = visits + 1 where id = ?", bookmarkID)
	if err != nil {
		return log.Err("increase bookmark visits error", err)
	}

	return nil
}

// DeleteBookmark deletes bookmark.
//
// Auto delete associated files.
// Causes panic when a transaction error occurs.
func DeleteBookmark(bookmarkID int64) error {
	log := logger.New("bookmark_id", bookmarkID)

	tx := db.MustBegin()

	var unit struct {
		Name model.MyString `db:"name"`
		URL  model.MyString `db:"url"`
	}

	err := tx.Get(&unit, "select name, url from bookmarks where id = ?", bookmarkID)
	if err == sql.ErrNoRows {
		log.Warn("trying to delete a non-existent bookmark")
		return nil
	} else if err != nil {
		return log.Err("error getting information about bookmark to be deleted", err)
	}
	log.SetMeta("name", unit.Name, "url", unit.URL)

	list := make([]model.MyString, 0)

	err = getFilesByBookmarkID(tx, bookmarkID, &list)
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from bookmarks where id = ?", bookmarkID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("delete bookmark error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	//? delete associated files.
	for _, item := range list {
		path := model.UPLOAD_FILES_PATH + "/" + item.String()
		//? delete real file errors will not cause the transaction to rollback.
		fErr := data.DeleteFile(path)
		if fErr == nil {
			log.Info("auto delete file", "path", path)
		} else {
			log.Warn("auto delete file error", "path", path)
		}
	}

	log.Info("deleted bookmark")

	return nil
}

// BookmarkBatchHandler handles bookmark in batches.
//
// Causes panic when a transaction error occurs.
func BookmarkBatchHandler(body *model.BatchPatchBody) error {
	if len(body.DataSet) == 0 {
		return nil
	}

	tx := db.MustBegin()

	//? File list
	pathList := make([]model.MyString, 0)

	if body.Action == "delete" {
		//? Read file list
		fQeury, fArgs, fErr := sqlx.In("select path from files where bookmark_id in (?)", body.DataSet)
		if fErr != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("in for select files statement error", fErr)
		}

		qErr := tx.Select(&pathList, fQeury, fArgs...)
		if qErr != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("select files error", qErr)
		}

		query, args, err := sqlx.In("delete from bookmarks where id in (?)", body.DataSet)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("in for delete bookmarks statement error", err)
		}
		_, err = tx.Exec(query, args...)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("delete bookmarks error", err)
		}
	} else if body.Action == "setPrivacy" {
		b, ok := body.Payload.(bool)
		if ok {
			log := logger.New("column", "privacy")

			query, args, err := sqlx.In("update bookmarks set modified_time = datetime('now', 'localtime'), privacy = ? where id in (?)", b, body.DataSet)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update bookmarks statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update bookmarks error", err)
			}
		} else {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return model.ErrQueryParamMissing
		}
	} else if body.Action == "setHideInOther" {
		b, ok := body.Payload.(bool)
		if ok {
			log := logger.New("column", "hide_in_other")

			query, args, err := sqlx.In("update bookmarks set modified_time = datetime('now', 'localtime'), hide_in_other = ? where id in (?)", b, body.DataSet)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update bookmarks statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update bookmarks error", err)
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

			var params []any
			str := ""
			if body.Action == "setWeight" {
				str = "update bookmarks set modified_time = datetime('now', 'localtime'), weight = ? where id in (?)"
				params = []any{i, body.DataSet}
			} else {
				str = "update bookmarks set modified_time = datetime('now', 'localtime'), weight = weight " + data.Int16ToStringWithSign(int16(math.Round(i))) + " where id in (?)"
				params = []any{body.DataSet}
			}

			query, args, err := sqlx.In(str, params...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update bookmarks statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update bookmarks error", err)
			}
		} else {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return model.ErrQueryParamMissing
		}
	} else if body.Action == "clearVisits" {
		log := logger.New("column", "visits")

		query, args, err := sqlx.In("update bookmarks set visits = 0 where id in (?)", body.DataSet)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return log.Err("in for update bookmarks statement error", err)
		}
		_, err = tx.Exec(query, args...)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return log.Err("update bookmarks error", err)
		}
	} else if body.Action == "setFolder" {
		i, ok := body.Payload.(float64)
		if ok {
			log := logger.New("column", "folder_id")

			query, args, err := sqlx.In("update bookmarks set modified_time = datetime('now', 'localtime'), folder_id = ? where id in (?)", model.NullInt64(math.Round(i)), body.DataSet)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update bookmarks statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update bookmarks error", err)
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

	if body.Action == "delete" {
		logger.Info("deleted bookmarks in batches", "data", body.DataSet)
	}

	//! Delete files
	for _, item := range pathList {
		path := model.UPLOAD_FILES_PATH + "/" + item.String()
		//? Delete real file errors will not cause the transaction to rollback.
		fErr := data.DeleteFile(path)
		if fErr == nil {
			logger.Info("auto delete file", "path", path)
		} else {
			logger.Warn("auto delete file error", "path", path)
		}
	}

	return nil
}

// ImportBookmarks imports bookmarks and returns number of imports.
//
// Causes panic when a transaction error occurs.
func ImportBookmarks(list *[]model.BookmarkImportItem) (int64, error) {
	useList := make([]model.BookmarkImportItem, 0)
	for _, item := range *list {
		if item.URL != "" {
			useList = append(useList, item)
		}
	}

	importLen := len(useList)
	if importLen == 0 {
		return int64(importLen), nil
	}

	tx := db.MustBegin()

	nstmt, err := tx.PrepareNamed("insert into bookmarks(url, name, description, icon, privacy, weight, hide_in_other, created_time, modified_time) values(:url, :name, :description, :icon, :privacy, :weight, :hide_in_other, datetime('now', 'localtime'), datetime('now', 'localtime'))")
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, logger.Err("prepare insert bookmark statement error", err)
	}
	defer nstmt.Close()

	for _, useItem := range useList {
		_, err := nstmt.Exec(useItem)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return 0, logger.Err("insert bookmark error", err, "bookmark_url", useItem.URL)
		}
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	logger.Info("imported bookmarks", "amount", importLen)

	return int64(importLen), nil
}

// getBookmarkFilterCondition returns bookmark list filter condition.
func getBookmarkFilterCondition(cond *model.BookmarkListCondition) (string, map[string]any) {
	condList := make([]string, 0)
	params := map[string]any{}

	if cond.Privacy != nil {
		condList = append(condList, "(b.privacy = :privacy)")
		params["privacy"] = *cond.Privacy
	}

	if cond.HideInOther != nil {
		condList = append(condList, "(b.hide_in_other = :hide_in_other)")
		params["hide_in_other"] = *cond.HideInOther
	}

	if cond.Folder != nil {
		if *cond.Folder == 0 {
			condList = append(condList, "(b.folder_id is null)")
		} else {
			condList = append(condList, "(b.folder_id = :folder_id)")
			params["folder_id"] = *cond.Folder
		}
	}

	if cond.Name != "" {
		condList = append(condList, `b.name like :name escape '\'`)
		params["name"] = cond.Name.LikeMatchingString()
	}

	if cond.Des != "" {
		condList = append(condList, `b.description like :description escape '\'`)
		params["description"] = cond.Des.LikeMatchingString()
	}

	if cond.URL != "" {
		condList = append(condList, `b.url like :url escape '\'`)
		params["url"] = cond.URL.LikeMatchingString()
	}

	return strings.Join(condList, " and "), params
}

// getBookmarkListOrder returns bookmark list order condition.
//
// Optional order values: created-time | modified-time | visits | folder-weight | weight.
func getBookmarkListOrder(order string) string {
	var o string

	switch order {
	case "modified-time":
		o = "b.modified_time desc"
	case "visits":
		o = "b.visits desc"
	case "folder-weight":
		o = "f.large desc, f.weight desc, b.weight desc"
	case "weight":
		o = "b.weight desc"
	case "created-time":
		fallthrough
	default:
		o = "b.created_time desc"
	}

	return o
}

// isBookmarkExists returns true when the bookmark exists.
func isBookmarkExists(bookmarkID int64) (bool, error) {
	log := logger.New("bookmark_id", bookmarkID)

	var readID model.MyInt64
	err := db.Get(&readID, "select id from bookmarks where id = ?", bookmarkID)
	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	} else {
		return false, log.Err("an error occurred while checking whether the bookmark exists", err)
	}
}
