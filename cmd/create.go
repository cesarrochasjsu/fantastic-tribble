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
	// "golang.org/x/crypto/bcrypt"
)

var forumMode bool

func getModeratorId(user manga.User) (int64, error) {
	var reviewer_id int64
	var password string
	row := db.QueryRow(`SELECT moderator_id, password
				FROM user u, moderator m
				WHERE m.moderator_email = ? and m.user_id = u.user_id`, user.Email)
	if err := row.Scan(&reviewer_id, &password); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("getModeratorId %s: unknown user", user.Email)
		}
		return 0, fmt.Errorf("getModeratorId %s", user.Email)
	}
	if CheckPasswordHash(user.Password, password) {
		return reviewer_id, nil
	}
	return 0, fmt.Errorf("Wrong password")
}

func addForum(forum manga.Forum) (int64, error) {
	result, err := db.Exec("INSERT INTO forum (title, description) VALUES (?, ?)", forum.Title, forum.Description)
	if err != nil {
		return 0, fmt.Errorf("addForum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addForum: %v", err)
	}
	return id, nil
}

func addManga(manga manga.Manga) (int64, error) {
	result, err := db.Exec("INSERT INTO manga (title, description) VALUES (?, ?)", manga.Title, manga.Description)
	if err != nil {
		return 0, fmt.Errorf("addManga: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addManga: %v", err)
	}
	return id, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [title] [description]",
	Short: "Adds a manga or forum to the database",
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
		userId, err := getModeratorId(manga.User{
			Email:    u,
			Password: pw,
		})
		if err != nil {
			log.Fatal(err)
		}
		if forumMode {
			forumId, err := addForum(manga.Forum{
				Title:       args[0],
				Description: args[1],
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v added %v\n", userId, forumId)
		} else {
			mangaId, err := addManga(manga.Manga{
				Title:       args[0],
				Description: args[1],
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v added %v\n", userId, mangaId)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&u, "moderator", "m", "", "Moderator Email (required if password is set)")
	createCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
	createCmd.Flags().BoolVarP(&forumMode, "forum", "f", false, "Creates a new forum or board")
	createCmd.MarkFlagRequired("moderator")
	createCmd.MarkFlagsRequiredTogether("moderator", "password")
}
