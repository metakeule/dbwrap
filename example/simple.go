package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-on/pq"
	"github.com/metakeule/dbwrap"
	"os"
)

// we need pqdrv since lib/pq does not export it's driver and export a toplevel Open method instead
type pqdrv int

// fullfill the driver.Driver interface
func (d pqdrv) Open(name string) (driver.Conn, error) { return pq.Open(name) }

func main() {
	// creates a new wrap and registers it with the name "debug"
	wrap := dbwrap.New("debug", pqdrv(0))

	// use the new name to connect
	db, err := sql.Open("debug", connectString)
	if err != nil {
		panic(err.Error())
	}

	// sets the search_path to a scheme before each connection
	// the conn is already made but not yet returned to the sql library
	wrap.HandleOpen = func(name string, conn driver.Conn) (driver.Conn, error) {
		conn.(driver.Execer).Exec(`SET search_path = "public"`, []driver.Value{})
		// return the conn to the library
		return conn, nil
	}

	// is used instead of conn.Exec
	wrap.HandleExec = func(exec driver.Execer, query string, args []driver.Value) (driver.Result, error) {
		fmt.Println("exec: ", query)
		// do the real Exec and return the result
		return exec.Exec(query, args)
	}

	// called before each method call
	wrap.BeforeAll = func(conn driver.Conn, event string, data ...interface{}) {
		vals := []interface{}{"before: ", event}
		vals = append(vals, data...)
		fmt.Println(vals...)
	}

	query(db)
}

func query(db *sql.DB) {
	r, err := db.Query(`SELECT 'Donald Duck' as "name", 'Hi' as "message"`)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	defer r.Close()
	for r.Next() {
		var name, msg string
		err = r.Scan(&name, &msg)
		if err != nil {
			fmt.Println("ERROR while scanning: " + err.Error())
			return
		}
		fmt.Printf("%s says %#v\n", name, msg)
	}
}

var connectString string

func init() {
	if os.Getenv("PG_URL") == "" {
		panic("please set the environment variable PG_URL as database connection string")
	}
	var err error
	connectString, err = pq.ParseURL(os.Getenv("PG_URL"))
	if err != nil {
		panic(err.Error())
	}
}
