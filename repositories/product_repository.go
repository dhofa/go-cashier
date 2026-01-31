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
// =====================
// GET ALL
// =====================
// =====================
// GET ALL
// =====================
func (repo *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
			c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var catID *int
		var catName, catDesc *string

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &catID, &catName, &catDesc); err != nil {
			return nil, err
		}

		if catID != nil {
			p.Category = &models.Category{
				ID:          *catID,
				Name:        *catName,
				Description: "",
			}
			if catDesc != nil {
				p.Category.Description = *catDesc
			}
		}

		products = append(products, p)
	}

	return products, nil
}

// =====================
// CREATE
// =====================
// =====================
// CREATE
// =====================
func (repo *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return repo.db.
		QueryRow(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).
		Scan(&product.ID)
}

// =====================
// GET BY ID
// =====================
// =====================
// GET BY ID
// =====================
// =====================
// GET BY ID
// =====================
func (repo *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
			c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	var catID *int
	var catName, catDesc *string

	err := repo.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &catID, &catName, &catDesc)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	if catID != nil {
		p.Category = &models.Category{
			ID:          *catID,
			Name:        *catName,
			Description: "",
		}
		if catDesc != nil {
			p.Category.Description = *catDesc
		}
	}

	return &p, nil
}

// =====================
// UPDATE
// =====================
// =====================
// UPDATE
// =====================
func (repo *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, stock = $3, category_id = $4
		WHERE id = $5
	`

	tag, err := repo.db.Exec(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
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
