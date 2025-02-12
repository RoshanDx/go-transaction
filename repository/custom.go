package repository

import "context"

func (q *Queries) CustomCreateUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRow(ctx, insertUser, arg.Username, arg.Firstname, arg.Activated)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Activated,
		&i.CreatedAt,
	)
	return i, err
}
