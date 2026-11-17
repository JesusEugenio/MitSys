package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/jackc/pgx/v5"
)

type UsersStore interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type userStore struct {
	db DBTX
}

func NewUserStore(db DBTX) UsersStore {
	return &userStore{db: db}
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `SELECT id, email, full_name, password, is_admin, is_cashier, is_delivery FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.FullName, &user.Password, &user.IsAdmin, &user.IsCashier, &user.IsDelivery)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("error getting user by email, not found: %w", err)
		}

		return user, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}
