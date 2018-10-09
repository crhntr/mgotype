package mgoface

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Session interface {
	BuildInfo() (info mgo.BuildInfo, err error)
	Clone() *mgo.Session
	Close()
	Copy() *mgo.Session
	DB(name string) *Database
	DatabaseNames() (names []string, err error)
	EnsureSafe(safe *mgo.Safe)
	FindRef(ref *mgo.DBRef) *mgo.Query
	Fsync(async bool) error
	FsyncLock() error
	FsyncUnlock() error
	LiveServers() (addrs []string)
	Login(cred *mgo.Credential) error
	LogoutAll()
	Mode() mgo.Mode
	New() *mgo.Session
	Ping() error
	Refresh()
	ResetIndexCache()
	Run(cmd interface{}, result interface{}) error
	Safe() (safe *mgo.Safe)
	SelectServers(tags ...bson.D)
	SetBatch(n int)
	SetBypassValidation(bypass bool)
	SetCursorTimeout(d time.Duration)
	SetMode(consistency mgo.Mode, refresh bool)
	SetPoolLimit(limit int)
	SetPoolTimeout(timeout time.Duration)
	SetPrefetch(p float64)
	SetSafe(safe *mgo.Safe)
	SetSocketTimeout(d time.Duration)
	SetSyncTimeout(d time.Duration)
}

type Database interface {
	AddUser(username, password string, readOnly bool) error
	C(name string) *mgo.Collection
	CollectionNames() (names []string, err error)
	CreateView(view string, source string, pipeline interface{}, collation *mgo.Collation) error
	DropDatabase() error
	FindRef(ref *mgo.DBRef) *mgo.Query
	// GridFS(prefix string) *mgo.GridFS
	Login(user, pass string) error
	Logout()
	RemoveUser(user string) error
	Run(cmd interface{}, result interface{}) error
	UpsertUser(user *mgo.User) error
	With(s *mgo.Session) *Database
}

type Collection interface {
	Bulk() *mgo.Bulk
	Count() (n int, err error)
	Create(info *mgo.CollectionInfo) error
	DropAllIndexes() error
	DropCollection() error
	DropIndex(key ...string) error
	DropIndexName(name string) error
	EnsureIndex(index mgo.Index) error
	EnsureIndexKey(key ...string) error
	Find(query interface{}) *mgo.Query
	FindId(id interface{}) *mgo.Query
	Indexes() (indexes []mgo.Index, err error)
	Insert(docs ...interface{}) error
	NewIter(session *Session, firstBatch []bson.Raw, cursorId int64, err error) *mgo.Iter
	Pipe(pipeline interface{}) *mgo.Pipe
	Remove(selector interface{}) error
	RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error)
	RemoveId(id interface{}) error
	Repair() *mgo.Iter
	Update(selector interface{}, update interface{}) error
	UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	UpdateId(id interface{}, update interface{}) error
	UpdateWithArrayFilters(selector, update, arrayFilters interface{}, multi bool) (*mgo.ChangeInfo, error)
	Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	Watch(pipeline interface{}, options mgo.ChangeStreamOptions) (*mgo.ChangeStream, error)
	With(s *Session) *mgo.Collection
}

type Query interface {
	All(result interface{}) error
	Apply(change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error)
	Batch(n int) *mgo.Query
	Collation(collation *mgo.Collation) *mgo.Query
	Comment(comment string) *mgo.Query
	Count() (n int, err error)
	Distinct(key string, result interface{}) error
	Explain(result interface{}) error
	For(result interface{}, f func() error) error
	Hint(indexKey ...string) *mgo.Query
	Iter() *Iter
	Limit(n int) *mgo.Query
	LogReplay() *mgo.Query
	MapReduce(job *mgo.MapReduce, result interface{}) (info *mgo.MapReduceInfo, err error)
	Max(max interface{}) *mgo.Query
	Min(min interface{}) *mgo.Query
	One(result interface{}) (err error)
	Prefetch(p float64) *mgo.Query
	Select(selector interface{}) *mgo.Query
	SetMaxScan(n int) *mgo.Query
	SetMaxTime(d time.Duration) *mgo.Query
	Skip(n int) *mgo.Query
	Snapshot() *mgo.Query
	Sort(fields ...string) *mgo.Query
	Tail(timeout time.Duration) *Iter
}

type Iter interface {
	Err() error
	Close() error
	Done() bool
	Timeout() bool
	Next(result interface{}) bool
	All(result interface{}) error
}

type Pipe interface {
	All(result interface{}) error
	AllowDiskUse() *Pipe
	Batch(n int) *Pipe
	Collation(collation *mgo.Collation) *Pipe
	Explain(result interface{}) error
	Iter() *Iter
	One(result interface{}) error
	SetMaxTime(d time.Duration) *Pipe
}

type ChangeStream interface {
	Close() error
	Err() error
	Next(result interface{}) bool
	ResumeToken() *bson.Raw
	Timeout() bool
}

type Bulk interface {
	Insert(docs ...interface{})
	Remove(selectors ...interface{})
	RemoveAll(selectors ...interface{})
	Run() (*mgo.BulkResult, error)
	Unordered()
	Update(pairs ...interface{})
	UpdateAll(pairs ...interface{})
	Upsert(pairs ...interface{})
}
