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

func (r *ContactRepository) GetAllContacts(ctx context.Context, teamID int) ([]models.Contact, error) {
	query := `
		SELECT id, name, telegram_contact, team_id, created_at
		FROM contacts
		WHERE team_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, teamID)
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
			&contact.TeamID,
			&contact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (r *ContactRepository) GetContactByID(ctx context.Context, teamID, id int) (*models.Contact, error) {
	query := `
		SELECT id, name, telegram_contact, team_id, created_at
		FROM contacts
		WHERE id = $1 AND team_id = $2
	`
	
	var contact models.Contact
	err := r.db.Pool.QueryRow(ctx, query, id, teamID).Scan(
		&contact.ID,
		&contact.Name,
		&contact.TelegramContact,
		&contact.TeamID,
		&contact.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &contact, nil
}

func (r *ContactRepository) CreateContact(ctx context.Context, teamID int, contact *models.Contact) error {
	query := `
		INSERT INTO contacts (name, telegram_contact, team_id)
		VALUES ($1, $2, $3)
		RETURNING id, team_id, created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		contact.Name,
		contact.TelegramContact,
		teamID,
	).Scan(&contact.ID, &contact.TeamID, &contact.CreatedAt)
	
	return err
}

func (r *ContactRepository) UpdateContact(ctx context.Context, teamID, id int, contact *models.Contact) error {
	query := `
		UPDATE contacts 
		SET name = $3, telegram_contact = $4
		WHERE id = $1 AND team_id = $2
		RETURNING team_id, created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
		teamID,
		contact.Name,
		contact.TelegramContact,
	).Scan(&contact.TeamID, &contact.CreatedAt)
	
	if err != nil {
		return err
	}
	
	contact.ID = id
	return nil
}

func (r *ContactRepository) DeleteContact(ctx context.Context, teamID, id int) error {
	query := `DELETE FROM contacts WHERE id = $1 AND team_id = $2`
	
	_, err := r.db.Pool.Exec(ctx, query, id, teamID)
	return err
}
