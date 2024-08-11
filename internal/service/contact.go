package service

import domain "contact-list/internal/domain/contact"

type Contacts struct {
}

func (c *Contacts) All() ([]domain.Contact, error) {
	return nil, nil
}

func (c *Contacts) GetOne(id int64) (*domain.Contact, error) {
	return &domain.Contact{}, nil
}

func NewContacts() *Contacts {
	return &Contacts{}
}