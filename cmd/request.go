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

func request(reviewer_id int64, title string) error {
	_, err := db.Exec("INSERT INTO request (reviewer_id, title) VALUES (?, ?)", reviewer_id, title)
	if err != nil {
		return fmt.Errorf("request: %v", err)
	}
	return nil
}

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request [title]",
	Short: "Requests for title to be added into database",
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
		reviwerId, err := getReviewerId(manga.User{
			Name:     u,
			Password: pw,
		})
		if err != nil {
			log.Fatal(err)
		}
		requestErr := request(reviwerId, args[0])
		if err != nil {
			log.Fatal(requestErr)
		}
		fmt.Printf("Request for %s successfully sent to be approved by moderator team\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	requestCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
	requestCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	requestCmd.MarkFlagsRequiredTogether("username", "password")
}
