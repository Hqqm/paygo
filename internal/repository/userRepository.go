package repository

import (
	"context"

	"github.com/Hqqm/paygo/internal/entities"
	"github.com/Hqqm/paygo/internal/interfaces"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (userRep *UserRepository) AddUserInfoToAccount(ctx context.Context, user *entities.User) error {
	userUUID, err := uuid.FromString(user.ID)
	if err != nil {
		return err
	}

	query := `	
		INSERT INTO users(id, first_name, last_name, patronymic)
		VALUES (:id, :first_name, :last_name, :patronymic)
	`

	_, err = userRep.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         userUUID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"patronymic": user.Patronymic,
	})

	return err
}

func (userRep *UserRepository) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	user := &entities.User{}
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	err = userRep.db.Get(user, "SELECT * FROM users WHERE id=$1", userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
