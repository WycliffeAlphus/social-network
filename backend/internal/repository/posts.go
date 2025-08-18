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
    ) AS comment_count
FROM posts p
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
ORDER BY p.created_at DESC`, id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var post model.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Visibility, &post.ImageUrl, &post.CreatedAt, &post.CommentCount); err != nil {
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
