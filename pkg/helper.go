package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	// prov "Z0y6h0kS9X/libby-chapterizer/provider"
)

var (
	authorRegex   = regexp.MustCompile(`^aut(hor)?$`)
	narratorRegex = regexp.MustCompile(`^n(arrator|rt)?$`)
)

func FormatDuration(seconds float64) string {

	lengthRaw := time.Duration(seconds) * time.Second
	length := fmt.Sprintf("%02d:%02d:%02d.%03d",
		int(lengthRaw.Hours()),
		int(lengthRaw.Minutes())%60,
		int(lengthRaw.Seconds())%60,
		int(lengthRaw.Milliseconds()))

	length = strings.TrimRight(length, "0")

	return length

}

func GetComplexDuration(file1 string, file1Start, file2End float64) (string, error) {

	file1Duration, err := GetFileDuration(file1)
	if err != nil {
		log.Println("Error getting file1 duration:", err)
		return "", err
	}

	// Calculates duration using file duration and start time
	duration1 := file1Duration - file1Start

	// file 2 will always start as 0, so no need to get duration, it will be whatever file2End is
	duration2 := file2End

	// Adds the duration of the 2 file pieces together
	duration := duration1 + duration2

	// Formats the duration
	durationFormatted := FormatDuration(duration)

	return durationFormatted, nil
}

func GetSimpleDuration(start, end float64) (string, error) {

	lengthRaw := end - start
	length := FormatDuration(lengthRaw)

	return length, nil
}

func GetFileNameAndSeconds(path string) (string, float64) {
	fileName := ""
	seconds := 0.000

	// Discard everything up to and including 'Fmt425-'
	index := strings.Index(path, "Fmt425-")
	if index != -1 {
		path = path[index+len("Fmt425-"):]
	}

	// Split on '#', if it exists
	if strings.Contains(path, "#") {
		parts := strings.Split(path, "#")
		fileName = parts[0]
		seconds, _ = strconv.ParseFloat(parts[1], 64)
	} else {
		fileName = path
	}

	return fileName, seconds
}

func GetPartFromMp3File(mp3File string) string {

	part := ""
	index := (strings.Index(mp3File, "-Part") - 4)
	if index != -1 {
		part = mp3File[index+len("-Part"):]
	}

	return part
}

func GetBitRate(path string) (int, error) {
	// ffprobe -v error -show_entries format=bit_rate -of default=noprint_wrappers=1 input.mp3
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", path)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	// Trim whitespace from the output
	output = []byte(strings.TrimSpace(string(output)))

	// Parse the output to an integer
	bitRate, err := strconv.Atoi(string(output))
	if err != nil {
		fmt.Println("Error parsing bit rate:", err)
		return 0, err
	}

	return bitRate, nil
}

func NormalizeName(filename string) string {

	outputFileNormal := strings.Map(func(r rune) rune {
		switch {
		case r == '<' || r == '>' || r == ':' || r == '"' || r == '/' || r == '\\' || r == '|' || r == '?' || r == '*':
			return '-'
		default:
			return r
		}
	}, filename)

	return outputFileNormal

}

func GetOutputDirPath(details BookDetails, asin, outPath string) (string, error) {

	// Adds the first author, series name, and title to the path
	outName := details.Title
	if asin != "" {
		outName = outName + " (" + asin + ")"
	}

	if details.SeriesPrimary.Position != "" {
		floatNumber, err := strconv.ParseFloat(details.SeriesPrimary.Position, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
		}
		padded := fmt.Sprintf("%04.1f", floatNumber)
		outName = "[" + padded + "]. " + outName
	}

	// Normalizes the fields
	outName = NormalizeName(outName)
	seriesName := NormalizeName(details.SeriesPrimary.Name)
	author := NormalizeName(details.Authors[0].Name)

	outputDir := path.Join(outPath, author, seriesName, outName)

	return outputDir, nil

}

