package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
)

type userCollectionRepository struct {
	db SQLExecutor
}

func NewUserCollectionRepository(db *sql.DB) repository.UserCollectionRepository {
	return &userCollectionRepository{
		db: db,
	}
}

func (uc *userCollectionRepository) List(ctx context.Context, userID string) ([]*model.UserCollection, error) {
	executor := uc.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM User_Collections
	WHERE user_id = ?
	`

	rows, err := executor.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCollections []*model.UserCollection
	for rows.Next() {
		var userCollection model.UserCollection
		if err = rows.Scan(
			&userCollection.UserID,
			&userCollection.CollectionID,
		); err != nil {
			return nil, err
		}
		userCollections = append(userCollections, &userCollection)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userCollections, nil
}

func (uc *userCollectionRepository) Get(ctx context.Context, userID, collectionID string) (*model.UserCollection, error) {
	executor := uc.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM User_Collections
	WHERE user_id = ? AND collection_id = ?
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, userID, collectionID)

	var userCollection model.UserCollection
	if err := row.Scan(
		&userCollection.UserID,
		&userCollection.CollectionID,
	); err != nil {
		return nil, err
	}
	return &userCollection, nil
}

func (uc *userCollectionRepository) Create(ctx context.Context, userCollection model.UserCollection) error {
	executor := uc.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO User_Collections (
	user_id, collection_id
	)
	VALUES (?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		userCollection.UserID,
		userCollection.CollectionID,
	); err != nil {
		return err
	}
	return nil
}

func (uc *userCollectionRepository) BatchCreate(ctx context.Context, userCollections []*model.UserCollection) error {
	executor := uc.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO User_Collections (user_id, collection_id) VALUES `
	values := make([]interface{}, 0, len(userCollections)*2) //nolint:gomnd // 2 is the number of columns

	for i, userCollection := range userCollections {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?)"
		values = append(values, userCollection.UserID, userCollection.CollectionID)
	}

	if _, err := executor.ExecContext(ctx, query, values...); err != nil {
		return err
	}
	return nil
}

func (uc *userCollectionRepository) Delete(ctx context.Context, userID, collectionID string) error {
	executor := uc.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM User_Collections
	WHERE user_id = ? AND collection_id = ?
	`

	if _, err := executor.ExecContext(ctx, query, userID, collectionID); err != nil {
		return err
	}
	return nil
}
