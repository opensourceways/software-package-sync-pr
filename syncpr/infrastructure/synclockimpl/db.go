package synclockimpl

type dbClient interface {
	FirstOrCreate(filter, result interface{}) error
	UpdateRecord(filter, update interface{}) error
	DeleteRecord(filter interface{}) error
	IsRowNotFound(err error) bool
	IsRowExists(err error) bool
}
