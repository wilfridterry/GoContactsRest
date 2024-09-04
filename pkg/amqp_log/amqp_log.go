package amqplog

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfigOptions struct {
	Username string
	Password string
	Host     string
	Port     int
	Queue    string
}

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	cf   *ConfigOptions
}

func New(cf *ConfigOptions) (*Client, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", cf.Username, cf.Password, cf.Host, cf.Port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, ch: ch, cf: cf}, nil
}

func (c *Client) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Log(msg map[string]any) error {
	q, err := c.ch.QueueDeclare(
		c.cf.Queue, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		return err
	}

	msgBts, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBts,
		},
	)
}

func (c *Client) GetLogs() (<-chan amqp.Delivery, error) {
	q, err := c.ch.QueueDeclare(
		c.cf.Queue, // name
		false,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		return nil, err
	}

	err = c.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return nil, err
	}

	return c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
