package main

import (
	"fmt"
	"gopkg.in/metakeule/dbwrap.v2"
)

var fake, db = dbwrap.NewFake()

func q1() {
	fake.SetNumInputs(1)
	db.Query("Select ?", "hiho")
	q, v := fake.LastQuery()
	fmt.Println(q, v)
}

func q2() {
	fake.SetNumInputs(0)
	db.Exec("Delete * from mytable")
	q, v := fake.LastQuery()
	fmt.Println(q, v)
}

func q3() {
	fake.SetNumInputs(1)
	stmt, _ := db.Prepare("select * from ?")
	stmt.Query("mytable")
	q, v := fake.LastQuery()
	fmt.Println(q, v)
}

func q4() {
	fake.SetNumInputs(1)
	db.QueryRow("select 1 from ? ", "mytable")
	q, v := fake.LastQuery()
	fmt.Println(q, v)
}

func main() {
	q1()
	q2()
	q3()
	q4()
}
