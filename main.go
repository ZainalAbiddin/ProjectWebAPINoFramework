package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// initial server
	http.HandleFunc("/", func(wr http.ResponseWriter, rq *http.Request) {
		wr.Header().Set("Content-type", "application/json")
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte(`{"Pesan":"server aktif"}`))
	})
	// cek error dan rute alamat
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
