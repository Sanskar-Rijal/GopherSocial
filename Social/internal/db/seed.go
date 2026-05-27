package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"social/internal/store"
)

func Seed(store store.Storage, db *sql.DB) {

	ctx := context.Background()

	//creating a user
	users := generateUser(100)

	tx, _ := db.BeginTx(ctx, nil)

	for _, value := range users {
		if err := store.Users.Create(ctx, tx, value); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating users: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, value := range posts {
		if err := store.Posts.Create(ctx, value); err != nil {
			log.Println("Error creating posts: ", err)
			return
		}
	}

	postComments := generateComments(400, posts, users)
	for _, value := range postComments {
		if err := store.Comments.Create(ctx, value); err != nil {
			log.Println("Error creating comments: ", err)
			return
		}
	}
	log.Printf("Seeding complete ........")
	return
}

func generateUser(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		//grabbing a random user
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	postComments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		//graabbing a random user
		user := users[rand.Intn(len(users))]
		//grabbing a random post
		post := posts[rand.Intn(len(posts))]

		postComments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return postComments
}
