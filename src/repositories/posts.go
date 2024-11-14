package repositories

import (
	"blogger/src/models"
	"database/sql"
	"errors"
)

type posts struct {
	db *sql.DB
}

// PostRepository creates a posts repository.
func PostRepository(db *sql.DB) *posts {
	return &posts{db: db}
}

// ListingPost returns a listing of posts for the authenticated user.
func (p posts) ListingPost(userID uint64) ([]models.Post, error) {
	var posts []models.Post

	records, err := p.db.Query(`
		SELECT DISTINCT p.*, u.nickname
		FROM posts p
		INNER JOIN users u
		ON p.user_id = u.id
		INNER JOIN followers f
		ON p.user_id = f.user_id
		WHERE u.id = ? OR f.follower_id = ?
		ORDER BY p.id DESC
	`, userID, userID)

	if err != nil {
		return posts, err
	}
	defer records.Close()

	for records.Next() {
		var post models.Post

		if err = records.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return posts, nil
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// CreatePost inserts a brand new post in the database.
func (p posts) CreatePost(post models.Post) (uint64, error) {
	record, err := p.db.Prepare("INSERT INTO posts(title, content, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer record.Close()

	result, err := record.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

// GetPost returns a specific post.
func (p posts) GetPost(postID uint64) (models.Post, error) {
	record, err := p.db.Query(`
		SELECT p.*, u.nickname 
		FROM posts p
		INNER JOIN users u
		ON p.user_id = u.id
		WHERE p.id = ?
	`, postID)

	var post models.Post

	if err != nil {
		return post, err
	}
	defer record.Close()

	if record.Next() {
		if err = record.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return post, err
		}
	}

	return post, nil
}

// UpdatePost updates a specific post.
func (p posts) UpdatePost(postID uint64, post models.Post) error {
	record, err := p.db.Prepare("UPDATE posts SET title = ?, content = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer record.Close()

	if _, err = record.Exec(post.Title, post.Content, postID); err != nil {
		return err
	}

	return nil
}

// DeletePost deletes a specific post.
func (p posts) DeletePost(postID uint64) error {
	statement, err := p.db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(postID); err != nil {
		return err
	}

	return nil
}

// ListingPost returns a listing of posts for the authenticated user.
func (p posts) ListingUserPosts(userID uint64) ([]models.Post, error) {
	var posts []models.Post

	records, err := p.db.Query(`
		SELECT p.*, u.nickname
		FROM posts p
		INNER JOIN users u
		ON p.user_id = u.id
		WHERE u.id = ?
		ORDER BY p.id DESC
	`, userID)

	if err != nil {
		return posts, err
	}
	defer records.Close()

	for records.Next() {
		var post models.Post

		if err = records.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return posts, nil
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return posts, errors.New("invalid user")
	}

	return posts, nil
}

// LikePost increases the like property of the post.
func (p posts) LikePost(postID uint64) error {
	statement, err := p.db.Prepare("UPDATE posts SET likes = likes + 1 WHERE id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(postID); err != nil {
		return err
	}

	return nil
}

// UnlikePost decreases the like property of the post.
func (p posts) UnlikePost(postID uint64) error {
	statement, err := p.db.Prepare(`
		UPDATE posts SET likes = 
		CASE WHEN likes > 0 THEN
			likes - 1
		ELSE
			0
		END
		WHERE id = ?
	`)

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(postID); err != nil {
		return err
	}

	return nil
}
