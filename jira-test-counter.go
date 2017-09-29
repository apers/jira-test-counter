package main

import (
    "fmt"
    "net/http"
    "time"
    "log"
)

const ServerPort = "80";

func requestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Request recieved: %s", r.Body);
    fmt.Fprintf(w, "Hi there: %s", r.URL.Path[1:])
}

func main() {
    fmt.Println(time.Now())
    fmt.Println("Staring server..")

    http.HandleFunc("/", requestHandler);
    log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
