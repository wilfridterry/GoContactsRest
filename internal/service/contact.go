package service

import (
	"context"
	"time"

	"github.com/wilfridterry/contact-list/internal/domain"

	"github.com/sirupsen/logrus"
	// audit "github.com/wilfridterry/audit-log/pkg/domain"
)

type Contacts struct {
	repository  ContactRepository
	auditClient AuditClient
	auditLog    AuditLog
}

type ContactRepository interface {
	GetAll(context.Context) ([]domain.Contact, error)
	GetById(context.Context, int64) (*domain.Contact, error)
	Create(context.Context, *domain.SaveInputContact) (int64, error)
	Delete(context.Context, int64) error
	Update(context.Context, int64, *domain.SaveInputContact) error
}

func (c *Contacts) All(ctx context.Context) ([]domain.Contact, error) {
	return c.repository.GetAll(ctx)
}

func (service *Contacts) GetOne(ctx context.Context, id int64) (*domain.Contact, error) {
	contact, err := service.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action: audit.ACTION_GET,
	// 	Entity: audit.ENTITY_CONTACT,
	// 	EntityID: contact.ID,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Contacts.Get",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_GET,
		Entity:    ENTITY_CONTACT,
		EntityID:  contact.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Contacts.Get",
		}).Error("failed to send log request:", err)
	}

	return contact, nil
}

func (service *Contacts) Create(ctx context.Context, inp *domain.SaveInputContact) error {
	id, err := service.repository.Create(ctx, inp)
	if err != nil {
		return err
	}

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action: audit.ACTION_CREATE,
	// 	Entity: audit.ENTITY_CONTACT,
	// 	EntityID: id,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Contacts.Create",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_CREATE,
		Entity:    ENTITY_CONTACT,
		EntityID:  id,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Contacts.Create",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (service *Contacts) Update(ctx context.Context, id int64, inp *domain.SaveInputContact) error {
	err := service.repository.Update(ctx, id, inp)
	if err != nil {
		return err
	}

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action: audit.ACTION_UPDATE,
	// 	Entity: audit.ENTITY_CONTACT,
	// 	EntityID: id,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Contacts.Update",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_UPDATE,
		Entity:    ENTITY_CONTACT,
		EntityID:  id,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Contacts.Update",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (service *Contacts) Delete(ctx context.Context, id int64) error {
	err := service.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action: audit.ACTION_DELETE,
	// 	Entity: audit.ENTITY_CONTACT,
	// 	EntityID: id,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Contacts.Delete",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_DELETE,
		Entity:    ENTITY_CONTACT,
		EntityID:  id,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Contacts.Delete",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func NewContacts(repository ContactRepository, auditClient AuditClient, auditLog AuditLog) *Contacts {
	return &Contacts{
		repository:  repository,
		auditClient: auditClient,
		auditLog:    auditLog,
	}
}
