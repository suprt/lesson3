package main

import (
	"fmt"
	"sync"
	"time"
)

type Comment struct {
	UserID int
	Text   string
}

type Session struct {
	ID string
}

type DB struct {
	comment []Comment
	session Session
}

var (
	db     *DB
	dbOnce sync.Once
)

func getDB() *DB {
	dbOnce.Do(func() {
		fmt.Println("loading DB...")
		time.Sleep(500 * time.Millisecond)

		db = &DB{
			comment: []Comment{
				{UserID: 1, Text: "hello world"},
				{UserID: 2, Text: "hi"},
				{UserID: 1, Text: "goodbye world"},
			},
			session: Session{
				ID: "session-121",
			},
		}
	})
	return db
}

func loadComms() <-chan []Comment {

	out := make(chan []Comment)
	go func() {
		defer close(out)
		db := getDB()
		fmt.Println("loading comments...")
		time.Sleep(500 * time.Millisecond)
		out <- db.comment
	}()

	return out
}

func loadUsers(commentsCh <-chan []Comment) {
	fmt.Println("loading users...")

	loaded := make(map[int]struct{})
	for comments := range commentsCh {
		for _, comment := range comments {
			if _, exists := loaded[comment.UserID]; exists {
				continue
			}
			loaded[comment.UserID] = struct{}{}

			fmt.Printf("user %d loaded\n", comment.UserID)
			time.Sleep(500 * time.Millisecond)
		}

	}
}

func loadSession() <-chan Session {

	out := make(chan Session)
	go func() {
		defer close(out)
		db := getDB()
		fmt.Println("loading session...")
		time.Sleep(500 * time.Millisecond)
		out <- db.session
	}()
	return out
}

func loadAtatach(sessionID string) {
	fmt.Printf("loading attachments to session %s\n", sessionID)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Attachments loaded")
}

func main() {
	var wg sync.WaitGroup

	comments := loadComms()
	session := loadSession()

	wg.Add(1)

	go func() {
		defer wg.Done()
		loadUsers(comments)
	}()

	ses := <-session
	if ses.ID != "" {
		wg.Add(1)

		go func() {
			defer wg.Done()
			loadAtatach(ses.ID)
		}()
	}
	wg.Wait()
	fmt.Println("done")
}
