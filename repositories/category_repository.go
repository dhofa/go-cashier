package repositories

import (
	"context"
	"errors"
	"go-cashier/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// =====================
// GET ALL
// =====================
func (repo *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name, COALESCE(description, '') FROM categories`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// =====================
// CREATE
// =====================
func (repo *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`

	return repo.db.
		QueryRow(ctx, query, category.Name, category.Description).
		Scan(&category.ID)
}

// =====================
// GET BY ID
// =====================
func (repo *CategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := `
		SELECT id, name, COALESCE(description, '')
		FROM categories
		WHERE id = $1
	`

	var c models.Category
	err := repo.db.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("kategori tidak ditemukan")
		}
		return nil, err
	}

	return &c, nil
}

// =====================
// UPDATE
// =====================
func (repo *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`

	tag, err := repo.db.Exec(
		ctx,
		query,
		category.Name,
		category.Description,
		category.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}

// =====================
// DELETE
// =====================
func (repo *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	tag, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}
