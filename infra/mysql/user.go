package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
)

type userRepository struct {
	db SQLExecutor
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*model.User, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Users
	WHERE id = ?
	LIMIT 1`

	row := executor.QueryRowContext(ctx, query, id)

	var user model.User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Coins,
		&user.HighScore,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Create(ctx context.Context, user model.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Users (
	id, name, email, password, coins, high_score
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Coins,
		user.HighScore,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Update(ctx context.Context, user model.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `UPDATE Users
	SET name = ?, email = ?, password = ?, coins = ?, high_score = ?
	WHERE id = ?
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Coins,
		user.HighScore,
		user.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Delete(ctx context.Context, id string) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM Users
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
