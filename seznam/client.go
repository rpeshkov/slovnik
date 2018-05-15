package seznam

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/rpeshkov/slovnik"
)

const (
	scheme     = "https"
	seznamHost = "slovnik.seznam.cz"

	wordQueryVar      = "q"
	shortViewQueryVar = "shortView"
)

// Map of urls used for translation
var urls = map[slovnik.Language]string{
	slovnik.Cz: "cz-ru",
	slovnik.Ru: "ru-cz",
}

type Client struct {
	client *http.Client
}

// NewClient creates a client for accessing slovnik.seznam.cz portal
func NewClient(httpClient *http.Client) *Client {
	var c *http.Client

	if httpClient == nil {
		c = http.DefaultClient
	} else {
		c = httpClient
	}

	return &Client{
		client: c,
	}
}

// Get requests translation result page for provided word
func (c *Client) Get(word string, language slovnik.Language) (io.ReadCloser, error) {
	query := c.createURL(word, language)
	resp, err := c.client.Get(query.String())

	if err != nil {
		return nil, errors.Wrap(err, "get failed")
	}

	return resp.Body, nil
}

func (c *Client) createURL(word string, language slovnik.Language) url.URL {
	v := url.Values{}
	v.Add(wordQueryVar, word)
	v.Add(shortViewQueryVar, "0")

	return url.URL{
		Scheme:   scheme,
		Host:     seznamHost,
		Path:     urls[language],
		RawQuery: v.Encode(),
	}
}
