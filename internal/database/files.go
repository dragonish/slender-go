package database

import (
	"database/sql"
	"slender/internal/data"
	"slender/internal/logger"
	"slender/internal/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

// AddFile adds file record.
func AddFile(file string, bookmarkID model.NullInt64) (int64, error) {
	log := logger.New("file", file)

	tx := db.MustBegin()

	res, err := tx.Exec("insert into files(path, bookmark_id) values(?, ?)", file, bookmarkID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("add file record error", err)
	}

	fileID, err := res.LastInsertId()
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("read the added file record id error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	return fileID, nil
}

// DeleteFile deletes file.
//
// Causes panic when a transaction error occurs.
func DeleteFile(fileID int64, force bool) error {
	log := logger.New("file_id", fileID)

	tx := db.MustBegin()

	var fileInfo model.FileInfo
	err := getFile(tx, fileID, &fileInfo)
	if err != nil {
		return err
	}

	if fileInfo.BookmarkID != 0 && !force {
		//? nothing is actually deleted.
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return nil
	}

	_, err = tx.Exec("delete from files where id = ?", fileID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("delete file error", err)
	}

	path := model.UPLOAD_FILES_PATH + "/" + fileInfo.Path.String()
	dErr := data.DeleteFile(path)
	if dErr != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("delete real file error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	return nil
}

// GetFileList gets file list.
func GetFileList(cond *model.FileListCondition, body *model.FileListData) error {
	body.PageSize = cond.Size
	body.PageNo = cond.Page

	filter, params := getFileFilterCondition(cond)
	if filter != "" {
		filter = " where " + filter
	}

	gStmt, err := db.PrepareNamed("select count(*) from files f" + filter)
	if err != nil {
		return logger.Err("prepared file list count query statement error", err)
	}
	defer gStmt.Close()
	err = gStmt.Get(&body.Total, params)
	if err != nil {
		return logger.Err("get file list count error", err)
	} else if body.Total == 0 {
		return nil
	}

	qStmt, err := db.PrepareNamed("select f.id, '/assets/uploads/' || f.path as path, f.bookmark_id is not null as used from files f " + filter + " order by f.id desc limit :offset,:size")
	if err != nil {
		return logger.Err("prepared file list query statement error", err)
	}
	defer qStmt.Close()

	params["offset"] = cond.Size * (cond.Page - 1)
	params["size"] = cond.Size

	err = qStmt.Select(&body.List, params)
	if err != nil {
		return logger.Err("get file list error", err)
	}

	return nil
}

// RemoveUnusedFiles removes all unused files.
//
// Causes panic when a transaction error occurs.
func RemoveUnusedFiles() error {
	var list = make([]model.MyString, 0)

	tx := db.MustBegin()

	err := tx.Select(&list, "select path from files where bookmark_id is null")
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return logger.Err("select used files error", err)
	}

	if len(list) == 0 {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return nil
	}

	_, err = tx.Exec("delete from files where bookmark_id is null")
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return logger.Err("delete all unused files error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	for _, item := range list {
		path := model.UPLOAD_FILES_PATH + "/" + item.String()
		//? delete real file errors will not cause the transaction to rollback.
		data.DeleteFile(path)
	}

	return nil
}

// getFile returns the file.
// (transaction)
func getFile(tx *sqlx.Tx, fileID int64, info *model.FileInfo) error {
	log := logger.New("file_id", fileID)

	err := tx.Get(info, "select path, bookmark_id from files where id = ?", fileID)
	if err == sql.ErrNoRows {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return model.ErrNotExist
	} else if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("get bookmark id corresponding to the file error", err)
	}

	return nil
}

// getFilesByBookmarkID gets bookmark's file list.
// (transaction)
func getFilesByBookmarkID(tx *sqlx.Tx, bookmarkID int64, list *[]model.MyString) error {
	log := logger.New("bookmark_id", bookmarkID)

	err := tx.Select(list, "select path from files where bookmark_id = ?", bookmarkID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("get file list of the bookmark error", err)
	}

	return nil
}

// queryFiles query bookmark's file list.
func queryFiles(bookmarkID int64, list *[]model.FileBaseData) error {
	log := logger.New("bookmark_id", bookmarkID)

	err := db.Select(list, "select id, '/assets/uploads/' || path as path from files where bookmark_id = ?", bookmarkID)
	if err != nil {
		return log.Err("query the file list of the bookmark error", err)
	}

	return nil
}

// updateFile updates file.
// (transaction)
func updateFile(tx *sqlx.Tx, fileID model.MyInt64, bookmarkID model.MyInt64) error {
	log := logger.New("file_id", fileID, "bookmark_id", bookmarkID)

	_, err := tx.Exec("update files set bookmark_id = ? where id = ?", bookmarkID, fileID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("udpate file error", err)
	}

	return nil
}

// readFiles reads file id list.
// (transaction)
func readFiles(tx *sqlx.Tx, bookmarkID int64, list *[]model.MyInt64) error {
	log := logger.New("bookmark_id", bookmarkID)

	err := tx.Select(list, "select id from files where bookmark_id = ?", bookmarkID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("read file id list error", err)
	}

	return nil
}

// getFiles gets file path list by file id list.
// (transaction)
func getFiles(tx *sqlx.Tx, idList []model.MyInt64) ([]model.MyString, error) {
	pathList := make([]model.MyString, 0)

	if len(idList) == 0 {
		return pathList, nil
	}

	query, args, err := sqlx.In("select path from files where id in (?)", idList)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return pathList, logger.Err("in for select files statement error", err)
	}

	err = tx.Select(&pathList, query, args...)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return pathList, logger.Err("select files error", err)
	}

	return pathList, nil
}

// deleteFiles delete files.
// (transaction)
func deleteFiles(tx *sqlx.Tx, list []model.MyInt64) error {
	if len(list) == 0 {
		return nil
	}

	query, args, err := sqlx.In("delete from files where id in (?)", list)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return logger.Err("in for delete files statement error", err)
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return logger.Err("delete files error", err)
	}

	return nil
}

// updateFiles updates files.
// (transaction)
func updateFiles(tx *sqlx.Tx, bookmarkID int64, list []model.MyInt64) error {
	if len(list) == 0 {
		return nil
	}

	log := logger.New("bookmark_id", bookmarkID)

	query, args, err := sqlx.In("update files set bookmark_id = ? where id in (?)", bookmarkID, list)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("in for update files statement error", err)
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return log.Err("update files error", err)
	}

	return nil
}

// getFileFilterCondition returns file list filter condition.
func getFileFilterCondition(cond *model.FileListCondition) (string, map[string]interface{}) {
	condList := make([]string, 0)
	params := map[string]interface{}{}

	if cond.Use != nil {
		if *cond.Use {
			condList = append(condList, "(f.bookmark_id > 0)")
		} else {
			condList = append(condList, "(f.bookmark_id is null)")
		}
	}

	if cond.Path != "" {
		condList = append(condList, `('/assets/uploads/' || f.path) like :path escape '\'`)
		params["path"] = cond.Path.LikeMatchingString()
	}

	return strings.Join(condList, " and "), params
}
