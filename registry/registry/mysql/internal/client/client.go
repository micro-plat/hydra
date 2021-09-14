package client

import (
	"sync"
)

type Message struct {
	Action  string
	Payload string
}

type Client struct {
	mu       sync.Mutex
	channels map[string]*Channel
}

type Channel struct {
	chOnce  sync.Once
	chOnce1 sync.Once
	chValue chan *Message
	chChild chan *Message
}

const (
	ActionInsert = "insert"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// BulkRequest is used to send multi request in batch.
type BulkRequest struct {
	Action string
	Schema string
	Table  string
	Data   map[string]interface{}

	PkName  string
	PkValue string
}

// NewClient creates the Cient with configuration.
func NewClient() *Client {
	c := new(Client)
	return c
}

func (c *Client) PublishChild(channel, message string) {
	ch, ok := c.channels[channel]
	if !ok {
		return
	}
	msg := &Message{
		Action:  ActionUpdate,
		Payload: message,
	}
	ch.chChild <- msg
}

func (c *Client) PublishValue(channel, message string) {
	ch, ok := c.channels[channel]
	if !ok {
		return
	}
	msg := &Message{
		Action:  ActionUpdate,
		Payload: message,
	}
	ch.chValue <- msg
}

func (c *Client) Subscribe(channel string) *Channel {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.channels == nil {
		c.channels = make(map[string]*Channel)
	}
	r := &Channel{}
	c.channels[channel] = r
	return r
}

func (c *Channel) ChannelValue() <-chan *Message {
	c.chOnce.Do(func() {
		c.chValue = make(chan *Message, 100)
	})
	return c.chValue
}

func (c *Channel) ChannelChild() <-chan *Message {
	c.chOnce1.Do(func() {
		c.chChild = make(chan *Message, 100)
	})
	return c.chChild
}

// Bulk sends the bulk request.
func (c *Client) Bulk(reqs []*BulkRequest) error {
	for _, req := range reqs {
		msg := &Message{}

		ch, ok := c.channels[req.PkValue]
		if !ok {
			continue
		}

		switch req.Action {
		case ActionDelete:
			if cap(ch.chChild) == 0 {
				continue
			}
			msg.Action = ActionDelete
			ch.chChild <- msg
		case ActionUpdate:
			if cap(ch.chValue) == 0 {
				continue
			}
			msg.Action = ActionUpdate
			ch.chValue <- msg
		default:
			if cap(ch.chChild) == 0 {
				continue
			}
			msg.Action = ActionInsert
			ch.chChild <- msg
		}
	}

	return nil
}
