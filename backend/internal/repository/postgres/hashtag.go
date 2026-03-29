package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HashtagRepo struct {
	pool *pgxpool.Pool
}

func NewHashtagRepo(pool *pgxpool.Pool) *HashtagRepo {
	return &HashtagRepo{pool: pool}
}

type Hashtag struct {
	Name      string
	PostCount int
}

func (r *HashtagRepo) UpsertAndLink(ctx context.Context, postID uuid.UUID, tags []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, tag := range tags {
		var hashtagID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO hashtags (id, name, post_count) VALUES ($1, $2, 1)
			ON CONFLICT (name) DO UPDATE SET post_count = hashtags.post_count + 1
			RETURNING id`,
			uuid.New(), tag,
		).Scan(&hashtagID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO post_hashtags (post_id, hashtag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			postID, hashtagID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *HashtagRepo) GetTrending(ctx context.Context, limit int) ([]Hashtag, error) {
	query := `SELECT name, post_count FROM hashtags ORDER BY post_count DESC LIMIT $1`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Hashtag
	for rows.Next() {
		h := Hashtag{}
		if err := rows.Scan(&h.Name, &h.PostCount); err != nil {
			return nil, err
		}
		tags = append(tags, h)
	}
	return tags, nil
}

func (r *HashtagRepo) GetPostHashtags(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	query := `
		SELECT ph.post_id, h.name
		FROM post_hashtags ph
		JOIN hashtags h ON h.id = ph.hashtag_id
		WHERE ph.post_id = ANY($1)`
	rows, err := r.pool.Query(ctx, query, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]string)
	for rows.Next() {
		var postID uuid.UUID
		var name string
		if err := rows.Scan(&postID, &name); err != nil {
			return nil, err
		}
		result[postID] = append(result[postID], name)
	}
	return result, nil
}
