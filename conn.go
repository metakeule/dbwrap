package dbwrap

import "database/sql/driver"

type conn struct {
	driver.Conn
	*Wrapper
}

func (ø *conn) Begin() (t driver.Tx, err error) {
	if ø.BeforeAll != nil {
		ø.BeforeAll(ø.Conn, "Begin")
	}
	if ø.HandleBegin != nil {
		t, err = ø.HandleBegin(ø.Conn)
	} else {
		t, err = ø.Conn.Begin()
	}
	if ø.AfterAll != nil {
		ø.AfterAll(ø.Conn, "Begin", t, err)
	}
	return
}

func (ø *conn) Prepare(query string) (st driver.Stmt, err error) {
	if ø.BeforeAll != nil {
		ø.BeforeAll(ø.Conn, "Prepare", query)
	}
	if ø.HandlePrepare != nil {
		st, err = ø.HandlePrepare(ø.Conn, query)
	} else {
		st, err = ø.Conn.Prepare(query)
	}
	if ø.AfterAll != nil {
		ø.AfterAll(ø.Conn, "Prepare", st, err)
	}
	return
}

func (ø *conn) Close() (err error) {
	if ø.BeforeAll != nil {
		ø.BeforeAll(ø.Conn, "Close")
	}
	if ø.HandleClose != nil {
		err = ø.HandleClose(ø.Conn)
	} else {
		err = ø.Conn.Close()
	}
	if ø.AfterAll != nil {
		ø.AfterAll(ø.Conn, "Close", err)
	}
	return
}

type execConn struct {
	driver.Conn
	driver.Execer
	*Wrapper
}

func (ø *execConn) Exec(query string, args []driver.Value) (res driver.Result, err error) {
	if ø.BeforeAll != nil {
		ø.BeforeAll(ø.Execer.(driver.Conn), "Exec", query, args)
	}
	if ø.HandleExec != nil {
		res, err = ø.HandleExec(ø.Execer, query, args)
	} else {
		res, err = ø.Execer.Exec(query, args)
	}
	if ø.AfterAll != nil {
		ø.AfterAll(ø.Execer.(driver.Conn), "Exec", res, err)
	}
	return
}

type queryConn struct {
	driver.Conn
	driver.Queryer
	*Wrapper
}

func (ø *queryConn) Query(query string, args []driver.Value) (res driver.Rows, err error) {
	if ø.BeforeAll != nil {
		ø.BeforeAll(ø.Queryer.(driver.Conn), "Query", query, args)
	}
	if ø.HandleQuery != nil {
		res, err = ø.HandleQuery(ø.Queryer, query, args)
	} else {
		res, err = ø.Queryer.Query(query, args)
	}
	if ø.AfterAll != nil {
		ø.AfterAll(ø.Queryer.(driver.Conn), "Query", res, err)
	}
	return
}
