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
	"github.com/spf13/cobra/doc"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var moderator bool
var username, password string
var db *sql.DB

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// addUser adds the specified User to the database,
func addUser(user manga.User) (int64, error) {
	hash, _ := HashPassword(user.Password) // ignore error for the sake of simplicity
	result, err := db.Exec("INSERT INTO user (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hash)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	if moderator {
		reviewer_id, err := db.Exec("INSERT INTO moderator (user_id, moderator_email) VALUES (?, ?)", id, user.Email)
		if err != nil {
			return 0, fmt.Errorf("add Reviewer: %d %v", reviewer_id, err)
		}
	} else {
		reviewer_id, err := db.Exec("INSERT INTO reviewer (user_id, name) VALUES (?, ?)", id, user.Name)
		if err != nil {
			return 0, fmt.Errorf("add Reviewer: %d %v", reviewer_id, err)
		}
	}
	return id, nil
}

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register [name] [email] [password]",
	Short: "Registers a reviewer into the database",
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
		userId, err := addUser(manga.User{
			Name:     args[0],
			Email:    args[1],
			Password: args[2],
		})
		err = doc.GenMarkdownTree(cmd, "/tmp")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of added mangaum: %v\n", userId)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().BoolVarP(&moderator, "moderator", "m", false, "Create a moderator")
}
