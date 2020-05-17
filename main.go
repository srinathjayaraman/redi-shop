package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/valyala/fasthttp"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=redi password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Redi-shop started")
	fmt.Println("Awaiting requests...")
	err = fasthttp.ListenAndServe(":8000", getUserRouter(db))
	if err != nil {
		panic(err)
	}
}
