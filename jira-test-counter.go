package main

import (
    "fmt"
    "net/http"
    "time"
    "log"
)

const ServerPort = "80";

func requestHandler(writer http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(writer, "Hi there: %s", req.URL.Path[1:])
}

func main() {
    fmt.Println(time.Now())
    fmt.Println("Staring server..")

    http.HandleFunc("/", requestHandler);
    log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
