package main

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCreateCard(t *testing.T) {
    client := NewClient(
        os.Getenv("TRELLO_API_URL"),
        os.Getenv("TRELLO_API_KEY"),
        os.Getenv("TRELLO_API_TOKEN"),
    )

    assert.NotNil(t, client)

    card := Card{
        Name: "Test Card",
        Desc: "Test Description",
        IDList: "5bf52bcc81455a2f35ba8e11",
    }

    err := client.CreateCard(&card)

    assert.Nil(t, err)
    assert.NotEqual(t, card.ID, "", "the card must receive an ID from Trello")
}

func TestAddAttachment(t *testing.T) {
    // @TODO: Figure out how to test this.
}
