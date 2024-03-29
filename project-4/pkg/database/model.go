package database

import (
	"database/sql"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *Product) GetProduct(db *sql.DB) error {

	err := db.QueryRow("SELECT name, price FROM products WHERE id = $1", p.ID).Scan(&p.Name, &p.Price)

	return err
}

func (p *Product) UpdateProduct(db *sql.DB) error {

	if p.Name == "" {
		_, err := db.Exec("UPDATE products SET price=$1 WHERE id=$2", p.Price, p.ID)
		return err
	} else if p.Price == 0 {
		_, err := db.Exec("UPDATE products SET name=$1 WHERE id=$2", p.Name, p.ID)
		return err
	}

	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)

	return err
}

func (p *Product) DeleteProduct(db *sql.DB) error {

	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)
	return err
}

func (p *Product) CreateProduct(db *sql.DB) error {

	// The user won't know the id of product
	// the DB'll know the id as its is SERIAL value type
	err := db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2)  RETURNING id", p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetProducts(db *sql.DB) ([]Product, error) {

	products := []Product{}

	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price)

		if err != nil {
			return nil, err
		} else {
			products = append(products, p)
		}
	}

	return products, nil
}
