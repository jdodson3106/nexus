package ar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRecord struct {
	BaseRecord
	Name  string
	Email string
	Age   int
}

func TestBaseQuery(t *testing.T) {
	expected := &TestRecord{
		BaseRecord: BaseRecord{1},
		Name:       "Justin Dodson",
		Email:      "test@mail.com",
		Age:        35,
	}
	tr := TestRecord{Name: "Justin Dodson", Email: "test@mail.com", Age: 35}
	saved, err := tr.Create(&tr)
	assert.NoError(t, err)
	assert.Equal(t, expected, saved)
}
