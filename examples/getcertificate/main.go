package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Rhymond/api-demo/utils"
)

type AllowListRequest struct {
	CIDR     string `json:"cidrBlock"`
	RuleType string `json:"ruleType"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run main.go <cluster-id>")
	}

	c := utils.NewClient()

	resp, err := c.Do(http.MethodGet, "/v2/clusters/"+os.Args[1]+"/certificate", nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = utils.PrettyPrint(resp)
}
