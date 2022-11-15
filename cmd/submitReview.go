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
	"golang.org/x/crypto/bcrypt"
)

var u, pw string

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getUserId returns the reviewer id
func getUserId(user manga.User) (int64, error) {
	var user_id int64
	var password string
	row := db.QueryRow("SELECT user_id, password from user WHERE name = ?", user.Name)
	if err := row.Scan(&user_id, &password); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("submitReview %s: unknown user", user.Name)
		}
		return 0, fmt.Errorf("submitReview %s", user.Name)
	}
	if CheckPasswordHash(user.Password, password) {
		return user_id, nil
	}
	return 0, fmt.Errorf("Wrong password")
}

// getUserId returns the reviewer id
func postReview(userID int64, review manga.Review) (int64, error) {
	result, err := db.Exec("INSERT INTO review (user_id, title, description) VALUES (?, ?, ?)", userID, review.Title, review.Description)
	if err != nil {
		return 0, fmt.Errorf("postReview: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("postReview: %v", err)
	}
	return id, nil
}

// submitReviewCmd represents the submitReview command
var submitReviewCmd = &cobra.Command{
	Use:   "submitReview",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
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
		userId, err := getUserId(manga.User{
			Name:     u,
			Password: pw,
		})
		if err != nil {
			log.Fatal(err)
		}
		reviewId, err := postReview(userId, manga.Review{
			Title:       args[0],
			Description: args[1],
		})
		fmt.Printf("ID of review: %d\n", reviewId)
	},
}

func init() {
	rootCmd.AddCommand(submitReviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitReviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// submitReviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	submitReviewCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
	submitReviewCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	submitReviewCmd.MarkFlagsRequiredTogether("username", "password")
}
