package ar

import (
	"fmt"
	"reflect"
)

type QueryParams struct {
	Column string
	Value  interface{}
}

type ActiveRecord interface {
	Get(params ...QueryParams) (ActiveRecord, error)
	Create(ActiveRecord) (ActiveRecord, error)
	Update() (ActiveRecord, error)
	Delete() error
}

// BaseRecord gets embedded in all ActiveRecord models
// and provides base fields
type BaseRecord struct {
	Id int
}

func (b *BaseRecord) updateCaller(caller ActiveRecord) ActiveRecord {
	el := reflect.ValueOf(caller)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	// set the baserecord
	el.FieldByName("BaseRecord").Set(reflect.ValueOf(b).Elem())
	addr := el.Addr().Interface()

	if c, ok := addr.(ActiveRecord); ok {
		fmt.Print("all good\n")
		return c
	}
	fmt.Printf("\n\n##############\n\tERROR :: %v\n##############\n\n", caller)
	return nil
}

func (b *BaseRecord) Create(record ActiveRecord) (ActiveRecord, error) {
	b.Id = 1
	return record, nil
}

func (b *BaseRecord) Get(params ...QueryParams) (ActiveRecord, error) {
	return nil, nil
}

func (b *BaseRecord) Update() (ActiveRecord, error) {
	return nil, nil
}

func (b *BaseRecord) Delete() error {
	return nil
}

/*
	type Model struct {

	}


*/
