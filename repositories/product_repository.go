package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"

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
func (repo *ProductRepository) GetAll(ctx context.Context, params models.PaginationParams) ([]models.Product, int, error) {
	// Query dasar
	baseQuery := `
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`
	
	// Filter Search
	var conditions []string
	var args []interface{}
	argIdx := 1

	if params.Search != "" {
		conditions = append(conditions, "(p.name ILIKE $"+strconv.Itoa(argIdx)+" OR CAST(p.price AS TEXT) ILIKE $"+strconv.Itoa(argIdx)+" OR CAST(p.stock AS TEXT) ILIKE $"+strconv.Itoa(argIdx)+")")
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Hitung Total Data (untuk pagination)
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalItems int
	err := repo.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Query Data dengan Pagination
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
			c.id, c.name, c.description
	` + baseQuery + whereClause + `
		ORDER BY p.id DESC
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	
	args = append(args, params.Limit, (params.Page-1)*params.Limit)

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var catID *int
		var catName, catDesc *string

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &catID, &catName, &catDesc); err != nil {
			return nil, 0, err
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

	return products, totalItems, nil
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
