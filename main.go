// This is used to emulate the Google Cloud Function caller
package main

import (
	"net/http"

	p "github.com/TimothyCole/FossaPoints/ls"
)

func main() {
	http.HandleFunc("/", p.LastSeen)

	panic(http.ListenAndServe(":8080", nil))
}
