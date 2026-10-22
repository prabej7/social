package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	db "github.com/prabej7/social/internal/db/seed"
	"github.com/prabej7/social/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var usernames []string = []string{
	"alice", "bob", "charlie", "david", "emma",
	"olivia", "liam", "noah", "sophia", "ethan",
	"ava", "lucas", "mia", "james", "isabella",
	"benjamin", "amelia", "elijah", "harper", "logan",
	"prabej", "rohan", "suman", "nischal", "sanjay",
	"aarav", "anisha", "sujan", "bibek", "sabin",
	"admin", "moderator", "guest", "developer", "tester",
	"john_doe", "jane_doe", "alex_smith", "sam_wilson", "chris_brown",
	"user_001", "user_002", "user_003", "user_004", "user_005",
	"nepal_dev", "go_programmer", "code_master", "tech_lover", "debugger",
}

func Seed(storage *store.Storage, db *sql.DB) {
	ctx := context.Background()

	numOfUsers := 50
	users := generateUser(&numOfUsers)

	createdUsers := make([]*store.User, 0, len(users))

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Println("Error while initializing transaction:", err)
		return
	}

	for _, user := range users {
		if err := storage.Users.Create(ctx, tx, user); err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Println("Rollback error on creating user:", err)
			}
			log.Println("Error creating user:", err)
			return
		}
		createdUsers = append(createdUsers, user)
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing users:", err)
		return
	}

	posts := generatePosts(createdUsers)
	createdPosts := make([]*store.Post, 0, len(posts))

	for _, post := range posts {
		if err := storage.Posts.Create(ctx, post); err != nil {
			log.Println("Error while creating post:", err)
			continue
		}

		// Only append posts that were successfully created.
		if post.ID == 0 {
			log.Println("Post created without a valid ID")
			continue
		}

		createdPosts = append(createdPosts, post)
	}

	comments := generateComments(createdPosts, createdUsers)
	for _, comment := range comments {
		if err := storage.Comments.Create(ctx, comment); err != nil {
			log.Println("Error while creating comment:", err)
			return
		}
	}

	log.Println("Seeding completed successfully!")

}

func HashPassword(text string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func generateUser(num *int) []*store.User {
	*num = min(*num, len(usernames))

	users := make([]*store.User, 0, *num)

	for i := range *num {
		hashedPassword, err := HashPassword("Password@123")
		if err != nil {
			log.Println("Error while hashing password:", err)
			panic(err)
		}
		users = append(users, &store.User{
			UserName: usernames[i],
			Email:    fmt.Sprintf("%s@social.com", usernames[i]),
			Password: string(hashedPassword),
			RoleID:   1,
		})
	}

	return users
}

func generatePosts(users []*store.User) []*store.Post {
	if len(users) == 0 {
		return nil
	}

	postsPerUser := min(5, len(db.Posts))
	generatedPosts := make([]*store.Post, 0, postsPerUser*len(users))

	for i, user := range users {
		start := i * postsPerUser
		end := min(start+postsPerUser, len(db.Posts))

		for _, post := range db.Posts[start:end] {
			generatedPosts = append(generatedPosts, &store.Post{
				Title:   post.Title,
				Content: post.Content,
				Tags:    post.Tags,
				UserID:  user.ID,
			})
		}
	}

	return generatedPosts
}

func generateComments(posts []*store.Post, users []*store.User) []*store.Comment {
	if len(users) == 0 {
		return nil
	}

	if len(db.Comments) == 0 {
		return nil
	}

	comments := make([]*store.Comment, 0, len(posts))

	for _, post := range posts {
		if post == nil {
			continue
		}

		user := users[rand.Intn(len(users))]
		if user == nil {
			continue
		}

		comments = append(comments, &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: db.Comments[rand.Intn(len(db.Comments))],
		})
	}

	return comments
}
