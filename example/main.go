package main

import (
	"fmt"
	"log"

	"github.com/vincer2040/lexiGo/pkg/lexigo"
)

func main() {
	client := lexigo.New("127.0.0.1:5173")
	defer client.Close()
	err := client.Connect()
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	res, err := client.Ping()
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(res)

	res, err = client.Set("foo", "bar")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(res)

	res, err = client.Set("bar", "foo")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(res)

	res, err = client.Set("baz", "foo")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(res)

	res, err = client.Get("foo")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(res)

	keys, err := client.Keys()
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", keys)

	t, err := client.Type("foo")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", t)

	resi, err := client.Incr("myValue")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Printf("%d\n", resi)

	resi, err = client.Del("foo")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	fmt.Println(resi)
}
