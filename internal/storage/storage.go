package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"rent-watcher/internal/models"
)

type Storage interface {
	SaveProperty(property *models.Property, rawData string) (bool, error)
}

type SQLStorage struct {
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) Storage {
	return &SQLStorage{db: db}
}

func (s *SQLStorage) SaveProperty(property *models.Property, rawData string) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var existingID string
	err = tx.QueryRow("SELECT id FROM properties WHERE id = ?", property.ID).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("failed to check existing property: %w", err)
	}

	isNew := errors.Is(err, sql.ErrNoRows)

	if isNew {
		err = s.insertProperty(tx, property, rawData)
	} else {
		err = s.updateProperty(tx, property, rawData)
	}

	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return isNew, nil
}

func (s *SQLStorage) insertProperty(tx *sql.Tx, property *models.Property, rawData string) error {
	_, err := tx.Exec(`
		INSERT INTO properties 
		(id, first_photo, price, logradouro, bairro, cidade, metragem, quartos, banheiros, suites, garagens, tipo_imovel, distance_meters) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		property.ID, property.FirstPhoto, property.Price, property.Logradouro, property.Bairro, property.Cidade,
		property.Metragem, property.Quartos, property.Banheiros, property.Suites, property.Garagens, property.TipoImovel,
		property.DistanceMeters)
	if err != nil {
		return fmt.Errorf("failed to insert property: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO raw_data 
		(id, json_data) 
		VALUES (?, ?)`,
		property.ID, rawData)
	if err != nil {
		return fmt.Errorf("failed to insert raw data: %w", err)
	}

	return nil
}

func (s *SQLStorage) updateProperty(tx *sql.Tx, property *models.Property, rawData string) error {
	_, err := tx.Exec(`
		UPDATE properties 
		SET first_photo = ?, price = ?, logradouro = ?, bairro = ?, cidade = ?, metragem = ?, 
			quartos = ?, banheiros = ?, suites = ?, garagens = ?, tipo_imovel = ?, distance_meters = ?
		WHERE id = ?`,
		property.FirstPhoto, property.Price, property.Logradouro, property.Bairro, property.Cidade,
		property.Metragem, property.Quartos, property.Banheiros, property.Suites, property.Garagens, property.TipoImovel,
		property.DistanceMeters, property.ID)
	if err != nil {
		return fmt.Errorf("failed to update property: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE raw_data 
		SET json_data = ?
		WHERE id = ?`,
		rawData, property.ID)
	if err != nil {
		return fmt.Errorf("failed to update raw data: %w", err)
	}

	return nil
}
