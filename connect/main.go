package main

import (
	"fmt"
	"log"
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: go run main.go <endpoint> <bucket> <username> <password>")
	}

	endpoint := os.Args[1]
	bucketName := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://"+endpoint+"?ssl=no_verify", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)
	col := bucket.DefaultCollection()

	// Create a N1QL Primary Index (but ignore if it exists)
	err = cluster.QueryIndexes().CreatePrimaryIndex(bucketName, &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	type User struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}

	id := uuid.New().String()

	// Create and store a Document
	_, err = col.Upsert(id,
		User{
			ID:        id,
			Name:      "Arthur",
			Email:     "kingarthur@couchbase.com",
			Interests: []string{"Holy Grail", "African Swallows"},
		}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the document back
	getResult, err := col.Get(id, nil)
	if err != nil {
		log.Fatal(err)
	}

	var inUser User
	err = getResult.Content(&inUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %v\n", inUser)

	// Perform a N1QL Query
	queryResult, err := cluster.Query(
		fmt.Sprintf("SELECT id, name FROM `%s` WHERE $1 IN interests", bucketName),
		&gocb.QueryOptions{PositionalParameters: []interface{}{"African Swallows"}},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print each found Row
	for queryResult.Next() {
		var result interface{}
		err := queryResult.Row(&result)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}
}
