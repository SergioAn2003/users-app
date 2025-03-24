package repository

import (
	"context"
	"errors"
	"fmt"
	"users-app/internal/entity"

	"github.com/gofrs/uuid/v5"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	sqlQuery := `
	select id, name, email, age, balance
	from users
	where id = $1`

	var user entity.User

	if err := r.pool.QueryRow(ctx, sqlQuery, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.Balance); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user with id %s not found", id)
		}

		return entity.User{}, fmt.Errorf("failed to get user with id %s: %w", id, err)
	}

	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user entity.User) error {
	constraintCode := "23505"

	sqlQuery := `
	insert into users
	(id, name, email, age, balance)
	values ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, sqlQuery, user.ID, user.Name, user.Email, user.Age, user.Balance)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == constraintCode {
			return fmt.Errorf("user with email %s %w", user.Email, entity.ErrAlreadyExists)
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user entity.User) error {
	sqlQuery := `
	update users
	set name = $2, email = $3, age = $4, balance = $5
	where id = $1`

	result, err := r.pool.Exec(ctx, sqlQuery, user.ID, user.Name, user.Email, user.Age, user.Balance)
	if err != nil {
		return fmt.Errorf("failed to update user with id %s: %w", user.ID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s %w", user.ID, entity.ErrNotFound)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	sqlQuery := `
	delete from users
	where id = $1`

	result, err := r.pool.Exec(ctx, sqlQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %s: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s %w", id, entity.ErrNotFound)
	}

	return nil
}
