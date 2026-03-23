package postgres

import (
	"context"
	"errors"
	"fmt"

	"cmpnion.space/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		insert into users (login, password_hash)
		values ($1, $2)
		returning id;
	`

	if err := r.pool.QueryRow(ctx, query, user.Login, user.PasswordHash).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrLoginTaken
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	const query = `
		select id, login, password_hash
		from users
		where login = $1;
	`

	var user domain.User
	if err := r.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by login: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const query = `
		select id, login, password_hash
		from users
		where id = $1;
	`

	var user domain.User
	if err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
