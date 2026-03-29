package postgres

import (
	"github.com/jackc/pgx/v5"
	"github.com/kuzuokatakumi/qube/internal/model"
)

func scanUsers(rows pgx.Rows) ([]*model.User, error) {
	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash,
			&u.Bio, &u.AvatarURL, &u.HeaderURL, &u.Location, &u.Website,
			&u.IsVerified, &u.IsPrivate, &u.FollowerCount, &u.FollowingCount, &u.PostCount,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
