package main

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "encoding/base64"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/jacobvaneijk/bugtracker-server/trello"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API is up and running!")
}

func CreateBugHandler(w http.ResponseWriter, r *http.Request) {
    client := trello.NewClient(
        os.Getenv("TRELLO_API_URL"),
        os.Getenv("TRELLO_API_KEY"),
        os.Getenv("TRELLO_API_TOKEN"),
    )

    card := trello.Card{
        Name: r.FormValue("title"),
        Desc: r.FormValue("description"),
        IDList: r.FormValue("list"),
    }

    err := client.CreateCard(&card)
    if err != nil {
        panic(err)
    }

    image, err := base64.StdEncoding.DecodeString(r.FormValue("screenshot"))
    if err != nil {
        panic(err)
    }

    err = card.AddAttachment(image)
    if err != nil {
        panic(err)
    }
}

func main() {
    err := godotenv.Load()
    if err != nil {
        panic(err)
    }

    r := mux.NewRouter()

    r.HandleFunc("/status", StatusHandler)
    r.HandleFunc("/bug", CreateBugHandler).Methods("POST")

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":4321", nil))
}
