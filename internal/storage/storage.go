package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"rent-watcher/internal/models"
)

type Storage interface {
	GetProperty(propertyID string) (*models.Property, error)
	PropertyExists(propertyID string) (bool, error)
	SaveOrUpdateProperty(property *models.Property, rawData string) error
}

type SQLStorage struct {
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) Storage {
	return &SQLStorage{db: db}
}

func (s *SQLStorage) GetProperty(propertyID string) (*models.Property, error) {
	var property models.Property
	err := s.db.QueryRow(`
		SELECT id, first_photo, price, logradouro, bairro, cidade, metragem, quartos, banheiros, suites, garagens, tipo_imovel, distance_meters
		FROM properties WHERE id = ?`, propertyID).Scan(
		&property.ID, &property.FirstPhoto, &property.Price, &property.Logradouro, &property.Bairro, &property.Cidade,
		&property.Metragem, &property.Quartos, &property.Banheiros, &property.Suites, &property.Garagens, &property.TipoImovel,
		&property.DistanceMeters)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("property not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get property: %w", err)
	}
	return &property, nil
}

func (s *SQLStorage) SaveOrUpdateProperty(property *models.Property, rawData string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := s.propertyExistsInTx(tx, property.ID)
	if err != nil {
		return fmt.Errorf("failed to check property existence: %w", err)
	}

	if exists {
		err = s.updatePropertyData(tx, property)
	} else {
		err = s.insertPropertyData(tx, property)
	}

	if err != nil {
		return err
	}

	err = s.upsertRawData(tx, property.ID, rawData)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLStorage) propertyExistsInTx(tx *sql.Tx, propertyID string) (bool, error) {
	var id string
	err := tx.QueryRow("SELECT id FROM properties WHERE id = ?", propertyID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SQLStorage) insertProperty(tx *sql.Tx, property *models.Property, rawData string) error {
	if err := s.insertPropertyData(tx, property); err != nil {
		return err
	}
	return s.insertRawData(tx, property.ID, rawData)
}

func (s *SQLStorage) updateProperty(tx *sql.Tx, property *models.Property, rawData string) error {
	if err := s.updatePropertyData(tx, property); err != nil {
		return err
	}
	return s.updateRawData(tx, property.ID, rawData)
}

func (s *SQLStorage) insertPropertyData(tx *sql.Tx, property *models.Property) error {
	_, err := tx.Exec(`
		INSERT INTO properties 
		(id, first_photo, price, logradouro, bairro, cidade, metragem, quartos, banheiros, suites, garagens, tipo_imovel, distance_meters, condominio, total_price)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		property.ID, property.FirstPhoto, property.Price, property.Logradouro, property.Bairro, property.Cidade,
		property.Metragem, property.Quartos, property.Banheiros, property.Suites, property.Garagens, property.TipoImovel,
		property.DistanceMeters, property.Condominio, property.TotalPrice)
	if err != nil {
		return fmt.Errorf("failed to insert property: %w", err)
	}
	return nil
}

func (s *SQLStorage) updatePropertyData(tx *sql.Tx, property *models.Property) error {
	_, err := tx.Exec(`
		UPDATE properties 
		SET first_photo = ?, price = ?, logradouro = ?, bairro = ?, cidade = ?, metragem = ?, 
			quartos = ?, banheiros = ?, suites = ?, garagens = ?, tipo_imovel = ?, condominio = ?, total_price = ?
		WHERE id = ?`,
		property.FirstPhoto, property.Price, property.Logradouro, property.Bairro, property.Cidade,
		property.Metragem, property.Quartos, property.Banheiros, property.Suites, property.Garagens, property.TipoImovel,
		property.Condominio, property.TotalPrice, property.ID)
	if err != nil {
		return fmt.Errorf("failed to update property: %w", err)
	}
	return nil
}

func (s *SQLStorage) upsertRawData(tx *sql.Tx, id, rawData string) error {
	_, err := tx.Exec(`
		INSERT OR REPLACE INTO raw_data (id, json_data)
		VALUES (?, ?)`,
		id, rawData)
	if err != nil {
		return fmt.Errorf("failed to upsert raw data: %w", err)
	}
	return nil
}

func (s *SQLStorage) insertRawData(tx *sql.Tx, id, rawData string) error {
	_, err := tx.Exec(`
		INSERT INTO raw_data 
		(id, json_data) 
		VALUES (?, ?)`,
		id, rawData)
	if err != nil {
		return fmt.Errorf("failed to insert raw data: %w", err)
	}
	return nil
}

func (s *SQLStorage) updateRawData(tx *sql.Tx, id, rawData string) error {
	_, err := tx.Exec(`
		UPDATE raw_data 
		SET json_data = ?
		WHERE id = ?`,
		rawData, id)
	if err != nil {
		return fmt.Errorf("failed to update raw data: %w", err)
	}
	return nil
}

func (s *SQLStorage) PropertyExists(propertyID string) (bool, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM properties WHERE id = ?", propertyID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check property existence: %w", err)
	}
	return true, nil
}
