package main

import (
	"fmt"
	"net/http"

	"github.com/mcsymiv/web-hello-world/pkg/hand"
)

const port = ":8080"

func main() {
	http.HandleFunc("/", hand.Index)
	http.HandleFunc("/home", hand.Home)
	http.HandleFunc("/about", hand.About)
	http.HandleFunc("/exit", hand.Exit)

	fmt.Println("Started app.\nListen on port :8080")
	_ = http.ListenAndServe(port, nil)
}
