package repository

import (
	"backend/internal/model"
	"backend/pkg/getusers"
	"database/sql"
	"fmt"
)

func GetPosts(id string, db *sql.DB) (*[]model.Post, error) {
	var posts []model.Post

	rows, err := db.Query(`
		SELECT 
			p.id, p.user_id, p.title, p.content, p.visibility, p.post_image, p.created_at,
			(
				SELECT COUNT(1) FROM comments c WHERE c.post_id = p.id AND c.parent_id IS NULL
			) AS comment_count,
			COALESCE(SUM(CASE WHEN r.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
			COALESCE(SUM(CASE WHEN r.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count,
			COALESCE(MAX(CASE WHEN r.user_id = ? THEN r.type ELSE '' END), '') AS user_reaction
		FROM posts p
		LEFT JOIN reactions r ON p.id = r.post_id
		WHERE p.visibility = 'public' 
		   OR p.user_id = ?
		   OR (
				p.visibility = 'almostprivate'
				AND EXISTS (
					SELECT 1
					FROM followers
					WHERE (followers.followed_id = ? OR followers.follower_id = ?)
					  AND followers.status = 'accepted'
				)
			)
		   OR (
				p.visibility = 'private'
				AND EXISTS (
					SELECT 1 
					FROM private_posts 
					WHERE private_posts.post_id = p.id
					  AND private_posts.user_id = ?
				)
			)
		GROUP BY p.id
		ORDER BY p.created_at DESC`, id, id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var post model.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Visibility, &post.ImageUrl, &post.CreatedAt, &post.CommentCount, &post.LikeCount, &post.DislikeCount, &post.UserReaction); err != nil {
			fmt.Println(err.Error())

			return nil, err
		}
		user, err := getusers.GetUserByID(db, post.UserId)
		if err != nil {
			fmt.Println("User not found:", post.UserId)
		}
		post.Creator = user.FirstName + " " + user.LastName
		post.CreatorImg = user.ImgURL
		posts = append(posts, post)
	}
	return &posts, nil

}

func GetPostByID(db *sql.DB, postID string, userID string) (*model.Post, error) {
	var post model.Post

	row := db.QueryRow(`
		SELECT 
			p.id, p.user_id, p.title, p.content, p.visibility, p.post_image, p.created_at,
			(
				SELECT COUNT(1) FROM comments c WHERE c.post_id = p.id AND c.parent_id IS NULL
			) AS comment_count,
			COALESCE(SUM(CASE WHEN r.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
			COALESCE(SUM(CASE WHEN r.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count,
			COALESCE(MAX(CASE WHEN r.user_id = ? THEN r.type ELSE '' END), '') AS user_reaction
		FROM posts p
		LEFT JOIN reactions r ON p.id = r.post_id
		WHERE p.id = ?
		GROUP BY p.id
	`, userID, postID)

	if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Visibility, &post.ImageUrl, &post.CreatedAt, &post.CommentCount, &post.LikeCount, &post.DislikeCount, &post.UserReaction); err != nil {
		return nil, err
	}

	user, err := getusers.GetUserByID(db, post.UserId)
	if err != nil {
		fmt.Println("User not found:", post.UserId)
	}
	post.Creator = user.FirstName + " " + user.LastName
	post.CreatorImg = user.ImgURL

	return &post, nil
}

func GetPostOwnerID(db *sql.DB, postID string) (string, error) {
	var ownerID string
	err := db.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		return "", err
	}
	return ownerID, nil
}
