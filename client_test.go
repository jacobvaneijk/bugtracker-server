package main

import (
    "os"
    "testing"
)

func TestNewClient(t *testing.T) {
    client := NewClient(
        os.Getenv("TRELLO_API_URL"),
        os.Getenv("TRELLO_API_KEY"),
        os.Getenv("TRELLO_API_TOKEN"),
    )

    if client == nil {
        t.Errorf("NewClient() returned nil")
    }
}
