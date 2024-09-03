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
	TaskQueue string
}

type amqpClient struct {
	conn *amqp.Connection
	ch *amqp.Channel
	cf *ConfigOptions
}

type Log struct {
	client amqpClient
}

type level string

type MesageLog struct {
	Level level `json:"level"`
	Value string `json:"value"`
}

const (
	INFO level = "INFO"
	ERROR level = "ERROR"
)

func New(cf *ConfigOptions) (*Log, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", cf.Username, cf.Password, cf.Host, cf.Port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	
	client := amqpClient{conn: conn, ch: ch, cf: cf}

	return &Log{client: client}, nil
}

func (c *amqpClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (l *Log) Info(lvl level, msg string) (error) {
	return l.log(map[string]any{
		"level": lvl,
		"message": msg,
	})
}

func (l *Log) Error(lvl level, msg string) (error) {
	return l.log(map[string]any{
		"level": lvl,
		"message": msg,
	})
}

func (l *Log) log(msg map[string]any) (error) {
	q, err := l.client.ch.QueueDeclare(
		l.client.cf.TaskQueue,
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

	return l.client.ch.Publish(
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
