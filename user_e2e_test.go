package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// This test communicates with a running instance of redi-shop and concurrently creates,
// updates and removes users.

func TestUser(t *testing.T) {
	start := 500
	var wg sync.WaitGroup
	wg.Add(start)

	for i := 0; i < start; i++ {
		go func() {
			checkUserE2E(t)
			wg.Done()
		}()
	}
	fmt.Println("started all")

	wg.Wait()

	fmt.Println("done")
}

func checkUserE2E(t *testing.T) {
	client := req.New()
	assert := assert.New(t)
	server := "http://localhost:8000"

	resp, err := client.Post(server + "/users/create")
	if err != nil {
		assert.FailNow(err.Error())
	}

	userID, err := resp.ToString()
	assert.NoError(err)

	resp, err = client.Post(server + "/users/credit/add/" + userID + "/43")
	checkErr(assert, err)

	respString, err := resp.ToString()
	checkErr(assert, err)

	success := respString == "success"
	if !success {
		log.Error("adding credit failed")
	}

	resp, err = client.Post(server + "/users/credit/subtract/" + userID + "/1")
	checkErr(assert, err)

	respString, err = resp.ToString()
	checkErr(assert, err)

	if respString != "success" {
		log.Error("adding credit failed")
	}

	resp, err = client.Get(server + "/users/find/" + userID)
	checkErr(assert, err)

	respString, err = resp.ToString()
	checkErr(assert, err)
	if respString != fmt.Sprintf("(%s, 42)", userID) {
		log.Error("invalid value for user, should be (userID, 42), but was: " + respString)
	}

	resp, err = client.Get(server + "/users/credit/" + userID)
	checkErr(assert, err)

	respString, err = resp.ToString()
	checkErr(assert, err)
	if respString != "42" {
		log.Error("invalid value for user credit, should be 42 but was: " + respString)
	}

	resp, err = client.Delete(server + "/users/remove/" + userID)
	checkErr(assert, err)

	respString, err = resp.ToString()
	checkErr(assert, err)

	if respString != "success" {
		log.Error("removing user failed")
	}

	resp, err = client.Get(server + "/users/find/" + userID)
	checkErr(assert, err)

	if resp.Response().StatusCode != http.StatusNotFound {
		log.Error("user should not be found after deleting")
	}

	fmt.Printf("Done for user %s\n", userID)
}

func checkErr(assert *assert.Assertions, err error) {
	assert.NoError(err)
	if err != nil {
		panic("")
	}
}
