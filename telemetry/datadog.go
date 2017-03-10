package telemetry

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

type Datadog interface {
	Incr(name string, tags []string)
	Failure(title, text string)
	Gauge(name string, tags []string, value int)
}

type Client struct {
	client      *statsd.Client
	initialized bool
}

func (c *Client) on_error(err error) {
	log.Println(err)
}

func (c *Client) Failure(title, text string) {
	if c.initialized {
		title = c.client.Namespace + title
		ev := statsd.NewEvent(title, text)
		ev.AlertType = statsd.Error
		err := errors.Wrap(c.client.Event(ev), "Failed sending event to DataDog")
		if err != nil {
			c.on_error(err)
		}
	}
}

func (c *Client) Incr(name string, tags []string) {
	if c.initialized {
		err := c.client.Incr(name, tags, 1.0)
		if err != nil {
			c.on_error(err)
		}
	}
}

func (c *Client) Gauge(name string, tags []string, value int) {
	if c.initialized {
		err := c.client.Gauge(name, float64(value), tags, 1.0)
		if err != nil {
			c.on_error(err)
		}
	}
}

func NewDatadog(addr string, namespace string, keyspace string) Datadog {
	c, err := statsd.New(addr)
	if err != nil {
		client := &Client{}
		client.on_error(err)
		return client
	}
	c.Namespace = namespace + "."
	c.Tags = append(c.Tags, "keyspace:"+keyspace)
	return &Client{client: c, initialized: true}
}
