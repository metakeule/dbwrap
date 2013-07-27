package dbwrap

import (
	"database/sql"
	"database/sql/driver"
)

/*
   registeres a Wrapper as a database driver with the given name
   and returns it.

   the Wrapper is based on the given inner driver and will forward any
   method call to the underlying innerDriver

   before you open a connection to the Wrapper, you may set the different handlers on it
   it set, they are called instead of the method of the innerDriver and they are given the
   underlying driver.Conn, so that you may call it

   to open a connection, use it the usual way, e.g.

       sql.Open(name, connectString)

   where name is the given name here and connectString is the same you would
   pass to innerDriver

*/
func New(name string, innerDriver driver.Driver) (ø *Wrapper) {
	ø = &Wrapper{Driver: innerDriver}
	sql.Register(name, ø)
	return
}

type Wrapper struct {
	driver.Driver

	// is called after a new connection (innerConn) has been successfully returned by innerDriver
	// you have to return a driver.Conn here to the library. in most cases you
	// will want to return innerConn here
	HandleOpen func(name string, innerConn driver.Conn) (driver.Conn, error)

	// is called instead of innerConn.Begin, which you might want to call at some point
	HandleBegin func(innerConn driver.Conn) (driver.Tx, error)

	// is called instead of innerConn.Prepare, which you might want to call at some point
	HandlePrepare func(innerConn driver.Conn, query string) (driver.Stmt, error)

	// is called instead of innerConn.Close, which you might want to call at some point
	HandleClose func(innerConn driver.Conn) error

	// is called instead of innerConn.(driver.Execer), if innerConn may be casted to a driver.Execer
	// you might want to call conn.Exec at some point
	HandleExec func(conn driver.Execer, query string, args []driver.Value) (driver.Result, error)

	// is called instead of innerConn.(driver.Queryer), if innerConn may be casted to a driver.Queryer
	// you might want to call conn.Query at some point
	HandleQuery func(conn driver.Queryer, query string, args []driver.Value) (driver.Rows, error)

	// is called before each method
	BeforeAll func(conn driver.Conn, event string, data ...interface{})

	// is called after each method
	AfterAll func(conn driver.Conn, event string, data ...interface{})
}

func (ø *Wrapper) Open(name string) (con driver.Conn, err error) {
	c, err := ø.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	con = &conn{c, ø}
	ex, isEx := c.(driver.Execer)
	if isEx {
		con = &execConn{con, ex, ø}
	}

	qr, isQu := c.(driver.Queryer)
	if isQu {
		con = &queryConn{con, qr, ø}
	}
	if ø.HandleOpen != nil {
		return ø.HandleOpen(name, con)
	}
	return
}
