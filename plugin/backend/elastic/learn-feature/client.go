package dp_golearn_to_rank

import (
	"gopkg.in/olivere/elastic.v5"
	"context"
	"net/url"
)

const (
	DefaultUrl = "http://127.0.0.1:9200"
)

// Embedds an elasticsearch client, which by default performs health checks every
// 60 seconds
type Client struct {
	c *elastic.Client
	ctx context.Context
}

func NewClient() (*Client, error) {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	c := &Client{
		c: client,
		ctx: ctx,
	}

	return c, nil
}

func (client *Client) BaseClient() (*elastic.Client) {
	// Returns the base elasticsearch client
	return client.c
}

func (client *Client) PerformRequest(method, path string, params url.Values, body interface{}, ignoreErrors ...int) (*elastic.Response, error) {
	return client.c.PerformRequest(client.ctx, method, path, params, body, ignoreErrors...)
}

func (client *Client) FeatureService() (*FeatureService) {
	return NewFeatureService(client)
}
