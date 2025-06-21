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
		return 0, fmt.Errorf("failed to add the parcel %s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get entry id %s", err)
	}
	return int(id), nil
}

// Get извлекает запись с конкретным идентификатором из базы
func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number=?", number)
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("failed to get entry with id %d from database %s", number, err)
	}

	return p, nil
}

// GetByClients извлекает все записи по конкретному клиенту из базы данных
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client=?", client)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries for client %d from database %s", client, err)
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		row := Parcel{}
		err := rows.Scan(&row.Number, &row.Client, &row.Status, &row.Address, &row.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to read entries for client %d from database %s", client, err)
		}
		res = append(res, row)
	}
	return res, nil
}

// SetStatus обновляет статус для указанной записи
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("No parcel with number %d exists", number)
	}
	_, err1 := s.db.Exec("UPDATE parcel SET status=?", status)
	if err1 != nil {
		return fmt.Errorf("Cant update status of parcel with number %d", number)
	}
	return nil
}

// SetAddress обновляет адрес для указанной записи
func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("No parcel with number %d exists", number)
	}
	if parcel.Status == ParcelStatusRegistered {
		_, err1 := s.db.Exec("UPDATE parcel SET address=?", address)
		if err1 != nil {
			return fmt.Errorf("Cant update parcel with number %d", number)
		}
	}
	return nil
}

// Delete удаляет указанную запись
func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("No parcel with number %d exists", number)
	}
	if parcel.Status == ParcelStatusRegistered {
		_, err1 := s.db.Exec("DELETE FROM parcel WHERE number=?", number)
		if err1 != nil {
			return fmt.Errorf("Cant delete parcel with number %d", number)
		}
	}
	return nil
}
