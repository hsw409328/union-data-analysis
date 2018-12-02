package driver

import (
	"database/sql"
	"github.com/hsw409328/gofunc/go_hlog"
	_ "github.com/mattn/go-sqlite3"
)

var (
	SQLiteDriverAnalysis = NewSQLite("./db/2018-12-02-003925.db")
	SQLiteDriverWeb      = NewSQLite("./db/ad-data-analysis.db")
	lg                   *go_hlog.Logger
	sqlLog               bool
)

func init() {
	lg = go_hlog.GetInstance("")
	sqlLog = true
}

type SQLite struct {
	Source string
	DB     *sql.DB
}

func NewSQLite(source string) *SQLite {
	return &SQLite{Source: source}
}

func (ctx *SQLite) Init() {
	if ctx.DB != nil {
		return
	}
	db, err := sql.Open("sqlite3", ctx.Source)
	if err != nil {
		lg.Error(err.Error())
	}
	ctx.DB = db
}

func (ctx *SQLite) GetOne(sqlStmt string) *sql.Row {
	ctx.Init()
	if sqlLog {
		lg.Debug(sqlStmt)
	}
	r := ctx.DB.QueryRow(sqlStmt)
	return r
}

func (ctx *SQLite) GetAll(sqlStmt string) (*sql.Rows, error) {
	ctx.Init()
	if sqlLog {
		lg.Debug(sqlStmt)
	}
	r, err := ctx.DB.Query(sqlStmt)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}
	return r, nil
}

func (ctx *SQLite) Insert(sqlStmt string, v ...interface{}) (int64, error) {
	ctx.Init()
	if sqlLog {
		lg.Debug(sqlStmt)
	}
	stmt, err := ctx.DB.Prepare(sqlStmt)
	if err != nil {
		lg.Error(err)
		return 0, err
	}
	result, err := stmt.Exec(v...)
	if err != nil {
		lg.Error(err.Error())
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		lg.Error(err)
		return 0, err
	}
	return lastID, err
}
