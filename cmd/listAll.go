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

// albumsByArtist queries for albums that have the specified artist name.
func showAll() ([]manga.Manga, error) {
	// An albums slice to hold data from returned rows.
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

// listAllCmd represents the listAll command
var listAllCmd = &cobra.Command{
	Use:   "listAll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(0),
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
		mangas, err := showAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Albums found: %v\n", mangas)
	},
}

func init() {
	rootCmd.AddCommand(listAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
