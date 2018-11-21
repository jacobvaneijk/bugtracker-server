package main

import (
    "net/http"
    "os"
    "log"
    "fmt"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API is up and running!");
}

func CreateBugHandler(w http.ResponseWriter, r *http.Request) {
    req, err := http.NewRequest(
        "POST",
        os.Getenv("TRELLO_API_URL"),
        nil,
    )

    if err != nil {
        fmt.Printf("http.NewRequest() error: %v\n", err)
        return
    }

    q := req.URL.Query()

    q.Add("key", os.Getenv("TRELLO_API_KEY"))
    q.Add("token", os.Getenv("TRELLO_API_TOKEN"))
    q.Add("name", r.FormValue("title"))
    q.Add("desc", r.FormValue("description"))
    q.Add("idList", r.FormValue("list"))

    req.URL.RawQuery = q.Encode()
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    c := &http.Client{}
    resp, err := c.Do(req)

    if err != nil {
        fmt.Printf("http.Do() error: %v\n", err)
        return
    }

    defer resp.Body.Close()

    fmt.Fprintf(w, "Created new card!");
}

func main() {
    err := godotenv.Load()

    if err != nil {
        log.Fatal("Couldn't read `.env` file.")
    }

    r := mux.NewRouter()

    r.HandleFunc("/status", StatusHandler)
    r.HandleFunc("/bug", CreateBugHandler).Methods("POST")

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":4321", nil))
}
