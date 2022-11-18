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

func checkForumId(forum manga.Forum) (int64, error) {
	var forumId int64
	// Query for a value based on a single row.
	if err := db.QueryRow("select forum_id from forum where forum_id = ?;", forum.Forum_id).Scan(&forumId); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: unknown manga", forumId)
		}
		return 0, fmt.Errorf("getMangaId %s", forumId)
	}
	return forumId, nil
}

func insertPost(post manga.Post) (manga.Post, error) {
	_, err := db.Exec("INSERT INTO post (article_id, forum_id, reviewer_id) VALUES (?, ?, ?)", post.ArticleId, post.ForumId, post.ReviewerId)
	if err != nil {
		return manga.Post{}, fmt.Errorf("insertPost: %v %v", post, err)
	}
	return post, nil
}

func postNewArticle(post manga.Article) (int64, error) {
	result, err := db.Exec("INSERT INTO forum_article (forum_id, title, content) VALUES (?, ?, ?)", post.ForumId, post.Title, post.Content)
	if err != nil {
		return 0, fmt.Errorf("postNew: %d %v", post.ForumId, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// postNewCmd represents the postNew command
var postNewCmd = &cobra.Command{
	Use:   "postNew [ForumId][title] [content]",
	Short: "Inserts a post into the database",
	Args:  cobra.ExactArgs(3),
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
		forumId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		verifiedForumId, err := checkForumId(manga.Forum{
			Forum_id: int64(forumId),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of forum exists: %v\n", verifiedForumId)
		postId, err := postNewArticle(manga.Article{
			ForumId: verifiedForumId,
			Title:   args[1],
			Content: args[2],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of new post: %v\n", postId)
		newPost, err := insertPost(manga.Post{
			ArticleId:  postId,
			ForumId:    verifiedForumId,
			ReviewerId: reviewerId,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successful Post: %v\n", newPost)
	},
}

func init() {
	rootCmd.AddCommand(postNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postNewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postNewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	postNewCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
	postNewCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	postNewCmd.MarkFlagsRequiredTogether("username", "password")
}
