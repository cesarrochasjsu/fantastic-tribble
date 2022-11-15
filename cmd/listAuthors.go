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

// listAuthorsCmd represents the listAuthors command
var listAuthorsCmd = &cobra.Command{
	Use:   "listAuthors",
	Short: "List all the authors we have on the database",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		authors, err := showAuthors()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Authors found: %v\n", authors)
	},
}

func init() {
	rootCmd.AddCommand(listAuthorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listAuthorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listAuthorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
