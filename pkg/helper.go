package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	authorRegex   = regexp.MustCompile(`^aut(hor)?$`)
	narratorRegex = regexp.MustCompile(`^n(arrator|rt)?$`)
)

type ChapterInfo struct {
	ID         int
	Title      string
	Start      int
	Duration   int
	FilePath   string
	FileLength Duration
}

// GetFileNameAndMilliseconds splits the path to extract the file name and milliseconds, if applicable.
// It takes a path string as input and returns the part string and milli integer.
func GetFileNameAndMilliseconds(path string) (string, int) {
	// Split the path to extract the file name and start time, if applicable.
	regex := `Fmt\d+-`
	leadingRegex := regexp.MustCompile(regex)
	if leadingRegex.MatchString(path) {
		regexMatch := leadingRegex.FindString(path)
		padCharsLength := utf8.RuneCountInString(regexMatch)
		path = path[strings.Index(path, regexMatch)+padCharsLength:]
	}

	// Split the Part from the start (Seconds).
	var part string
	var milli int
	tempStr := strings.Split(path, "#")
	if len(tempStr) > 1 {
		part = tempStr[0]
		tempFloat, err := strconv.ParseFloat(tempStr[1], 64)

		if err != nil {
			fmt.Println("Error Parsing to Float\n\n", err)
			return "", 0
		}
		tempInt := int(tempFloat * 1000)
		milli = tempInt
	} else {
		part = path
		milli = 0
	}

	return part, milli
}

// NormalizeName normalizes a filename by replacing special characters with hyphens.
func NormalizeName(filename string) string {
	// Use strings.Map to iterate over each rune in the filename
	// and replace special characters with hyphens
	outputFileNormal := strings.Map(func(r rune) rune {
		switch {
		// Replace special characters with hyphens
		case r == '<' || r == '>' || r == ':' || r == '"' || r == '/' || r == '\\' || r == '|' || r == '?' || r == '*':
			return '-'
		default:
			return r
		}
	}, filename)

	return outputFileNormal
}

// GetOutputDirPath generates the output directory path based on the metadata, ASIN, and the output path.
// It follows certain rules to format the output name and normalize the fields.
func GetOutputDirPath(meta Metadata, asin, outPath string) (string, error) {

	// Adds the first author, series name, and title to the path
	outName := meta.Title
	if asin != "" {
		outName = outName + " (" + asin + ")"
	}

	if meta.Series.Position != 0.0 {
		padded := fmt.Sprintf("%04.1f", meta.Series.Position)
		outName = "[" + padded + "]. " + outName
	}

	// Normalizes the fields
	outName = NormalizeName(outName)
	seriesName := NormalizeName(meta.Series.Name)
	author := NormalizeName(meta.Author)

	outputDir := path.Join(outPath, author, seriesName, outName)

	return outputDir, nil

}

// GetPrimaryAuthor returns the primary author of a book.
func GetPrimaryAuthor(book Openbook) string {

	var authors []string

	// Iterate through the creators of the book
	for _, creator := range book.Creator {
		// Check if the creator's role matches the author regular expression
		if authorRegex.MatchString(creator.Role) {
			// Add the creator's name to the list of authors
			authors = append(authors, creator.Name)
		}
	}

	// Return the first author in the list
	if len(authors) == 0 {
		return ""
	} else {
		return authors[0]
	}

}

// GetPrimaryNarrator returns the primary narrator of a book.
func GetPrimaryNarrator(book Openbook) string {
	// Create a slice to store the narrators
	var narrators []string

	// Iterate through the creators of the book
	for _, creator := range book.Creator {
		// Check if the creator has a role that matches the narrator regex
		if narratorRegex.MatchString(creator.Role) {
			// Add the creator's name to the narrators slice
			narrators = append(narrators, creator.Name)
		}
	}

	// Return the first narrator in the slice
	if len(narrators) == 0 {
		return ""
	} else {
		return narrators[0]
	}
}

// GetAllMp3Files returns a list of all the .mp3 files in the given directory and its subdirectories.
// The function takes a `path` parameter, which is the directory path to search for .mp3 files.
// It returns a slice of strings containing the absolute paths of the .mp3 files found, and an error if any.
func GetAllMp3Files(path string) ([]string, error) {

	var mp3Files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Check if the current file has a .mp3 extension
		if strings.HasSuffix(path, ".mp3") {
			mp3Files = append(mp3Files, filepath.ToSlash(path))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// return the list of .mp3 files
	return mp3Files, nil
}

// JSONFileToOpenBook reads a JSON file and returns an Openbook struct.
// It takes the path to the JSON file as input and returns the Openbook struct
// and an error if any.
func JSONFileToOpenBook(jsonPath string) (Openbook, error) {
	openBook := Openbook{} // Creates a new Openbook

	file, err := os.Open(jsonPath) // Opens the JSON file
	if err != nil {
		fmt.Println("Error opening file:", err)
		return openBook, err
	}
	defer file.Close()

	data, err := io.ReadAll(file) // Reads the file
	if err != nil {
		fmt.Println("Error reading file:", err)
		return openBook, err
	}

	err = json.Unmarshal(data, &openBook) // Unmarshals JSON
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return openBook, err
	}

	return openBook, nil
}

// CalculateDuration calculates the duration in hours, minutes, seconds, and milliseconds
// based on the given number of milliseconds.
func CalculateDuration(milliseconds int) Duration {
	// Create a new Duration object
	duration := Duration{}

	// Convert milliseconds to time.Duration
	temp := time.Duration(milliseconds) * time.Millisecond

	// Set duration properties
	duration.Hours = int(temp.Hours())
	duration.Minutes = int(temp.Minutes()) % 60
	duration.Seconds = int(temp.Seconds()) % 60
	duration.Milliseconds = int(temp.Nanoseconds()/int64(time.Millisecond)) % 1000
	duration.TotalMinutes = int(temp.Minutes())
	duration.TotalSeconds = int(temp.Seconds())
	duration.TotalMilliseconds = milliseconds

	// Return duration
	return duration
}
