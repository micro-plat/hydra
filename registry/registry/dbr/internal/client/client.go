package client

import (
	"sync"
)

type Message struct {
	Action  string
	Payload string
}

type Client struct {
	mu      sync.Mutex
	chanMap map[string][]chan *Message
}

// NewClient creates the Cient with configuration.
func NewClient() *Client {
	c := &Client{
		chanMap: make(map[string][]chan *Message),
	}
	return c
}

func (c *Client) Publish(channel, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.chanMap[channel] {
		v <- &Message{Action: ActionUpdate, Payload: message}
	}
}

func (c *Client) Subscribe(channel string) chan *Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.chanMap[channel]; ok {
		r := make(chan *Message)
		c.chanMap[channel] = append(c.chanMap[channel], r)
		return r
	}
	r := make(chan *Message)
	c.chanMap[channel] = []chan *Message{r}
	return r
}
