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

// AddSearchEngine adds new search engine item.
//
// Causes panic when a transaction error occurs.
func AddSearchEngine(body *model.SearchEnginePostBody) (int64, error) {
	log := logger.New("name", body.Name, "url", body.URL)

	tx := db.MustBegin()

	res, err := tx.NamedExec("insert into search_engines(name, method, url, body, icon, weight, created_time, modified_time) values(:name, :method, :url, :body, :icon, :weight, datetime('now', 'localtime'), datetime('now', 'localtime'))", body)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("add serach engine error", err)
	}

	searchEngineID, err := res.LastInsertId()
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return 0, log.Err("read the created search engine id error", err)
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	log.Info("add new search engine", "search_engine_id", searchEngineID)

	return searchEngineID, nil
}

// GetSearchEngines gets search engine list.
func GetSearchEngines(cond *model.SearchEngineListCondition, body *model.SearchEngineListData) error {
	body.PageSize = cond.Size
	body.PageNo = cond.Page

	filter, params := getSearchEngineFilterCondition(cond)
	if filter != "" {
		filter = " where " + filter
	}

	gStmt, err := db.PrepareNamed("select count(*) from search_engines s" + filter)
	if err != nil {
		return logger.Err("parepared search engine list count query statement error", err)
	}
	defer gStmt.Close()
	err = gStmt.Get(&body.Total, params)
	if err != nil {
		return logger.Err("get search engine list count error", err)
	} else if body.Total == 0 {
		return nil
	}

	o := getSearchEngineListOrder(cond.Order)

	qStmt, err := db.PrepareNamed("select s.id, s.name, s.method, s.url, s.body, s.icon, s.weight, s.created_time, s.modified_time from search_engines s " + filter + " order by " + o + " limit :offset,:size")
	if err != nil {
		return logger.Err("prepared search engine list query statement error", err)
	}
	defer qStmt.Close()

	params["offset"] = cond.Size * (cond.Page - 1)
	params["size"] = cond.Size

	err = qStmt.Select(&body.List, params)
	if err != nil {
		return logger.Err("get search engine list error", err)
	}
	return nil
}

// GetHomeSearchEngines gets search engine list used by the homepage.
func GetHomeSearchEngines(list *[]model.HomeSearchEngineListItem) error {
	err := db.Select(list, "select s.id, s.name, s.method, s.url, s.body, s.icon from search_engines s order by s.weight desc, s.id")
	if err != nil {
		return logger.Err("get search engine list used by the homepage error", err)
	}

	return nil
}

// GetSearchEngine gets search engine item data.
func GetSearchEngine(searchEngineID int64, body *model.SearchEngineBaseData) error {
	log := logger.New("search_engine_id", searchEngineID)

	err := db.Get(body, "select * from search_engines where id = ?", searchEngineID)
	if err == sql.ErrNoRows {
		log.Info("search engine does not exist")
		return model.ErrNotExist
	} else if err != nil {
		return log.Err("get search engine info error", err)
	}

	return nil
}

// UpdateSearchEngine updates search engine.
//
// Causes panic when a transaction error occurs.
func UpdateSearchEngine(searchEngineID int64, body *model.SearchEnginePatchBody) error {
	log := logger.New("search_engine_id", searchEngineID)

	cond := make([]string, 0)
	m := data.StructToMap(*body)

	for k, v := range m {
		if k == "name" || k == "url" {
			//? Not allowed to set the name|url as empty.
			switch s := v.(type) {
			case string, model.MyString:
				if s == "" {
					delete(m, k)
				} else {
					cond = append(cond, k+" = :"+k)
				}
			default:
				delete(m, k)
			}
		} else {
			cond = append(cond, k+" = :"+k)
		}
	}

	//* First check if the search engine exists
	state, qErr := isSearchEngineExists(searchEngineID)
	if qErr == nil {
		if !state {
			log.Info("search engine not exist")
			return model.ErrNotExist
		}
	} else {
		return qErr
	}

	if len(cond) > 0 {
		tx := db.MustBegin()

		m["id"] = searchEngineID
		_, err := tx.NamedExec("update search_engines set modified_time = datetime('now', 'localtime'), "+strings.Join(cond, ", ")+" where id = :id", m)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return log.Err("update search engine error", err)
		}

		if cErr := tx.Commit(); cErr != nil {
			panic(cErr)
		}
	} else {
		log.Info("no data that needs to be processed")
	}

	return nil
}

