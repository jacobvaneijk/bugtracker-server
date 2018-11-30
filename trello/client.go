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

/*
func (c *Client) Multipart(path string, headers map[string]string, body *bytes.Buffer) error {
    url := fmt.Sprintf("%s/%s", c.BaseURL, path)

    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return err
    }

    q := req.URL.Query()

    if c.Key != "" {
        q.Add("key", c.Key)
    }

    if c.Token != "" {
        q.Add("token", c.Token)
    }

    for key, val := range headers {
        req.Header.Set(key, val)
    }

    req.URL.RawQuery = q.Encode()

    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    return nil
}
*/

func (c *Client) Post(path string, params map[string]string, target interface{}) error {
    url := fmt.Sprintf("%s/%s", c.BaseURL, path)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return err
    }

    q := req.URL.Query()

    if c.Key != "" {
        q.Add("key", c.Key)
    }

    if c.Token != "" {
        q.Add("token", c.Token)
    }

    for key, val := range params {
        q.Add(key, val)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    req.URL.RawQuery = q.Encode()

    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if target != nil {
        decoder := json.NewDecoder(resp.Body)
        err = decoder.Decode(target)
        if err != nil {
            return err
        }
    }

    return nil
}
