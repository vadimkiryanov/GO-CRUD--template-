package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
)

type RatingRepository interface {
	SetVote(userID, postID int, action domains.VoteAction) error
	GetUserVote(userID, postID int) (*domains.Vote, error)
	GetStats(postID int) (*domains.RatingStats, error)
	GetStatsWithUserVote(userID, postID int) (*domains.RatingStatsWithUserVote, error)
	DeleteVote(userID, postID int) error
}

type RatingPostgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *RatingPostgres {
	return &RatingPostgres{db: db}
}

func (r *RatingPostgres) SetVote(userID, postID int, action domains.VoteAction) error {
	voteType := string(action)
	query := fmt.Sprintln(`
        INSERT INTO votes (user_id, post_id, vote_type) 
        VALUES ($1, $2, $3) 
        ON CONFLICT (user_id, post_id) DO UPDATE SET vote_type = $3, created_at = NOW()`)

	_, err := r.db.Exec(query, userID, postID, voteType)

	return err
}

func (r *RatingPostgres) GetUserVote(userID, postID int) (*domains.Vote, error) {
	vote := &domains.Vote{}
	err := r.db.Select(vote, `
        SELECT id, user_id, post_id, vote_type FROM votes 
        WHERE user_id = $1 AND post_id = $2`, userID, postID)
	if err != nil {
		return nil, nil // none
	}
	return vote, err
}

func (r *RatingPostgres) GetStats(postID int) (*domains.RatingStats, error) {
	stats := &domains.RatingStats{}
	err := r.db.Get(stats, `
        SELECT
            COALESCE(likes.cnt, 0) as likes,
            COALESCE(dislikes.cnt, 0) as dislikes
        FROM (SELECT COUNT(*) as cnt FROM votes WHERE post_id = $1 AND vote_type = 'like') likes,
             (SELECT COUNT(*) as cnt FROM votes WHERE post_id = $1 AND vote_type = 'dislike') dislikes`, postID)
	stats.NetScore = stats.Likes - stats.Dislikes
	return stats, err
}

// GetStatsWithUserVote gets rating stats along with user's vote status for a specific post
func (r *RatingPostgres) GetStatsWithUserVote(userID, postID int) (*domains.RatingStatsWithUserVote, error) {
	// First get the rating stats
	stats := &domains.RatingStats{}
	err := r.db.Get(stats, `
        SELECT
            COALESCE(likes.cnt, 0) as likes,
            COALESCE(dislikes.cnt, 0) as dislikes
        FROM (SELECT COUNT(*) as cnt FROM votes WHERE post_id = $1 AND vote_type = 'like') likes,
             (SELECT COUNT(*) as cnt FROM votes WHERE post_id = $1 AND vote_type = 'dislike') dislikes`, postID)
	if err != nil {
		return nil, err
	}

	// Then get the user's vote
	var userVote string
	rows, err := r.db.Query(`
		SELECT vote_type
		FROM votes 
		WHERE post_id = $1 AND user_id = $2 
		LIMIT 1`, postID, userID)
	if err != nil {
		// If query fails, userVote will remain empty
		userVote = ""
	} else {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&userVote)
			if err != nil {
				userVote = ""
			}
		} else {
			// No vote found
			userVote = ""
		}
	}

	result := &domains.RatingStatsWithUserVote{
		Likes:    stats.Likes,
		Dislikes: stats.Dislikes,
		NetScore: stats.Likes - stats.Dislikes,
		UserVote: userVote,
	}

	return result, nil
}

func (r *RatingPostgres) DeleteVote(userID, postID int) error {
	query := fmt.Sprintln(`
        DELETE FROM votes WHERE user_id = $1 AND post_id = $2`)

	_, err := r.db.Exec(query, userID, postID)

	return err
}
