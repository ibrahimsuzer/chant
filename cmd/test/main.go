package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/genjidb/genji/sql/driver"
	_ "github.com/genjidb/genji/sql/driver"
)

type User struct {
	ID      int64
	Name    string
	Age     uint32
	Address struct {
		City    string
		ZipCode string
	}
}

func main() {
	
	// Create a sql/database DB instance
	db, err := sql.Open("genji", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table. Genji tables are schemaless by default, you don't need to specify a schema.
	_, err = db.Exec("create table user")
	if err != nil {
		panic(err)
	}

	// Create an index.
	_, err = db.Exec("create index idx_user_name on user (name)")
	if err != nil {
		panic(err)
	}

	// Insert some data
	_, err = db.Exec("insert into user (id, name, age) values (?, ?, ?)", 10, "foo", 15)
	if err != nil {
		panic(err)
	}

	// Insert some data using document notation
	_, err = db.Exec(`insert into user values {id: 12, "name": "bar", age: ?, address: {city: "Lyon", zipcode: "69001"}}`, 16)
	if err != nil {
		panic(err)
	}

	// Structs can be used to describe a document
	_, err = db.Exec("insert into user values ?, ?", &User{ID: 1, Name: "baz", Age: 100}, &User{ID: 2, Name: "bat"})
	if err != nil {
		panic(err)
	}

	// Query some documents
	stream, err := db.Query("select * from user where id > ?", 1)
	if err != nil {
		panic(err)
	}
	// always close the result when you're done with it
	defer stream.Close()
	
	for stream.Next() {
		var u User

		err := stream.Scan(driver.Scanner(&u))
		if err != nil {
			fmt.Printf("err %v: %s", u, err)

			return 
		}
		fmt.Println(u)

	}
	

}
