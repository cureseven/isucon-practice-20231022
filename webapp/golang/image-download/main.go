package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strconv"
)

type Post struct {
	ID      int    `db:"id"`
	Imgdata []byte `db:"imgdata"`
	Mime    string `db:"mime"`
}

func main() {
	host := os.Getenv("ISUCONP_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ISUCONP_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCONP_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCONP_DB_PASSWORD")
	dbname := os.Getenv("ISUCONP_DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	defer db.Close()

	posts := make([]Post, 0)
	err = db.Select(&posts, "SELECT id, imgdata, mime FROM `posts`")
	if err != nil {
		log.Print(err)
		return
	}

	for _, post := range posts {
		ext := "jpg"
		if post.Mime == "image/png" {
			ext = "png"
		}
		if post.Mime == "image/gif" {
			ext = "gif"
		}

		filename := fmt.Sprintf("images/%d.%s", post.ID, ext)
		if err := os.WriteFile(filename, post.Imgdata, 0666); err != nil {
			log.Print(err)
			return
		}
	}
}
