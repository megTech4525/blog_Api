package main

import (
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
		}

		// UserID := c.MustGet("UserID").int()

		var postID PostStructure

		query :=
			`INSERT INTO posts (user_id, title, content)
		VALUES($1,$2,$3)
		RETURNING id, created_at`

		if err := db.QueryRowx(query, posts.Title, posts.Contents).Scan(postID.ID, postID.CreatedAt); err != nil {
			c.JSON(401, gin.H{
				"ERROR": "Post not found",
			})
			return
		}

	}
}
