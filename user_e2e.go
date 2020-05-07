package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

func main() {
	start := 500
	var wg sync.WaitGroup
	wg.Add(start)

	for i := 0; i < start; i++ {
		go func() {
			checkUserE2E()
			wg.Done()
		}()
	}
	fmt.Println("started all")

	wg.Wait()

	fmt.Println("done")
}

func checkUserE2E() {
	server := "http://localhost:8000"

	resp, err := req.Post(server + "/users/create")
	checkErr(err)

	userID, err := resp.ToString()
	checkErr(err)

	resp, err = req.Post(server + "/users/credit/add/" + userID + "/43")
	checkErr(err)

	respString, err := resp.ToString()
	checkErr(err)

	success := respString == "success"
	if !success {
		log.Error("adding credit failed")
	}

	resp, err = req.Post(server + "/users/credit/subtract/" + userID + "/1")
	checkErr(err)

	respString, err = resp.ToString()
	checkErr(err)

	if respString != "success" {
		log.Error("adding credit failed")
	}

	resp, err = req.Get(server + "/users/find/" + userID)
	checkErr(err)

	respString, err = resp.ToString()
	if respString != fmt.Sprintf("(%s, 42)", userID) {
		log.Error("invalid value for user, should be (userID, 42), but was: " + respString)
	}

	resp, err = req.Get(server + "/users/credit/" + userID)
	checkErr(err)

	respString, err = resp.ToString()
	if respString != "42" {
		log.Error("invalid value for user credit, should be 42 but was: " + respString)
	}

	resp, err = req.Delete(server + "/users/remove/" + userID)
	checkErr(err)

	respString, err = resp.ToString()
	checkErr(err)

	if respString != "success" {
		log.Error("removing user failed")
	}

	resp, err = req.Get(server + "/users/find/" + userID)
	checkErr(err)

	if resp.Response().StatusCode != http.StatusNotFound {
		log.Error("user should not be found after deleting")
	}

	fmt.Printf("Done for user %s\n", userID)
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
