package database

import (
	"database/sql"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"slender/internal/data"
	"slender/internal/logger"
	"slender/internal/model"
)

var (
	db   *sqlx.DB
	once sync.Once
)

// Load handles connect and initialize the database.
//
// journal_mode=WAL and synchronous=NORMAL when wal is true.
//
// If an error occurs during the connection, the program is aborted.
// The filename parameter does not take an extension name.
func Load(filename string, wal bool) {
	dbFile := getDBFilePath(filename)

	fileLog := logger.New("db_file", dbFile)
	if data.IsPathExists(dbFile) {
		fileLog.Debug("database file exists")
	} else {
		fileLog.Info("create database file")
		file, err := os.Create(dbFile)
		if err != nil {
			fileLog.Fatal("error craeting database file", err)
		}
		file.Close()
	}

	once.Do(func() {
		log := logger.New("db_file", dbFile)

		//* enable foreign key feature.
		iDb, err := sqlx.Open("sqlite3", dbFile+"?_fk=on")
		if err != nil {
			log.Fatal("database initial connection error", err)
		} else {
			log.Debug("initial connect database")
		}
		defer iDb.Close()

		_, err = iDb.Exec(`create table if not exists users(
			username varchar(50) primary key,
			salt varchar(64) not null,
			secret varchar(16) not null
		);`)
		createTableFatal(err, "users")

		_, err = iDb.Exec(`create table if not exists folders(
			id integer primary key autoincrement,
			name varchar(100) not null,
			description text,
			large bool not null default false,
			privacy bool not null default false,
			weight smallint not null default 0 check(weight >= -32768 and weight <= 32767),
			created_time datetime not null,
			modified_time datetime not null
		);`)
		createTableFatal(err, "folders")

		_, err = iDb.Exec(`create table if not exists bookmarks(
			id integer primary key autoincrement,
			url text not null,
			name varchar(255) not null,
			description text,
			icon text,
			privacy bool not null default false,
			weight smallint not null default 0 check(weight >= -32768 and weight <= 32767),
			created_time datetime not null,
			modified_time datetime not null,
			visits int unsigned not null default 0,
			folder_id integer,
			foreign key(folder_id) references folders(id) on delete set null on update cascade
		);`)
		createTableFatal(err, "bookmarks")

		_, err = iDb.Exec("create index if not exists idx_bookmarks_created_time on bookmarks(created_time)")
		createIndexFatal(err, "bookmarks", "idx_bookmarks_created_time")

		_, err = iDb.Exec(`create table if not exists search_engines(
			id integer primary key autoincrement,
			name varchar(100) not null,
			method char(10) not null,
			url text not null,
			body text,
			icon text,
			weight smallint not null default 0 check(weight >= -32768 and weight <= 32767),
			created_time datetime not null,
			modified_time datetime not null
		);`)
		createTableFatal(err, "search_engines")

		_, err = iDb.Exec(`create table if not exists files(
			id integer primary key autoincrement,
			path text,
			bookmark_id integer,
			foreign key(bookmark_id) references bookmarks(id) on delete cascade on update cascade
		);`)
		createTableFatal(err, "files")

		_, err = iDb.Exec(`create table if not exists logins(
			login_id varchar(36) primary key,
			login_time datetime not null,
			ip varchar(45) not null,
			ua text,
			is_admin bool not null default false,
			max_age smallint unsigned not null default 0,
			active bool default null
		);`)
		createTableFatal(err, "logins")

		_, err = iDb.Exec("create index if not exists idx_logins_login_time on logins(login_time)")
		createIndexFatal(err, "logins", "idx_logins_login_time")

		//? add new column max_age to logins table
		ageMeta := []any{"table", "logins", "column", "max_age"}
		var ageCol model.MyString
		ageErr := iDb.Get(&ageCol, "select name from pragma_table_info(?) where name = ?", "logins", "max_age")
		if ageErr == sql.ErrNoRows {
			log.Debug("add new column to table", ageMeta...)
			tx := iDb.MustBegin()

			_, ageErr = tx.Exec("alter table logins add column max_age smallint unsigned not null default 0")
			if ageErr != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				log.Fatal("failed to add column to table", ageErr, ageMeta...)
			}

			if cErr := tx.Commit(); cErr != nil {
				panic(cErr)
			}
		} else if ageErr != nil {
			log.Fatal("failed to read column info from table", ageErr, ageMeta...)
		}

		//? add new column active to logins table
		activeMeta := []any{"table", "logins", "column", "active"}
		var activeCol model.MyString
		activeErr := iDb.Get(&activeCol, "select name from pragma_table_info(?) where name = ?", "logins", "active")
		if activeErr == sql.ErrNoRows {
			log.Debug("add new column to table", activeMeta...)
			tx := iDb.MustBegin()

			_, activeErr = tx.Exec("alter table logins add column active bool default null")
			if activeErr != nil {
				if rErr := tx.Rollback(); rErr != nil {
					panic(rErr)
				}
				log.Fatal("failed to add column to table", activeErr, activeMeta...)
			}

			if cErr := tx.Commit(); cErr != nil {
				panic(cErr)
			}
		} else if activeErr != nil {
			log.Fatal("failed to read column info from table", activeErr, activeMeta...)
		}
	})

	connect(dbFile, wal)
}

func connect(dbFile string, wal bool) {
	var err error
	log := logger.New("db_file", dbFile)

	if wal {
		db, err = sqlx.Open("sqlite3", dbFile+"?_fk=on&_journal=WAL&_sync=NORMAL")
	} else {
		//* enable foreign key feature.
		db, err = sqlx.Open("sqlite3", dbFile+"?_fk=on")
	}

	if err != nil {
		log.Fatal("database connection error", err)
	} else {
		log.Info("connect database")
	}

	//* register
	rErr := Register()
	if rErr != nil {
		log.Fatal("register error", rErr)
	}
}

func createTableFatal(err error, tableName string) {
	if err != nil {
		logger.Fatal("table creation error", err, "table", tableName)
	}
}

func createIndexFatal(err error, tableName, indexName string) {
	if err != nil {
		logger.Fatal("index creation error", err, "table", tableName, "index", indexName)
	}
}

// getDBFilePath returns database file full path.
//
// The filename parameter does not take an extension name.
func getDBFilePath(filename string) string {
	return model.DATA_DIR + "/" + filename + ".db"
}
