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
)

// addManga adds the specified album to the database,
// returning the album ID of the new entry
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

var db *sql.DB

// mangaCmd represents the manga command
var mangaCmd = &cobra.Command{
	Use:   "manga",
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
		mangaId, err := addManga(manga.Manga{
			Title:       args[0],
			Description: args[1],
		})
		err = doc.GenMarkdownTree(cmd, "/tmp")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID of added mangaum: %v\n", mangaId)
	},
}

func init() {
	rootCmd.AddCommand(mangaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mangaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mangaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
