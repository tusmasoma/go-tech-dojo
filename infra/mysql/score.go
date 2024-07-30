package mysql

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
)

type scoreRepository struct {
	db SQLExecutor
}

func NewScoreRepository(db *sql.DB) repository.ScoreRepository {
	return &scoreRepository{
		db: db,
	}
}

func (sr *scoreRepository) Get(ctx context.Context, id string) (*model.Score, error) {
	executor := sr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Scores
	WHERE id = ?
	LIMIT 1`

	row := executor.QueryRowContext(ctx, query, id)

	var score model.Score
	if err := row.Scan(
		&score.ID,
		&score.UserID,
		&score.Value,
	); err != nil {
		return nil, err
	}
	return &score, nil
}

func (sr *scoreRepository) Create(ctx context.Context, score model.Score) error {
	executor := sr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Scores (id, user_id, value)
	VALUES (?, ?, ?)`

	if _, err := executor.ExecContext(
		ctx,
		query,
		score.ID,
		score.UserID,
		score.Value,
	); err != nil {
		return err
	}

	return nil
}

func (sr *scoreRepository) Delete(ctx context.Context, id string) error {
	executor := sr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM Scores
	WHERE id = ?`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
