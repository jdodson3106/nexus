package db

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	SELECT     = "select"
	SELECT_ALL = "SELECT_ALL"
	INSERT     = "INSERT"
	UPDATE     = "UPDATE"
	DELETE     = "DELETE"
)

type IModel interface {
	Select(...string) IModel

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
	varCount     int
	errors       []error
}

func (db DbModel) getBaseFields() []string {
	return []string{
		"queryType",
		"queryStarted",
		"queryString",
		"conn",
		"carCount",
		"errors",
	}
}

func (db DbModel) isBaseField(prop string) bool {
	for _, f := range db.getBaseFields() {
		if prop == f {
			return true
		}
	}
	return false
}

func (db DbModel) isValidField(prop string) bool {
	t := reflect.TypeOf(db)
	for i := 0; i < t.NumField(); i++ {
		if prop == t.Field(i).Name {
			return true
		}
	}
	return false
}

func (db *DbModel) Select(props ...string) IModel {
	if db.checkQueryStarted("Select()") {
		return db
	}
	db.queryType = SELECT
	var qBuilder strings.Builder
	qBuilder.Write([]byte("SELECT "))

	propStr, count, err := db.buildPropertyListString(props...)
	if err != nil {
		db.errors = append(db.errors, err)
		return db
	}

	qBuilder.Write(propStr)
	db.varCount += count
	db.queryString = qBuilder.String()
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
		return nil, db.errors[0]
	}

	// TODO: Build the actual interface to run the sql.DB.QueryRow(s) against
	return nil, nil
}

func (db *DbModel) checkQueryStarted(action string) bool {
	if db.queryStarted {
		t := reflect.TypeOf(db).Name()
		msg := fmt.Sprintf("query already started when calling %s on %s", action, t)
		db.errors = append(db.errors, &QueryBuilderError{ErrorMessage: msg})
		return true
	}
	return false
}

func (db DbModel) buildPropertyListString(props ...string) (b []byte, count int, err error) {
	var propCount int
	buf := make([]byte, 256)

	if len(props) == 0 {
		propCount = 0
	} else {
		propCount = len(props)
		for _, prop := range props {
			if db.isValidField(prop) {
				buf = append(buf, []byte(fmt.Sprintf("%s, ", prop))...)
			} else {
				msg := fmt.Sprintf("invalid property provided \"%s\"", prop)
				return nil, 0, &QueryBuilderError{ErrorMessage: msg}
			}
		}
	}

	return buf, propCount, nil
}

func (db DbModel) buildAllPropsString() ([]byte, int) {
	basePropCount := len(db.getBaseFields())
	basePropsSeenCount := 0
	propCount := 0
	buf := make([]byte, 0)

	t := reflect.TypeOf(db)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if !db.isBaseField(name) {
			buf = append(buf, name...)
			propCount++
		} else {
			basePropsSeenCount++
		}

		// we've seen all the base model props, so just use the index count for the end
		baseDiff := basePropCount - basePropsSeenCount
		toEnd := (t.NumField() - i) - baseDiff

		if baseDiff == 0 && t.NumField()-i > 0 {
			// not at the end yet...
		}

		//
		if baseDiff-(t.NumField()-i) == 0 {
			// we are at the end but still have base fields to parse. so jump from loop
		}

		if (baseDiff == 0 || toEnd == 0) && i < t.NumField()-1 {
			buf = append(buf, ", "...) // add the comma separator
		}
	}

	return buf, propCount
}
