package service

import (
	domain "contact-list/internal/domain/contact"
	"context"
)

type Contacts struct {
	repository ContactRepository
}

type ContactRepository interface{
	GetAll(context.Context) ([]domain.Contact, error)
	GetById(context.Context, int64) (*domain.Contact, error)
	Create(context.Context, *domain.SaveInputContact) error
	Delete(context.Context, int64) error
	Update(context.Context, int64, *domain.SaveInputContact) error
}

func (c *Contacts) All(ctx context.Context) ([]domain.Contact, error) {
	return c.repository.GetAll(ctx)
}

func (c *Contacts) GetOne(ctx context.Context, id int64) (*domain.Contact, error) {
	return c.repository.GetById(ctx, id)
}

func (c *Contacts) Create(ctx context.Context, inp *domain.SaveInputContact) error {
	return c.repository.Create(ctx, inp)
}

func (c *Contacts) Update(ctx context.Context, id int64, inp *domain.SaveInputContact) error {
	return c.repository.Update(ctx, id, inp)
}

func (c *Contacts) Delete(ctx context.Context, id int64) error {
	return c.repository.Delete(ctx, id)
}

func NewContacts(repository ContactRepository) *Contacts {
	return &Contacts{repository}
}