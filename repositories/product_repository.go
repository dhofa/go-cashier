package repositories

import (
	"context"
	"errors"

	"go-cashier/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// =====================
// GET ALL
// =====================
func (repo *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, price, stock FROM products`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// =====================
// CREATE
// =====================
func (repo *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, price, stock)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return repo.db.
		QueryRow(ctx, query, product.Name, product.Price, product.Stock).
		Scan(&product.ID)
}

// =====================
// GET BY ID
// =====================
func (repo *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, price, stock
		FROM products
		WHERE id = $1
	`

	var p models.Product
	err := repo.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	return &p, nil
}

// =====================
// UPDATE
// =====================
func (repo *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, stock = $3
		WHERE id = $4
	`

	tag, err := repo.db.Exec(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

// =====================
// DELETE
// =====================
func (repo *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	tag, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}
