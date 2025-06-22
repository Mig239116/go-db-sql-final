package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет запись в базу данных
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO parcel(client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client,
		p.Status,
		p.Address,
		p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add the parcel %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get entry id %w", err)
	}
	return int(id), nil
}

// Get извлекает запись с конкретным идентификатором из базы
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number=?", number)
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("failed to get entry with id %d from database %w", number, err)
	}

	return p, nil
}

// GetByClients извлекает все записи по конкретному клиенту из базы данных
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client=?", client)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries for client %d from database %w", client, err)
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		row := Parcel{}
		err := rows.Scan(&row.Number, &row.Client, &row.Status, &row.Address, &row.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to read entries for client %d from database %w", client, err)
		}
		res = append(res, row)
	}
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to get entries for client %d from database %w", client, err)
    }
	return res, nil
}

// SetStatus обновляет статус для указанной записи
func (s ParcelStore) SetStatus(number int, status string) error {
	res, err := s.db.Exec("UPDATE parcel SET status=? WHERE number=?", status, number)
	if err != nil {
		return fmt.Errorf("cant update status of parcel with number %d", number)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
        return fmt.Errorf("failed to check rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("parcel with number %d not found", number)
    }
	return nil
}

// SetAddress обновляет адрес для указанной записи
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(
		"UPDATE parcel SET address=? WHERE number = ? AND status = ?", 
		address,
		number,
		ParcelStatusRegistered,
	)
	if err != nil {
		return fmt.Errorf("cant update parcel with number %d", number)
	}
	return nil
}

// Delete удаляет указанную запись
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number=? AND status=?", number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("cant delete parcel with number %d", number)
	}
	return nil
}
