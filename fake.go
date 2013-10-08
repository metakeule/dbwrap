package dbwrap

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

/*
   returns a new database and a fake object that you can use to check, what has been queried.
   Before running a query, first use SetNumInputs() to set the number of parameters you want to pass to
   Exec or Query and then run the query.
   After the Query you may run LastQuery() to get the last query and the last parameters that have been committed
   to the database.

   The fake database is not threadsafe and is not meant to be used in an asynchronous way.
   Instead, open a new fake database for every goroutine.
*/
func NewFake() (*Fake, *sql.DB) {
	name := fmt.Sprintf("fakedb_%v", time.Now().Nanosecond())
	f := &Fake{name: name}
	f.stmt = &FakeStmt{Fake: f}
	sql.Register(name, f)
	db, err := sql.Open(name, "")
	if err != nil {
		panic("can't open fakedb with name " + name + ": " + err.Error())
	}
	return f, db
}

type FakeStmt struct {
	*Fake
	lastValues []driver.Value
	lastQuery  string
	numInputs  int
}

func (f *FakeStmt) NumInput() int                                    { return f.numInputs }
func (f *FakeStmt) Exec(v []driver.Value) (r driver.Result, e error) { f.lastValues = v; return }
func (f *FakeStmt) Query(v []driver.Value) (r driver.Rows, e error)  { f.lastValues = v; return }

type Fake struct {
	name string
	stmt *FakeStmt
}

func (f *Fake) Open(name string) (con driver.Conn, err error) { return f, nil }
func (f *Fake) Begin() (tx driver.Tx, e error)                { return }
func (f *Fake) Close() (err error)                            { return }
func (f *Fake) Name() string                                  { return f.name }
func (f *Fake) SetNumInputs(i int)                            { f.stmt.numInputs = i }
func (f *Fake) LastQuery() (string, []driver.Value)           { return f.stmt.lastQuery, f.stmt.lastValues }
func (f *Fake) Prepare(q string) (driver.Stmt, error) {
	f.stmt.lastQuery, f.stmt.lastValues = q, []driver.Value{}
	return f.stmt, nil
}
