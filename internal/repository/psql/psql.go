package psql

import (
	domain "contact-list/internal/domain/contact"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Contacts struct {
	Conn *pgx.Conn
}

func NewContacts(conn *pgx.Conn) *Contacts {
	return &Contacts{conn}
}

func (repo *Contacts) GetAll(ctx context.Context) ([]domain.Contact, error) {
	rows, err := repo.Conn.Query(ctx, "SELECT * FROM contacts")
	if err != nil {
		return nil, err
	}

	var contacts []domain.Contact

	for rows.Next() {
		c := domain.Contact{}
		if err := rows.Scan(&c.ID, &c.Name, &c.LastName, &c.Phone, &c.Email, &c.Address, &c.Author, &c.CreatedAt, &c.UpdatedAt);
			err != nil {
				return nil, err
			}

		contacts = append(contacts, c)
	}

	return contacts, nil
}

func (repo *Contacts) GetById(ctx context.Context, id int64) (*domain.Contact, error) {
	c := domain.Contact{}
	row := repo.Conn.QueryRow(ctx, "SELECT * from contatcts WHERE id = $1", id)

	if err := row.Scan(&c.ID, &c.Name, &c.LastName, &c.Phone, &c.Email, &c.Address, &c.Author, &c.CreatedAt, &c.UpdatedAt);
		err != nil {
			return &domain.Contact{}, domain.ErrContactNotFound
		}

	return &c, nil
}

func (repo *Contacts) Create(ctx context.Context, inp *domain.SaveInputContact) error {
	_, err := repo.Conn.Exec(
		ctx, 
		"INSERT INTO contacts (name, last_name, phone, email, address, author) values ($1, $2, $3, $4, $5, $6)",
		inp.Name, inp.LastName, inp.Phone, inp.Email, inp.Address, inp.Author,
		)

	return err
}

func (repo *Contacts) Delete(ctx context.Context, id int64) error {
	_, err := repo.Conn.Exec(ctx, "DELETE FROM contacts WHERE id = $1", id)

	return err
}

func (repo *Contacts) Update(ctx context.Context, id int64, inp *domain.SaveInputContact) error {
	args := make([]interface{}, 0)
	fields := make([]string, 0)
	argInd := 1

	if inp.Name != nil {
		fields = append(fields, fmt.Sprintf("name=$%d", argInd))
		args[argInd] = inp.Name
		argInd++
	}

	if inp.LastName != nil {
		fields = append(fields, fmt.Sprintf("last_name=$%d", argInd))
		args[argInd] = inp.LastName
		argInd++
	}

	if inp.Phone != nil {
		fields = append(fields, fmt.Sprintf("phone=$%d", argInd))
		args[argInd] = inp.Phone
		argInd++
	}

	if inp.Email != nil {
		fields = append(fields, fmt.Sprintf("email=$%d", argInd))
		args[argInd] = inp.Email
		argInd++
	}

	if inp.Address != nil {
		fields = append(fields, fmt.Sprintf("address=$%d", argInd))
		args[argInd] = inp.Address
		argInd++
	}

	if inp.Author != nil {
		fields = append(fields, fmt.Sprintf("author=$%d", argInd))
		args[argInd] = inp.Author
		argInd++
	}

	argInd++
	args[argInd] = id
	setQuery := strings.Join(fields, ", ")
	query := fmt.Sprintf("UPDATE contacts set %s WHERE id=$%d", setQuery, argInd)

	_, err := repo.Conn.Exec(ctx, query, args...)

	return err
}