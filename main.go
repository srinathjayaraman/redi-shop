package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Redi-shop started")
	fmt.Println("Awaiting requests...")
	err := fasthttp.ListenAndServe(":8000", getRouter())
	if err != nil {
		panic(err)
	}
}
