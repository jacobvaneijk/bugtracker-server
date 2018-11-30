package trello

import (
    "fmt"
    "bytes"
    "net/http"
    "encoding/json"
    "mime/multipart"
)

type Client struct {
    Client  *http.Client
    BaseURL string
    Key     string
    Token   string
}

func NewClient(baseURL string, key string, token string) *Client {
    return &Client{
        Client:     http.DefaultClient,
        BaseURL:    baseURL,
        Key:        key,
        Token:      token,
    }
}

func (c *Client) NewRequest(method string, path string, body *bytes.Buffer) (*http.Request, error) {
    url := fmt.Sprintf("%s/%s", c.BaseURL, path)

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    q := req.URL.Query()

    if c.Key != "" {
        q.Add("key", c.Key)
    }

    if c.Token != "" {
        q.Add("token", c.Token)
    }

    req.URL.RawQuery = q.Encode()

    return req, nil
}

func (c *Client) NewUploadRequest(path string, body []byte) (*http.Request, error) {
    buffer := &bytes.Buffer{}
    writer := multipart.NewWriter(buffer)

    part, err := writer.CreateFormFile("file", "asdf.png")
    if err != nil {
        return nil, err
    }

    part.Write(body)

    writer.Close()

    req, err := c.NewRequest("POST", path, buffer)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if v != nil {
        decoder := json.NewDecoder(resp.Body)
        return decoder.Decode(v)
    }

    return nil
}
