package postgres

import (
	"context"
	"go-data-catalog/internal/models"
)

type ContactRepository struct {
	db *DB
}

func NewContactRepository(db *DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) GetAllContacts(ctx context.Context) ([]models.Contact, error) {
	query := `
		SELECT id, name, telegram_contact, created_at
		FROM contacts
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.TelegramContact,
			&contact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (r *ContactRepository) GetContactByID(ctx context.Context, id int) (*models.Contact, error) {
	query := `
		SELECT id, name, telegram_contact, created_at
		FROM contacts
		WHERE id = $1
	`
	
	var contact models.Contact
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&contact.ID,
		&contact.Name,
		&contact.TelegramContact,
		&contact.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &contact, nil
}

func (r *ContactRepository) CreateContact(ctx context.Context, contact *models.Contact) error {
	query := `
		INSERT INTO contacts (name, telegram_contact)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		contact.Name,
		contact.TelegramContact,
	).Scan(&contact.ID, &contact.CreatedAt)
	
	return err
}

func (r *ContactRepository) UpdateContact(ctx context.Context, id int, contact *models.Contact) error {
	query := `
		UPDATE contacts 
		SET name = $2, telegram_contact = $3
		WHERE id = $1
		RETURNING created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
		contact.Name,
		contact.TelegramContact,
	).Scan(&contact.CreatedAt)
	
	if err != nil {
		return err
	}
	
	contact.ID = id
	return nil
}

func (r *ContactRepository) DeleteContact(ctx context.Context, id int) error {
	query := `DELETE FROM contacts WHERE id = $1`
	
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}