// DeleteSearchEngine deletes search engine.
func DeleteSearchEngine(searchEngineID int64) error {
	log := logger.New("search_engine_id", searchEngineID)

	var unit struct {
		Name model.MyString `db:"name"`
		URL  model.MyString `db:"url"`
	}
	err := db.Get(&unit, "select name, url from search_engines where id = ?", searchEngineID)
	if err == sql.ErrNoRows {
		log.Warn("trying to delete a non-existent search engine")
		return nil
	} else if err != nil {
		return log.Err("error getting information about search engine to be deleted", err)
	}
	log.SetMeta("name", unit.Name, "url", unit.URL)

	_, err = db.Exec("delete from search_engines where id = ?", searchEngineID)
	if err != nil {
		return log.Err("delete search engine error", err)
	}

	log.Info("deleted search engine")

	return nil
}

func SearchEngineBatchHandler(body *model.BatchPatchBody) error {
	if len(body.DataSet) == 0 {
		return nil
	}

	tx := db.MustBegin()

	switch body.Action {
	case "delete":
		query, args, err := sqlx.In("delete from search_engines where id in (?)", body.DataSet)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("in for delete search engines statement error", err)
		}
		_, err = tx.Exec(query, args...)
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return logger.Err("delete search engines error", err)
		}
	case "setWeight", "incWeight":
		i, ok := body.Payload.(float64)
		if ok {
			log := logger.New("column", "weight")

			var params []any
			str := ""
			if body.Action == "setWeight" {
				str = "update search_engines set modified_time = datetime('now', 'localtime'), weight = ? where id in (?)"
				params = []any{i, body.DataSet}
			} else {
				str = "update search_engines set modified_time = datetime('now', 'localtime'), weight = weight " + data.Int16ToStringWithSign(int16(math.Round(i))) + " where id in (?)"
				params = []any{body.DataSet}
			}

			query, args, err := sqlx.In(str, params...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("in for update search engines statement error", err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				return log.Err("update search engines error", err)
			}
		} else {
			if rErr := tx.Rollback(); rErr != nil {
				panic(rErr)
			}
			return model.ErrQueryParamMissing
		}
	default:
		if rErr := tx.Rollback(); rErr != nil {
			panic(rErr)
		}
		return model.ErrDoNothing
	}

	if cErr := tx.Commit(); cErr != nil {
		panic(cErr)
	}

	if body.Action == "delete" {
		logger.Info("deleted search engines in batches", "data", body.DataSet)
	}

	return nil
}

// getSearchEngineFilterCondition returns search engine list filter condtion.
func getSearchEngineFilterCondition(cond *model.SearchEngineListCondition) (string, map[string]any) {
	condList := make([]string, 0)
	params := map[string]any{}

	if cond.Method != nil {
		condList = append(condList, ("(s.method = :method)"))
		params["method"] = *cond.Method
	}

	if cond.Name != "" {
		condList = append(condList, `s.name like :name escape '\'`)
		params["name"] = cond.Name.LikeMatchingString()
	}

	if cond.URL != "" {
		condList = append(condList, `s.url like :url escape '\'`)
		params["url"] = cond.URL.LikeMatchingString()
	}

	return strings.Join(condList, " and "), params
}

// getSearchEngineListOrder returns search engine list order condition.
//
// Optional order values: created-time | modified-time | weight.
func getSearchEngineListOrder(order string) string {
	var o string

	switch order {
	case "modified-time":
		o = "s.modified_time desc"
	case "weight":
		o = "s.weight desc"
	case "created-time":
		fallthrough
	default:
		o = "s.created_time desc"
	}

	return o
}

// isSearchEngineExists returns true when search engine exists.
func isSearchEngineExists(searchEngineID int64) (bool, error) {
	log := logger.New("search_engine_id", searchEngineID)

	var readID model.MyInt64
	err := db.Get(&readID, "select id from search_engines where id = ?", searchEngineID)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, log.Err("an error occurred while checking whether the search engine exists", err)
	}
}
