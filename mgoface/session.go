package mgoface

import (
	"time"

	"github.com/globalsign/mgo"
)

type Session interface {
	DB(name string) Database
	Copy() Session
	Clone() Session
	New() Session
	Close()

	Safe() *mgo.Safe
	SetSafe(safe *mgo.Safe)
	EnsureSafe(safe *mgo.Safe)
	DatabaseNames() (names []string, err error)
}

type Database interface {
	C(name string) Collection
	With(s Session) Database

	DropDatabase() error
	FindRef(ref *mgo.DBRef) Query
	CollectionNames() (names []string, err error)
}

type Collection interface {
	With(s Session) Collection
	DropCollection() error

	Find(query interface{}) Query
	FindId(id interface{}) Query
	Insert(docs ...interface{}) error
	Update(selector interface{}, update interface{}) error
	UpdateId(id interface{}, update interface{}) error
	UpdateWithArrayFilters(selector, update, arrayFilters interface{}, multi bool) (*mgo.ChangeInfo, error)
	UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	Remove(selector interface{}) error
	RemoveId(id interface{}) error
	RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error)

	Create(info *mgo.CollectionInfo) error
}

type Query interface {
	Batch(n int) Query
	Prefetch(p float64) Query
	Skip(n int) Query
	Limit(n int) Query
	Select(selector interface{}) Query
	Sort(fields ...string) Query
	Min(min interface{}) Query
	Max(max interface{}) Query

	Explain(result interface{}) error
	Hint(indexKey ...string) Query
	SetMaxScan(n int) Query
	SetMaxTime(d time.Duration) Query
	Snapshot() Query
	Comment(comment string) Query
	LogReplay() Query

	One(result interface{}) error
	All(result interface{}) error
	Count() (n int, err error)
	Apply(change Change, result interface{}) (info *mgo.ChangeInfo, err error)
}

type Iter interface {
	Err() error
	Close() error
	Done() bool
	Timeout() bool
	Next(result interface{}) bool
	All(result interface{}) error
}
