package main

import (
	"log"
	"net/http"

	"github.com/Rhymond/api-demo/utils"
)

func main() {
	c := utils.NewClient()

	resp, err := c.Do(http.MethodGet, "/v2/status", nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = utils.PrettyPrint(resp)
}
