package main

import (
    "os"
    "fmt"
    "log"
    "strconv"
    "net/http"
    "database/sql"
    "encoding/json"
    "encoding/base64"

    "github.com/gorilla/mux"
    "github.com/jmoiron/modl"
    _ "github.com/mattn/go-sqlite3"
    "github.com/jacobvaneijk/bugtracker-server/trello"
)

type App struct {
    Router *mux.Router
    DbMap *modl.DbMap
}

type BugReport struct {
    ID int `json:"-"`
    ProjectID int `json:"project_id"`
    Title string `json:"title"`
    Description string `json:"description"`
    TrelloID string `json:"id"`
    CurrentList string `json:"current_list"`
    SelectionWidth int `json:"selection_width"`
    SelectionHeight int `json:"selection_height"`
    PageWidth int `json:"page_width"`
    PageHeight int `json:"page_height"`
    DotX int `json:"dot_x"`
    DotY int `json:"dot_y"`
}

type Project struct {
    ID int `json:"-"`
    Name string `json:"name"`
    UnresolvedList string `json:"unresolved_list"`
    SiteUrl string `site_url:"site_url"`
}

func NewApp() *App {
    a := App{}

    dsn := os.Getenv("SQLITE_DSN")
    if dsn == "" {
        panic("SQLITE_DSN env not set.")
    }

    a.DbMap = modl.NewDbMap(connect(dsn), modl.SqliteDialect{})
    if a.DbMap == nil {
        panic("Couldn't create modl.DbMap internally.")
    }

    a.DbMap.AddTableWithName(Project{}, "projects").SetKeys(true, "id")
    a.DbMap.AddTableWithName(BugReport{}, "bug_reports").SetKeys(true, "id")

    err := a.DbMap.CreateTablesIfNotExists()
    if err != nil {
        panic(err.Error())
    }

    a.Router = mux.NewRouter()

    a.Router.HandleFunc("/status", a.statusHandler).Methods("GET")
    a.Router.HandleFunc("/bugs", a.getBugsHandler).Methods("GET")
    a.Router.HandleFunc("/bugs", a.createBugHandler).Methods("POST")
    a.Router.HandleFunc("/projects", a.createProjectHandler).Methods("POST")

    return &a
}

func connect(dsn string) *sql.DB {
    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        panic(err.Error())
    }

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return db
}

func (a *App) Run() {
    log.Fatal(http.ListenAndServe(":4321", a.Router))
}

func (a *App) statusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API is up and running!")
}

func (a *App) getBugsHandler(w http.ResponseWriter, r *http.Request) {
    bugs := []BugReport{}

    err := a.DbMap.Select(&bugs, "SELECT * FROM bug_reports")
    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")

    err = json.NewEncoder(w).Encode(bugs)
    if err != nil {
        panic(err)
    }
}

func (a *App) createBugHandler(w http.ResponseWriter, r *http.Request) {
    client := trello.NewClient(
        os.Getenv("TRELLO_API_URL"),
        os.Getenv("TRELLO_API_KEY"),
        os.Getenv("TRELLO_API_TOKEN"),
    )

    project := Project{}

    err := a.DbMap.Get(&project, r.FormValue("project_id"))
    if err != nil {
        panic(err)
    }

    selectionWidth, err := strconv.Atoi(r.FormValue("selection_width"))
    if err != nil {
        panic(err)
    }

    selectionHeight, err := strconv.Atoi(r.FormValue("selection_height"))
    if err != nil {
        panic(err)
    }

    pageWidth, err := strconv.Atoi(r.FormValue("page_width"))
    if err != nil {
        panic(err)
    }

    pageHeight, err := strconv.Atoi(r.FormValue("page_height"))
    if err != nil {
        panic(err)
    }

    dotX, err := strconv.Atoi(r.FormValue("dot_x"))
    if err != nil {
        panic(err)
    }

    dotY, err := strconv.Atoi(r.FormValue("dot_y"))
    if err != nil {
        panic(err)
    }

    bugReport := BugReport{
        ProjectID: project.ID,
        Title: r.FormValue("title"),
        Description: r.FormValue("description"),
        SelectionWidth: selectionWidth,
        SelectionHeight: selectionHeight,
        PageWidth: pageWidth,
        PageHeight: pageHeight,
        DotX: dotX,
        DotY: dotY,
    }

    card := trello.Card{
        Name: bugReport.Title,
        Desc: bugReport.Description,
        IDList: project.UnresolvedList,
    }

    err = client.CreateCard(&card)
    if err != nil {
        panic(err)
    }

    bugReport.TrelloID = card.ID;

    image, err := base64.StdEncoding.DecodeString(r.FormValue("screenshot"))
    if err != nil {
        panic(err)
    }

    err = card.AddAttachment(image)
    if err != nil {
        panic(err)
    }

    err = a.DbMap.Insert(&bugReport)
    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")

    err = json.NewEncoder(w).Encode(bugReport)
    if err != nil {
        panic(err)
    }
}

func (a *App) createProjectHandler(w http.ResponseWriter, r *http.Request) {
    project := Project{
        Name: r.FormValue("name"),
        UnresolvedList: r.FormValue("unresolved_list"),
        SiteUrl: r.FormValue("site_url"),
    }

    err := a.DbMap.Insert(&project)
    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")

    err = json.NewEncoder(w).Encode(project)
    if err != nil {
        panic(err)
    }
}
