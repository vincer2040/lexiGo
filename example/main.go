package main

import (
	"fmt"
	"log"

	"github.com/vincer2040/lexiGo/pkg/client"
)

func main() {
	client := client.New("127.0.0.1:5173")
	defer client.Close()
	err := client.Connect()
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	res, err := client.Set("foo", "bar")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	fmt.Println(res)
}