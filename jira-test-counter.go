package main

import (
    "fmt"
    "net/http"
    "time"
    "log"
    "io"
    "bytes"
)

const ServerPort = "80";

func readReader(reader io.Reader) *bytes.Buffer {
    buf := new(bytes.Buffer)
    buf.ReadFrom(reader)
    return buf;
}


func requestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there: %s", r.URL.Path[1:])
}

func main() {
    fmt.Println(time.Now())
    fmt.Println("Staring server..")

    http.HandleFunc("/webhook", requestHandler);
    log.Fatal(http.ListenAndServe(":"+ServerPort, nil))
}
