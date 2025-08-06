package repository

import (
	"backend/internal/model"
	"backend/pkg/getusers"
	"database/sql"
	"fmt"
)

func GetPosts(id string, db *sql.DB) (*[]model.Post, error) {
	var posts []model.Post

	rows, err := db.Query(`SELECT id, user_id, title, content, visibility, post_image, created_at FROM posts
WHERE visibility = 'public'
   OR (
        visibility = 'almostprivate'
        AND EXISTS (
            SELECT 1
            FROM followers
            WHERE followers.followed_id = ?
              OR followers.follower_id = ?
              AND followers.status = 'accepted'
        )
	OR(
		visibility = 'private'
		AND EXISTS(
			SELECT 1 FROM private_posts 
			WHERE private_posts.post_id = posts.id
			AND private_posts.user_id = ?
		)

	)
    )`, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var post model.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Visibility, &post.ImageUrl, &post.CreatedAt); err != nil {
			fmt.Println(err.Error())

			return nil, err
		}
		user, err := getusers.GetUserByID(db, post.UserId)
		if err != nil {
			fmt.Println("User not found:", post.UserId)
		}
		post.Creator = user.FirstName + " " + user.LastName
		posts = append(posts, post)
	}
	return &posts, nil

}
