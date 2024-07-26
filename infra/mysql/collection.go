package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
)

type collectionRepository struct {
	db SQLExecutor
}

func NewCollectionRepository(db *sql.DB) repository.CollectionRepository {
	return &collectionRepository{
		db: db,
	}
}

func (cr *collectionRepository) Get(ctx context.Context, id string) (*model.Collection, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Collections
	WHERE id = ?
	LIMIT 1`

	row := executor.QueryRowContext(ctx, query, id)

	var collection model.Collection
	if err := row.Scan(
		&collection.ID,
		&collection.Name,
		&collection.Rarity,
		&collection.Weight,
	); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (cr *collectionRepository) List(ctx context.Context) (model.Collections, error) {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}
	query := `SELECT *
	FROM Collections
	`

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections model.Collections
	for rows.Next() {
		var collection model.Collection
		if err = rows.Scan(
			&collection.ID,
			&collection.Name,
			&collection.Rarity,
			&collection.Weight,
		); err != nil {
			return nil, err
		}
		collections = append(collections, &collection)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return collections, nil
}

func (cr *collectionRepository) Create(ctx context.Context, collection model.Collection) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Collections (
	id, name, rarity, weight
	)
	VALUES (?, ?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		collection.ID,
		collection.Name,
		collection.Rarity,
		collection.Weight,
	); err != nil {
		return err
	}

	return nil
}

func (cr *collectionRepository) BatchCreate(ctx context.Context, collections model.Collections) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Collections (id, name, rarity, weight) VALUES `
	values := make([]interface{}, 0, len(collections)*4) //nolint:gomnd // 4 is the number of columns

	for i, collection := range collections {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?)"
		values = append(values, collection.ID, collection.Name, collection.Rarity, collection.Weight)
	}

	if _, err := executor.ExecContext(ctx, query, values...); err != nil {
		return err
	}
	return nil
}

func (cr *collectionRepository) Update(ctx context.Context, collection model.Collection) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `UPDATE Collections
	SET name = ?, rarity = ?, weight = ?
	WHERE id = ?
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		collection.Name,
		collection.Rarity,
		collection.Weight,
		collection.ID,
	); err != nil {
		return err
	}
	return nil
}

func (cr *collectionRepository) Delete(ctx context.Context, id string) error {
	executor := cr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM Collections
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
