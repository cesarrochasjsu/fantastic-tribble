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
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var forumFlag, genreFlag, favoritesFlag, authorFlag, reviewFlag, pendingFlag, commentFlag, postFlag bool
var title string
var articleId int

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list database contents",
	Long:  `List  information  about  the Database (manga titles by default).`,
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

		if reviewFlag {
			reviews, err := reviewsByTitle(title)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Reviews for %s found: %v\n", title, reviews)
		} else if commentFlag {
			article, err := checkArticleId(manga.Post{ArticleId: int64(articleId)})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Comments for %v found\n", article)
			comments, err := showComments(article)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Comments: %v found\n", comments)
		} else if authorFlag {
			authors, err := showAuthors()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Authors found: %v\n", authors)
		} else if forumFlag {
			forums, err := showForums()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Forums found: %v\n", forums)
		} else if genreFlag {
			genres, err := showGenres()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Genres found: %v\n", genres)
		} else if favoritesFlag {
			favorites, err := listByFavorite()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Most Popular: %v\n", favorites)
		} else if pendingFlag {
			requests, err := showRequests()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Requests found: %v\n", requests)
		} else {
			mangas, err := showMangas()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Mangas found: %v\n", mangas)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lsCmd.Flags().BoolVarP(&authorFlag, "authors", "a", false, "list all authors")
	lsCmd.Flags().BoolVarP(&forumFlag, "forums", "f", false, "list all forums")
	lsCmd.Flags().BoolVarP(&genreFlag, "genres", "g", false, "list all genres")
	lsCmd.Flags().BoolVarP(&favoritesFlag, "sort", "s", false, "list sorted by most favorited")
	lsCmd.Flags().BoolVarP(&pendingFlag, "requests", "R", false, "list manga requests by reviewers")
	lsCmd.Flags().BoolVarP(&reviewFlag, "review", "r", false, "list reviews by user (requires title of manga -t)")
	lsCmd.Flags().BoolVarP(&commentFlag, "comments", "c", false, "list comments (requires article_id -A)")
	lsCmd.Flags().StringVarP(&title, "title", "t", "", "title of a manga")
	lsCmd.Flags().IntVarP(&articleId, "ArticleId", "A", 0, "Id of an article")
	lsCmd.MarkFlagsRequiredTogether("review", "title")
	lsCmd.MarkFlagsRequiredTogether("comments", "ArticleId")
	// lsCmd.Flags().BoolVarP(&)
}

func showRequests() ([]manga.Requests, error) {
	var requests []manga.Requests

	rows, err := db.Query(`select request_id, reviewer_id, title from request`)
	if err != nil {
		return nil, fmt.Errorf("reviewsByTitle: %v", err)
	}
	for rows.Next() {
		var request manga.Requests
		if err := rows.Scan(&request.Request_id, &request.Reviewer_id, &request.Title); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return requests, nil
}

func showComments(post manga.Post) ([]manga.Comment, error) {
	var comments []manga.Comment

	rows, err := db.Query(`
select forum_id, article_id, reviewer_id, content 
from comment
where article_id = ?;`, post.ArticleId)
	if err != nil {
		return nil, fmt.Errorf("reviewsByTitle %q: %v", title, err)
	}
	for rows.Next() {
		var comment manga.Comment
		if err := rows.Scan(&comment.ForumId, &comment.ArticleId, &comment.ReviewerId, &comment.Content); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return comments, nil
}

func reviewsByTitle(title string) ([]manga.Review, error) {
	var reviews []manga.Review

	rows, err := db.Query(`SELECT review_id, reviewer_id, manga_id, title, description
FROM review join reviewer using(reviewer_id) WHERE title = ?`, title)
	if err != nil {
		return nil, fmt.Errorf("reviewsByTitle %q: %v", title, err)
	}
	for rows.Next() {
		var review manga.Review
		if err := rows.Scan(&review.ReviewId, &review.Manga_id, &review.Reviewer_id, &review.Title, &review.Description); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return reviews, nil

}

func showAuthors() ([]manga.Author, error) {
	var authors []manga.Author

	rows, err := db.Query("SELECT * FROM author")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var author manga.Author
		if err := rows.Scan(&author.A_ID, &author.A_name); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return authors, nil
}

func listByFavorite() ([]manga.Favorite, error) {
	var favorites []manga.Favorite

	rows, err := db.Query(`
select m.title, count(distinct reviewer_id) as favorites
FROM manga m JOIN favorite using(manga_ID)
GROUP BY manga_ID
ORDER BY count(distinct reviewer_id) desc;`)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var favorite manga.Favorite
		if err := rows.Scan(&favorite.Title, &favorite.Count); err != nil {
			return nil, fmt.Errorf("listByFavorite: %v", err)
		}
		favorites = append(favorites, favorite)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return favorites, nil
}

func showGenres() ([]manga.Genre, error) {
	var genres []manga.Genre

	rows, err := db.Query("SELECT * FROM genres")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var genre manga.Genre
		if err := rows.Scan(&genre.GId, &genre.GName); err != nil {
			return nil, fmt.Errorf("showGenre: %v", err)
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return genres, nil
}

func showForums() ([]manga.Forum, error) {
	var forums []manga.Forum

	rows, err := db.Query("SELECT * FROM forum")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var forum manga.Forum
		if err := rows.Scan(&forum.Forum_id, &forum.Title, &forum.Description); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		forums = append(forums, forum)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return forums, nil
}

func showMangas() ([]manga.Manga, error) {
	var mangas []manga.Manga

	rows, err := db.Query("SELECT * FROM manga")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var manga manga.Manga
		if err := rows.Scan(&manga.ID, &manga.Title, &manga.Description); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		mangas = append(mangas, manga)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return mangas, nil
}