func GetBookDetailsNoASIN(book Openbook) (BookDetails, error) {

	details := BookDetails{}

	// Create an author object
	author := struct {
		Asin string `json:"asin,omitempty"`
		Name string `json:"name,omitempty"`
	}{
		Asin: "",
		Name: GetPrimaryAuthor(book),
	}

	// Create a narrator object
	narrator := struct {
		Name string `json:"name,omitempty"`
	}{
		Name: GetPrimaryNarrator(book),
	}

	// Set the details
	details.Authors = []struct {
		Asin string `json:"asin,omitempty"`
		Name string `json:"name,omitempty"`
	}{author}

	details.Narrators = []struct {
		Name string `json:"name,omitempty"`
	}{narrator}

	details.Title = book.Title.Main
	details.SeriesPrimary.Name = book.Title.Collection
	details.Subtitle = book.Title.Subtitle
	details.Description = book.Description.Full

	return details, nil

}

func GetPrimaryAuthor(book Openbook) string {

	var authors []string

	// Get the primary author
	for _, creator := range book.Creator {
		if authorRegex.MatchString(creator.Role) {
			authors = append(authors, creator.Name)
			continue
		}
	}

	// Return the first author
	if len(authors) == 0 {
		return ""
	} else {
		return authors[0]
	}

}

func GetPrimaryNarrator(book Openbook) string {

	var narrators []string

	// Get the primary narrator
	for _, creator := range book.Creator {
		if narratorRegex.MatchString(creator.Role) {
			narrators = append(narrators, creator.Name)
		}
	}

	// Return the first narrator
	if len(narrators) == 0 {
		return ""
	} else {
		return narrators[0]
	}

}

func GetTimeBreakdown(input int) (hours, minutes, seconds, milliseconds int) {

	// Calculate hours, minutes, seconds, and milliseconds
	hours = (input / (1000 * 60 * 60)) % 24
	minutes = (input / (1000 * 60)) % 60
	seconds = (input / 1000) % 60
	milliseconds = input % 1000

	// Returns the values
	return hours, minutes, seconds, milliseconds
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

func GenerateChapterBlock(file, title string, duration, lastTimeMS int) string {

	// Gets the start time of the file
	start := lastTimeMS
	end := start + duration

	// Creates the chapter block
	chapterBlock := fmt.Sprintf("[CHAPTER]\nTIMEBASE=1/1000\nSTART=%d\nEND=%d\ntitle=%s\n", start, end, title)

	// Returns the chapter block
	return chapterBlock

}

func JSONFileToOpenBook(jsonPath string) (Openbook, error) {

	// Creates a new Openbook
	openBook := Openbook{}

	// Imports JSON files
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return openBook, err
	}
	defer file.Close()

	// Reads the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return openBook, err
	}

	// Unmarshals JSON
	err = json.Unmarshal(data, &openBook)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return openBook, err
	}

	return openBook, nil

}

// func GetTitleFromFilename(filename string) string {

// 	// Discard everything up to and including 'Fmt425-'
// 	index := strings.Index(filename, "Fmt425-")
// 	title := ""
// 	if index != -1 {
// 		title = filename[index+len("Fmt425-"):]
// 	}

// 	return title

// }

// Needs rework
func CalculateDuration(milliseconds int) Duration {

	// Creates a new Duration object
	duration := Duration{}
	temp := time.Duration(milliseconds) * time.Millisecond

	// Sets duration properties
	duration.Hours = int(temp.Hours())
	duration.Minutes = int(temp.Minutes()) % 60
	duration.Seconds = int(temp.Seconds()) % 60
	duration.Milliseconds = int(temp.Nanoseconds()/int64(time.Millisecond)) % 1000
	duration.TotalMinutes = int(temp.Minutes())
	duration.TotalSeconds = int(temp.Seconds())
	duration.TotalMilliseconds = milliseconds

	// returns duration
	return duration
}
