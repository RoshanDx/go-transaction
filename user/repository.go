package user

import (
	"context"
	"fmt"
	"go-transaction/repository"
)

type Repository struct {
	store repository.Repository
}

func NewUserStore(store repository.Repository) Repository {
	return Repository{store}
}

func (s Repository) CreateUser(ctx context.Context, user *User) error {

	if err := s.store.RunInTx(ctx, func(q repository.ExtendedQuerier) error {

		user, err := q.CustomCreateUser(ctx, repository.InsertUserParams{
			Username:  user.Username,
			Firstname: user.Firstname,
			Activated: user.Activate,
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

		return nil
	}); err != nil {
		return err
	}

	return nil
}
