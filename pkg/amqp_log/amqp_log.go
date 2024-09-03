package amqplog

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfigOptions struct {
	Username string
	Password string
	Host string
	Port int
	Queue string
}

type Client struct {
	conn *amqp.Connection
	ch *amqp.Channel
	cf *ConfigOptions
}

// type Log struct {
// 	client amqpClient
// }

// type level string

// type MesageLog struct {
// 	Level level `json:"level"`
// 	Value string `json:"value"`
// }

// const (
// 	INFO level = "INFO"
// 	ERROR level = "ERROR"
// )

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

// func (l *Log) Info(lvl level, msg string) (error) {
// 	return l.log(map[string]any{
// 		"level": lvl,
// 		"message": msg,
// 	})
// }

// func (l *Log) Error(lvl level, msg string) (error) {
// 	return l.log(map[string]any{
// 		"level": lvl,
// 		"message": msg,
// 	})
// }

func (c *Client) Log(msg map[string]any) (error) {
	q, err := c.ch.QueueDeclare(
		c.cf.Queue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgBts, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.ch.Publish(
		"",
		q.Name,
		false, 
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: msgBts,
		},
	)
}
