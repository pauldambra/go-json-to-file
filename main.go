package main

import (
	"fmt"
	"net/http"

	"github.com/pauldambra/filesaver/api"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(":3000", api.Handlers())
}
