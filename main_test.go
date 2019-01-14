package main

import (
    "os"
    "testing"

    "github.com/joho/godotenv"
)


func TestMain(m *testing.M) {
    godotenv.Load()
    /*
    if err != nil {
        panic(err)
    }
    */

    os.Exit(m.Run())
}
