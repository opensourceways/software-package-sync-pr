package postgresql

import (
	"errors"
)

var (
	errRowExists   = errors.New("row exists")
	errRowNotFound = errors.New("row not found")
)

type dbTable struct {
	name string
}

func NewDBTable(name string) dbTable {
	return dbTable{name: name}
}

func (t dbTable) FirstOrCreate(filter, result interface{}) error {
	query := db.Table(t.name).Where(filter).FirstOrCreate(result)

	if err := query.Error; err != nil {
		return err
	}

	if query.RowsAffected == 0 {
		return errRowExists
	}

	return nil
}

func (t dbTable) UpdateRecord(filter, update interface{}) (err error) {
	query := db.Table(t.name).Where(filter).Updates(update)
	if err = query.Error; err != nil {
		return
	}

	if query.RowsAffected == 0 {
		err = errRowNotFound
	}

	return
}

func (t dbTable) DeleteRecord(filter interface{}) (err error) {
	query := db.Table(t.name).Where(filter).Delete(nil)
	if err = query.Error; err != nil {
		return
	}

	if query.RowsAffected == 0 {
		err = errRowNotFound
	}

	return
}

func (t dbTable) IsRowNotFound(err error) bool {
	return errors.Is(err, errRowNotFound)
}

func (t dbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}
