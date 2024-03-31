package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	Id       int     `jason:"id"`
	Name     string  `jason:"name"`
	Quantity int     `jason:"quantity"`
	Price    float64 `jason:"price"`
}

func Getproducts(db *sql.DB) ([]product, error) {
	rows, err := db.Query("select * from products")

	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var data product
		err := rows.Scan(&data.Id, &data.Name, &data.Quantity, &data.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, data)
	}
	return products, nil
}

func (Prod *product) GetProduct(db *sql.DB) error {
	query := fmt.Sprintf("select name,quantity,price from products where id=%v", Prod.Id)
	row := db.QueryRow(query)
	err := row.Scan(&Prod.Name, &Prod.Price, &Prod.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (Prod *product) CreateProduct(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)", Prod.Name, Prod.Quantity, Prod.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	Prod.Id = int(id)
	return nil
}

func (Prod *product) UpdateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", Prod.Name, Prod.Quantity, Prod.Price, Prod.Id)
	result, _ := db.Exec(query)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return err
	}
	return nil
}

func (Prod *product) DeleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=%v", Prod.Id)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
