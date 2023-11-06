package main

import (
	p "Z0y6h0kS9X/libby-chapterizer/pkg"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "libby-chapterizer",
	Short: "A brief description of your application",
	Long:  "A longer description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		// This function will be executed when your application is run
		fmt.Println("Hello, World!")
	},
}

var jsonPath string

func init() {
	rootCmd.Flags().StringVarP(&jsonPath, "json", "j", "", "The path to the openbook.json file")
}

func main() {
	fmt.Println("Hello, World!")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Imports JSON files
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var book p.Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	var authors []string
	var narrators []string

	for _, creator := range book.Creator {
		if creator.Role == "author" {
			authors = append(authors, creator.Name)
		} else if creator.Role == "narrator" {
			narrators = append(narrators, creator.Name)
		}
	}

	authorString := strings.Join(authors, ", ")
	narratorString := strings.Join(narrators, ", ")

	fmt.Println("Title:", book.Title.Main)
	fmt.Println("Series:", book.Title.Collection)
	fmt.Println("Author:", authorString)
	fmt.Println("Narrator:", narratorString)

	// Gets the directory path from the json path
	var fileDir string
	if strings.Contains(jsonPath, "/") {
		fileDir = path.Dir(jsonPath)
	} else {
		jsonPath = strings.Replace(jsonPath, "\\", "/", -1)
		fileDir = path.Dir(jsonPath)
	}

	fmt.Println("Directory:", fileDir)

	// Gets a list of all the .mp3 files in the fileDir
	files, err := os.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var mp3Files []string
	for _, file := range files {
		if path.Ext(file.Name()) == ".mp3" {
			fmt.Println(path.Join(fileDir, file.Name()))
			mp3Files = append(mp3Files, path.Join(fileDir, file.Name()))
		}
	}

	// iterates through the book.nav.toc array and splits the path on "Fmt425-" and "#" to enumerate the seconds (if applicable) and match to the file paths above
	// for _, toc := range book.Nav.Toc {
	for i := 0; i < len(book.Nav.Toc); i++ {
		toc := book.Nav.Toc[i]
		part, seconds := p.GetFileNameAndSeconds(toc.Path)

		// Gets the next file in the mp3Files array and checks if it matches the path of the toc
		if i != len(book.Nav.Toc)-1 {
			toc2 := book.Nav.Toc[i+1]
			part2, seconds2 := p.GetFileNameAndSeconds(toc2.Path)

			if part == part2 {
				fmt.Println("Chapter:", toc.Title, "File:", part, "Duration:", seconds, "-", seconds2, "Seconds")
				continue
			} else {
				fmt.Println("Chapter:", toc.Title, "File:", part, "Starts At:", seconds, "Seconds")
			}
		}

	}

}
