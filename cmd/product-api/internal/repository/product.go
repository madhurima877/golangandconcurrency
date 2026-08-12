package repository

import (
	"database/sql"

	"product-api/internal/model"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(DB *sql.DB) *ProductRepository {
	return &ProductRepository{DB: DB}
}

func (r *ProductRepository) Create(p *model.Product) error {
	query := `
		INSERT INTO products (name, description, price, category, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.DB.QueryRow(
		query,
		p.Name,
		p.Description,
		p.Price,
		p.Category,
		p.Stock,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	p := &model.Product{}
	query := `
		SELECT id, name, description, price, category, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Category,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}
