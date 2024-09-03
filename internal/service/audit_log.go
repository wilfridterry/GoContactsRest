package service

import (
	"time"
)

type AuditLog interface {
	Log(logMsg LogMessage) error 
}

type action string
type entity string

const (
	ACTION_REGISTER action = "REGISTER"
	ACTION_LOGIN    action = "LOGIN"
	ACTION_CREATE   action = "CREATE"
	ACTION_GET      action = "GET"
	ACTION_UPDATE   action = "UPDATE"
	ACTION_DELETE   action = "DELETE"

	ENTITY_CONTACT entity = "CONTACT"
	ENTITY_USER    entity = "USER"
)

type LogMessage struct {
	Action    action
	Entity    entity
	EntityID  int64
	Timestamp time.Time
}

type AMQPClient interface {
	Log(msg map[string]any) error
}

type AuditLogService struct {
	client AMQPClient
}

func NewAuditLog(client AMQPClient) *AuditLogService {
	return &AuditLogService{client}
}

func (s *AuditLogService) Log(logMsg LogMessage) error {
	return s.client.Log(map[string]any{
		"ation":     logMsg.Action,
		"entity":    logMsg.Entity,
		"entity_id": logMsg.EntityID,
		"timestamp": logMsg.Timestamp,
	})
}
