package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return 0, err
	}
	defer db.Close()

	rows, err := db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Printf("Ошибка добавления записи в tracker.db: %s", err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := rows.LastInsertId()
	if err != nil {
		fmt.Printf("Ошибка получения идентификаора добавленной записи: %s", err)
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return Parcel{}, err
	}
	defer db.Close()

	p := Parcel{}

	row := db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Printf("Ошибка чтения данных tracker.db: %s", err)
		return Parcel{}, err
	}
	// заполните объект Parcel данными из таблицы

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return nil, err
	}
	defer db.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel

	rows, err := db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Printf("Ошибка получения строк из tracker.db: %s", err)
		return nil, err
	}
	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Printf("Ошибка чтения данных tracker.db: %s", err)
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Ошибка перебора строк в цикле: %s", err)
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Printf("Ошибка обновлеия статуса в tracker.db: %s", err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Printf("Ошибка обновления адреса в tracker.db: %s", err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Printf("Ошибка открытия файла tracker.db: %s", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Printf("Ошибка удаления строки в tracker.db: %s", err)
		return err
	}
	return nil
}
