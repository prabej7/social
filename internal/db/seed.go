package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	db "github.com/prabej7/social/internal/db/seed"
	"github.com/prabej7/social/internal/store"
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

func Seed(storage *store.Storage) {
	ctx := context.Background()

	numOfUsers := 50
	users := generateUser(&numOfUsers)

	createdUsers := make([]*store.User, 0, len(users))

	for _, user := range users {
		if err := storage.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
		createdUsers = append(createdUsers, user)
	}

	posts := generatePosts(createdUsers)
	createdPosts := make([]*store.Post, len(users)*5)

	for _, post := range posts {
		if err := storage.Posts.Create(ctx, post); err != nil {
			log.Println("Error while creating post:", err)
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

func generateUser(num *int) []*store.User {
	*num = min(*num, len(usernames))

	users := make([]*store.User, 0, *num)

	for i := range *num {
		users = append(users, &store.User{
			UserName: usernames[i],
			Email:    fmt.Sprintf("%s@social.com", usernames[i]),
			Password: "Password@123",
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
