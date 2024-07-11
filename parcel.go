package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	//добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, nil
	}

	//идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// чтение строки по заданному number
	// из таблицы возвращается только одна строка
	// заполнение объекта Parcel данными из таблицы
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// чтение строк из таблицы parcel по заданному client
	// из таблицы может вернуться несколько строк
	// заполнение среза Parcel данными из таблицы
	var res []Parcel
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status=:status WHERE number=:number",
		sql.Named("status", status), sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address=:address WHERE number=:number AND status=:status",
		sql.Named("status", ParcelStatusRegistered), sql.Named("number", number), sql.Named("address", address))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number=:number AND status=:status",
		sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	return err

}
