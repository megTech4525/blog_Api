package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Reg(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reg registerRequest

		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			log.Println(err.Error())
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password),
			bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			log.Println(err.Error())
			return
		}

		query :=
			`INSERT INTO users
		(name,email,password) 
		VALUES($1,$2,$3) 
		RETURNING id`
		var newID int
		err = db.QueryRowx(query, reg.Name, reg.Email, string(hashedPassword)).Scan(&newID)
		if err != nil {
			c.JSON(409, gin.H{"error": "email already taken"})
			log.Println(err.Error())
			return
		}

		c.JSON(201, gin.H{
			"name":  reg.Name,
			"email": reg.Email,
			// "created_at": reg.CreatedAt,
		})
	}

}

func login(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds loginRequest
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(400, gin.H{"error": "Unable to login"})
			fmt.Printf("log: %v\n", creds)
			return
		}

		var user User

		query := `SELECT id, email, password_hash FROM users WHERE email = $1`

		err := db.Get(&user, query, creds.Email)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid email or password"})
			return
		}
		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password_Hash), []byte(creds.Password)); err != nil {
			c.JSON(401, gin.H{"error": "invalid email or password"})
			return

		}
		token, err := generateToken(user.ID)
		if err != nil {
			c.JSON(409, gin.H{
				"error": "invalid token"})
			return
		}
		c.JSON(200, gin.H{
			"token": token,
		})

	}

}

func Createposts(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts createPostRequest
		if err := c.ShouldBindJSON(&posts); err != nil {
			c.JSON(400, gin.H{"error": ""})
			return
		}

		userID := c.MustGet("userID").(int)

		var post PostStructure

		query :=
			`INSERT INTO posts (user_id, title, content)
		VALUES($1,$2,$3)
		RETURNING id, created_at, updated_at`

		if err := db.QueryRowx(query, userID, posts.Title, posts.Contents).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			c.JSON(500, gin.H{
				"ERROR": "Failed to create Post",
			})
			return
		}
		c.JSON(200, gin.H{
			"ID":        post.ID,
			"UserID":    post.UserID,
			"Title":     posts.Title,
			"Content":   posts.Contents,
			"CreatedAt": post.CreatedAt,
			"UpdatedAt": post.UpdatedAt,
		})
	}
}

func updatePost(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		userID := c.MustGet("userID").(int)

		var req updatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		query := `UPDATE posts
		SET title = $1, content = $2, updated_at = now()
		WHERE id = $3 and user_id = $4
		RETURNING id, title, content, updated_at`

		var post PostStructure
		err := db.QueryRowx(query, req.Title, postID, userID).
			Scan(&post.ID, &post.Title, &post.Content, &post.UpdatedAt)
		if err != nil {
			// either the post doesn't exist, or it exists but isn't this user's

			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "post not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Faild to update post"})
			log.Println(err.Error())
			return
		}
		c.JSON(200, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": post.UpdatedAt})
	}
}
