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
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	// "golang.org/x/crypto/bcrypt"
)

func checkMangaId(id int64) (int64, error) {
	var manga_id int64
	// Query for a value based on a single row.
	if err := db.QueryRow("select manga_id from manga where manga_id=?", id).Scan(&manga_id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: unknown manga", id)
		}
		return 0, fmt.Errorf("getMangaId %s", id)
	}
	return manga_id, nil
}

func markFavorite(reviewerId, mangaId int64) ([]int64, error) {
	_, err := db.Exec("INSERT INTO favorite (reviewer_id, manga_id) VALUES (?, ?)", reviewerId, mangaId)
	if err != nil {
		return []int64{0, 0}, fmt.Errorf("markFavorite: %v", err)
	}
	return []int64{reviewerId, mangaId}, nil
}

// favoriteCmd represents the favorite command
var favoriteCmd = &cobra.Command{
	Use:   "favorite [manga_id]",
	Short: "Favorite a manga",
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
		reviewerId, err := getReviewerId(manga.User{
			Name:     u,
			Password: pw,
		})
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		verified_id, err := checkMangaId(int64(id))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reviewer number: %d\n", verified_id)
		markFavorite(reviewerId, int64(id))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reviewer number: %s Marked Favorite by %s\n", args[0], u)
	},
}

func init() {
	rootCmd.AddCommand(favoriteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// favoriteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// favoriteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	favoriteCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
	favoriteCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	createCmd.MarkFlagRequired("username")
	favoriteCmd.MarkFlagsRequiredTogether("username", "password")
}
