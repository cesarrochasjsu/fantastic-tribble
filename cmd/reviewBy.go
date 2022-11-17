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

	"github.com/spf13/cobra"
)

func reviewByUser(name string) ([]manga.Review, error) {
	var reviews []manga.Review

	rows, err := db.Query(`SELECT reviewer_id, manga_id, title, description
FROM review join reviewer using(reviewer_id) WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("reviewsByUser %q: %v", name, err)
	}
	for rows.Next() {
		var review manga.Review
		if err := rows.Scan(&review.Manga_id, &review.Reviewer_id, &review.Title, &review.Description); err != nil {
			return nil, fmt.Errorf("showAll: %v", err)
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return reviews, nil
}

// reviewByCmd represents the reviewBy command
var reviewByCmd = &cobra.Command{
	Use:   "reviewBy [user]",
	Short: "List the reviews from a user",
	Args:  cobra.ExactArgs(1),
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
		reviews, err := reviewByUser(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reviews found: %v\n", reviews)
	},
}

func init() {
	rootCmd.AddCommand(reviewByCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reviewByCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reviewByCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
