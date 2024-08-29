package db

import (
	"fmt"
	"reflect"
)

const (
	SELECT     = "select"
	SELECT_ALL = "SELECT_ALL"
	INSERT     = "INSERT"
	UPDATE     = "UPDATE"
	DELETE     = "DELETE"
)

type IModel interface {
	Select(string) IModel

	SelectAll(string) IModel

	Insert(string) IModel

	Update(string) IModel

	Delete(string) IModel

	Where(string) IModel

	Limit(int) IModel

	OrderBy(string) IModel

	OrderByAsc(string) IModel

	Exec() (*interface{}, error)
}

type QueryBuilderError struct {
	ErrorMessage string
}

func (q *QueryBuilderError) Error() string {
	return q.ErrorMessage
}

type DbModel struct {
	queryType    string
	queryStarted bool
	queryString  string
	conn         *DB
	errors       []QueryBuilderError
}

func (db *DbModel) Select(string) IModel {
	if db.checkQueryStarted("Select()") {
		return db
	}
	db.queryType = SELECT

	// TODO: Start building the SELECT statement
	return db
}

func (db *DbModel) SelectAll(string) IModel {
	if db.checkQueryStarted("SelectAll()") {
		return db
	}
	db.queryType = SELECT_ALL

	return db
}

func (db *DbModel) Insert(string) IModel {
	if db.checkQueryStarted("Insert") {
		return db
	}
	db.queryType = INSERT

	return db
}

func (db *DbModel) Update(string) IModel {
	if db.checkQueryStarted("Update()") {
		return db
	}
	db.queryType = UPDATE

	return db
}

func (db *DbModel) Delete(string) IModel {
	if db.checkQueryStarted("Delete()") {
		return db
	}
	db.queryType = DELETE

	return db
}

func (db *DbModel) Where(string) IModel { return nil }

func (db *DbModel) Limit(int) IModel { return nil }

func (db *DbModel) OrderBy(string) IModel { return nil }

func (db *DbModel) OrderByAsc(string) IModel { return nil }

func (db *DbModel) Exec() (*interface{}, error) {

	if !db.queryStarted {
		return nil, &QueryBuilderError{ErrorMessage: "cannot call Exec() on uninitialized query"}
	}

	if len(db.errors) > 0 {
		return nil, &db.errors[0]
	}

	// TODO: Build the actual interface to run the sql.DB.QueryRow(s) against
	return nil, nil
}

func (db *DbModel) checkQueryStarted(action string) bool {
	if db.queryStarted {
		t := reflect.TypeOf(db).Name()
		msg := fmt.Sprintf("query already started when calling %s on %s", action, t)
		db.errors = append(db.errors, QueryBuilderError{ErrorMessage: msg})
		return true
	}
	return false
}
