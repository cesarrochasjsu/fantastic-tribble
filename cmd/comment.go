/*
Copyright © 2022 Steve Francia <spf@spf13.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"github.com/cesarrochasjsu/myapp/users"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func checkArticleId(post manga.Post) (manga.Post, error) {
	var temp manga.Post
	// Query for a value based on a single row.
	if err := db.QueryRow("select article_id, forum_id, reviewer_id from post where article_id = ?;", post.ArticleId).Scan(&temp.ArticleId, &temp.ForumId, &temp.ReviewerId); err != nil {
		if err == sql.ErrNoRows {
			return manga.Post{}, fmt.Errorf("%v: unknown article", post)
		}
		return manga.Post{}, fmt.Errorf("checkArticleId %v", post)
	}
	return temp, nil
}

func postNewComment(comment manga.Comment) (int64, error) {
	result, err := db.Exec("INSERT INTO comment (forum_id, article_id, reviewer_id, content) VALUES (?, ?, ?, ?)", comment.ForumId, comment.ArticleId, comment.ReviewerId, comment.Content)
	if err != nil {
		return 0, fmt.Errorf("postNewComment: %d %v", comment.ForumId, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("postNewComment: %v", err)
	}
	return id, nil
}

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment [article_id] [content]",
	Short: "Writes a comment into the article",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Capture connection properties.
		cfg := mysql.Config{
			User:                 os.Getenv("DBUSER"),
			Passwd:               os.Getenv("DBPASS"),
			Net:                  "tcp",
			Addr:                 "127.0.0.1:3306",
			DBName:               "mangalist",
			AllowNativePasswords: true,
		}
		// Get a database handle.
		var err error
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatal(err)
		}

		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
		fmt.Println("Connected!")
		reviewerId, err := getReviewerId(manga.User{
			Name:     u,
			Password: pw,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of %s exists: %v\n", u, reviewerId)
		articleId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		newPost, err := checkArticleId(manga.Post{
			ArticleId: int64(articleId),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of article exists: %v\n", newPost)
		commentId, err := postNewComment(manga.Comment{
			ArticleId:  newPost.ArticleId,
			ForumId:    newPost.ForumId,
			ReviewerId: newPost.ReviewerId,
			Content:    args[1],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of new comment: %v\n", commentId)
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commentCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
	commentCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	commentCmd.MarkFlagsRequiredTogether("username", "password")
}
