package user

import (
	"context"
	"fmt"
	"go-transaction/repository"
)

type Service struct {
	store repository.Repository
}

func NewService(store repository.Repository) *Service {
	return &Service{store: store}
}

type User struct {
	ID        int64
	Username  string
	Firstname *string
	Activate  bool
	Role      string
}

func (s Service) CreateUser(req *User) (*User, error) {
	var result *User
	ctx := context.Background()

	// do some check

	// commit transaction
	if err := s.store.RunInTx(ctx, func(q repository.ExtendedQuerier) error {

		user, err := q.CustomCreateUser(ctx, repository.InsertUserParams{
			Username:  req.Username,
			Firstname: req.Firstname,
			Activated: req.Activate,
		})
		if err != nil {
			return fmt.Errorf("InsertUser: %w", err)
		}

		// simulate error
		role, err := q.GetRole(ctx, "regular")
		if err != nil {
			return fmt.Errorf("GetRole: %w", err)
		}

		if err := q.AssignUserRole(ctx, repository.AssignUserRoleParams{
			UserID: user.ID,
			RoleID: role.ID,
		}); err != nil {
			return fmt.Errorf("AssignUserRole: %w", err)
		}

		result = &User{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			Activate:  user.Activated,
			Role:      role.Name,
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return result, nil
}
