package dbwrap

import (
	"database/sql"
	"database/sql/driver"
	"testing"
)

var f, dbF = NewFake()

func TestFakeQuery(t *testing.T) {
	sq := "select 1"
	dbF.Query(sq)

	r1, _ := f.LastQuery()
	if r1 != sq {
		t.Errorf("last simple query should return %s, but returns %#v\n", sq, r1)
	}

	f.SetNumInputs(1)
	sq = "select ?"
	dbF.Query(sq, "1")

	r2, p2 := f.LastQuery()
	if r2 != sq {
		t.Errorf("last query should return %s, but returns %#v\n", sq, r2)
	}

	exp, _ := driver.DefaultParameterConverter.ConvertValue("1")

	if p2[0] != exp {
		t.Errorf("last query should return have parameters %#v, but returns %#v\n", exp, p2)
	}

	f.SetNumInputs(0)
	sq = "select ?"
	_, err := dbF.Query(sq, "1")
	if err == nil {
		t.Errorf("querying with wrong number of parameters should return error, but does not\n")
	}
}

func TestFakeExec(t *testing.T) {
	sq := "select 1"
	dbF.Exec(sq)

	r1, _ := f.LastQuery()
	if r1 != sq {
		t.Errorf("last simple exec should return %s, but returns %#v\n", sq, r1)
	}

	f.SetNumInputs(1)
	sq = "select ?"
	dbF.Exec(sq, "1")

	r2, p2 := f.LastQuery()
	if r2 != sq {
		t.Errorf("last exec should return %s, but returns %#v\n", sq, r2)
	}

	exp, _ := driver.DefaultParameterConverter.ConvertValue("1")

	if p2[0] != exp {
		t.Errorf("last exec should return have parameters %#v, but returns %#v\n", exp, p2)
	}

	f.SetNumInputs(0)
	sq = "select ?"
	_, err := dbF.Exec(sq, "1")
	if err == nil {
		t.Errorf("execing with wrong number of parameters should return error, but does not\n")
	}
}

func TestFakePrepare(t *testing.T) {
	f.SetNumInputs(1)
	sq := "select ?"
	stmt, _ := dbF.Prepare(sq)
	stmt.Exec("1")
	r1, p1 := f.LastQuery()
	if r1 != sq {
		t.Errorf("last prepare with exec should return %s, but returns %#v\n", sq, r1)
	}

	exp, _ := driver.DefaultParameterConverter.ConvertValue("1")

	if p1[0] != exp {
		t.Errorf("last prepare with exec should return have parameters %#v, but returns %#v\n", exp, p1)
	}

	stmt.Query("1")
	r2, p2 := f.LastQuery()
	if r2 != sq {
		t.Errorf("last prepare with query should return %s, but returns %#v\n", sq, r2)
	}

	if p2[0] != exp {
		t.Errorf("last prepare with query should return have parameters %#v, but returns %#v\n", exp, p2)
	}

	stmt.QueryRow("1")
	r3, p3 := f.LastQuery()
	if r3 != sq {
		t.Errorf("last prepare with query row should return %s, but returns %#v\n", sq, r3)
	}

	if p3[0] != exp {
		t.Errorf("last prepare with query row should return have parameters %#v, but returns %#v\n", exp, p3)
	}
}

func TestFakeQueryRow(t *testing.T) {
	sq := "select 1"
	dbF.QueryRow(sq)

	r1, _ := f.LastQuery()
	if r1 != sq {
		t.Errorf("last simple query row should return %s, but returns %#v\n", sq, r1)
	}

	f.SetNumInputs(1)
	sq = "select ?"
	dbF.QueryRow(sq, "1")

	r2, p2 := f.LastQuery()
	if r2 != sq {
		t.Errorf("last query row should return %s, but returns %#v\n", sq, r2)
	}

	exp, _ := driver.DefaultParameterConverter.ConvertValue("1")

	if p2[0] != exp {
		t.Errorf("last query row should return have parameters %#v, but returns %#v\n", exp, p2)
	}

	f.SetNumInputs(0)
	sq = "select ?"
	err := dbF.QueryRow(sq, "1")
	if err == nil {
		t.Errorf("querying a row with wrong number of parameters should return error, but does not\n")
	}
}

type wrapResults_ struct {
	Query  string
	Values []driver.Value
}

var results = &wrapResults_{}
var dbwrap = New("test", f)
var dbW, _ = sql.Open("test", "")

func init() {
	dbwrap.HandleExec = func(x driver.Execer, q string, v []driver.Value) (driver.Result, error) {
		results.Query = q
		results.Values = v
		return dbwrap.Driver.(driver.Execer).Exec(q, v)
	}

	dbwrap.HandleQuery = func(r driver.Queryer, qs string, v []driver.Value) (driver.Rows, error) {
		results.Query = qs
		results.Values = v
		return dbwrap.Driver.(driver.Queryer).Query(qs, v)
	}

	dbwrap.HandlePrepare = func(c driver.Conn, qs string) (driver.Stmt, error) {
		results.Query = qs
		return c.Prepare(qs)
	}
}

func TestDbwrapPassthrough(t *testing.T) {
	sq := "select 2"
	dbW.Exec(sq)

	l, _ := f.LastQuery()
	if l != sq {
		t.Errorf("drwap should pass an exec to the underlying driver, the inner fake should have %s, but has: %s", sq, l)
	}

	sq = "select 3"
	dbW.Query(sq)

	l, _ = f.LastQuery()
	if l != sq {
		t.Errorf("drwap should pass a query to the underlying driver, the inner fake should have %s, but has: %s", sq, l)
	}
	sq = "select 4"
	dbW.QueryRow(sq)

	l, _ = f.LastQuery()
	if l != sq {
		t.Errorf("drwap should pass a query row to the underlying driver, the inner fake should have %s, but has: %s", sq, l)
	}

	f.SetNumInputs(2)
	sq = "select ?, ?"
	dbW.Query(sq, "2", "3")

	l1, p1 := f.LastQuery()
	if l1 != sq {
		t.Errorf("drwap should pass a query row to the underlying driver, the inner fake should have %s, but has: %s", sq, l1)
	}

	exp1, _ := driver.DefaultParameterConverter.ConvertValue("2")
	exp2, _ := driver.DefaultParameterConverter.ConvertValue("3")

	if p1[0] != exp1 {
		t.Errorf("drwap should pass a query row to the underlying driver should have first parameter %#v, but has %#v\n", exp1, p1[0])
	}

	if p1[1] != exp2 {
		t.Errorf("drwap should pass a query row to the underlying driver should have second parameter %#v, but has %#v\n", exp2, p1[1])
	}

	f.SetNumInputs(1)
	sq = "select ?"
	stmt, _ := dbW.Prepare(sq)

	stmt.Exec("2")
	l2, p2 := f.LastQuery()
	if l2 != sq {
		t.Errorf("drwap should pass a query row to the underlying driver, the inner fake should have %s, but has: %s", sq, l2)
	}
	if p2[0] != exp1 {
		t.Errorf("drwap should pass a query row to the underlying driver should have second parameter %#v, but has %#v\n", exp1, p2[0])
	}

}

func TestDbwrapCatch(t *testing.T) {
	sq := "select 1"
	dbW.Exec(sq)
	if results.Query != sq {
		t.Errorf("drwap should have tacked the exec and have %s, but has: %s", sq, results.Query)
	}
	sq = "select 2"
	dbW.Query(sq)
	if results.Query != sq {
		t.Errorf("drwap should have tacked the query and have %s, but has: %s", sq, results.Query)
	}
}
