package trello

import (
    "fmt"
    "bytes"
)

type Card struct {
    Client *Client

    ID      string  `json:"id"`
    Name    string  `json:"name"`
    Url     string  `json:"url"`
    Desc    string  `json:"desc"`
    Closed  bool    `json:"closed"`
    IDList  string  `json:"idList"`
}

func (c *Client) CreateCard(card *Card) error {
    req, err := c.NewRequest("POST", "cards", &bytes.Buffer{})
    if err != nil {
        return err
    }

    q := req.URL.Query()

    q.Add("name", card.Name)
    q.Add("desc", card.Desc)
    q.Add("idList", card.IDList)

    req.URL.RawQuery = q.Encode()

    err = c.Do(req, card)
    if err != nil {
        return err
    }

    card.Client = c

    return nil
}

func (c *Card) AddAttachment(file []byte) error {
    path := fmt.Sprintf("cards/%s/attachments", c.ID)

    req, err := c.Client.NewUploadRequest(path, file)
    if err != nil {
        return err
    }

    return c.Client.Do(req, nil)
}
