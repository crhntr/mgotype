package mgofake

type Session struct{}

func (sess *Session) Clone() *Session {
	return sess
}

func (sess *Session) Close() error {
	return nil
}

func (sess *Session) DB(name string) Database {
	return Database{}
}

type Database struct{}

func (db *Database) C(name string) Collection {
	return Collection{}
}

type Collection struct{}
