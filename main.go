package main

import (
  "log"
  "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello niggha"))
}


func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)


  log.Println("Starting server on :6000")
  err := http.ListenAndServe(":6000", mux)
  log.Fatal(err)
}