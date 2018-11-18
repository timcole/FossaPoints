// This is used to emulate the Google Cloud Function caller
package main

import (
	"net/http"

	"github.com/TimothyCole/FossaPoints/p"
)

func main() {
	http.HandleFunc("/", p.GetPoints)

	panic(http.ListenAndServe(":8080", nil))
}
