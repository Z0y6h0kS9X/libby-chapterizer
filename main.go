package main

import (
	p "Z0y6h0kS9X/libby-chapterizer/pkg"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

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
var outPath string

func init() {
	rootCmd.Flags().StringVarP(&jsonPath, "json", "j", "", "The path to the openbook.json file")
	rootCmd.Flags().StringVarP(&outPath, "out", "o", "", "The path to the directory you want to output the files to")
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

	// Gets the directory path from the json path
	var fileDir string
	if strings.Contains(jsonPath, "/") {
		fileDir = path.Dir(jsonPath)
	} else {
		jsonPath = strings.Replace(jsonPath, "\\", "/", -1)
		fileDir = path.Dir(jsonPath)
	}

	// Gets the path
	var outDir string

	if outPath != "" {
		outDir = strings.Replace(outPath, "\\", "/", -1)

	} else {
		outDir = fileDir
	}

	// Adds the first author, series name, and title to the path
	outputPath := path.Join(outDir, authors[0], book.Title.Collection, book.Title.Main)

	// fmt.Println("Directory:", fileDir)

	// Gets a list of all the .mp3 files in the fileDir
	files, err := os.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var mp3Files []string
	for _, file := range files {
		if path.Ext(file.Name()) == ".mp3" {
			// fmt.Println(path.Join(fileDir, file.Name()))
			mp3Files = append(mp3Files, path.Join(fileDir, file.Name()))
		}
	}

	// $end = ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $file.FullName

	totalDuration := 0.0
	for _, mp3File := range mp3Files {
		cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", mp3File)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println("Error running ffprobe command:", err)
			continue
		}
		durationStr := strings.TrimSpace(string(stdout))
		duration, err := strconv.ParseFloat(durationStr, 64)
		if err != nil {
			fmt.Println("Error parsing duration:", err)
			continue
		}
		totalDuration += duration
	}

	lengthRaw := time.Duration(totalDuration) * time.Second
	length := fmt.Sprintf("%02d:%02d:%02d.%03d",
		int(lengthRaw.Hours()),
		int(lengthRaw.Minutes())%60,
		int(lengthRaw.Seconds())%60,
		int(lengthRaw.Milliseconds())%1000)

	fmt.Println("Title:", book.Title.Main)
	fmt.Println("Series:", book.Title.Collection)
	fmt.Println("Author:", authorString)
	fmt.Println("Narrator:", narratorString)
	fmt.Println("Duration:", length)
	fmt.Println("============================")

	var ProcessBlock []p.Process

	// iterates through the book.nav.toc array and splits the path on "Fmt425-" and "#" to enumerate the seconds (if applicable) and match to the file paths above
	// for _, toc := range book.Nav.Toc {
	for i := 0; i < len(book.Nav.Toc); i++ {

		// Makes an empty Process object
		var process p.Process

		toc := book.Nav.Toc[i]
		part, seconds := p.GetFileNameAndSeconds(toc.Path)

		// Iterates through the mp3Files array and checks if it matches part
		for _, mp3File := range mp3Files {
			if strings.Contains(mp3File, part) {
				process.Source = mp3File
				continue
			}
		}

		// Gets the next file in the mp3Files array and checks if it matches the path of the toc
		if i < len(book.Nav.Toc)-1 {
			toc2 := book.Nav.Toc[i+1]
			part2, seconds2 := p.GetFileNameAndSeconds(toc2.Path)

			if part == part2 {

				process.End = seconds2

			}

		}

		// Adds the title, start time, and output path
		process.Title = toc.Title
		process.Start = seconds

		outputFileNormal := strings.Map(func(r rune) rune {
			switch {
			case r == '<' || r == '>' || r == ':' || r == '"' || r == '/' || r == '\\' || r == '|' || r == '?' || r == '*':
				return '-'
			default:
				return r
			}
		}, toc.Title)

		process.Output = path.Join(outputPath, outputFileNormal+".mp3")

		// // Adds the command
		// if process.End != 0 {
		// 	// No endtime specified, just start on starttime and go to end of file
		// 	process.CommandArgs = fmt.Sprintf("-ss %d -i '%s' -acodec copy -to %d '%s' -hide_banner -loglevel error", process.Start, process.Source, process.End, process.Output)
		// } else {
		// 	process.CommandArgs = fmt.Sprintf("-ss %d -i '%s' -acodec copy '%s' -hide_banner -loglevel error", process.Start, process.Source, process.Output)
		// }

		// Adds the process to the ProcessBlock
		ProcessBlock = append(ProcessBlock, process)
	}

	// Checks if the folder exists and creates it if it does not
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Runs the commands to generate the output files
	var m3u []string
	m3u = append(m3u, "#EXTM3U")
	m3u = append(m3u, fmt.Sprintf("#PLAYLIST: %s", book.Title.Main))

	for _, process := range ProcessBlock {

		_, file := path.Split(process.Output)
		// title := strings.Replace(file, ".mp3", "", -1)
		fmt.Println("Processing Chapter:", process.Title)

		if process.End != 0 {

			cmd := exec.Command("ffmpeg",
				"-i", process.Source,
				"-ss", fmt.Sprintf("%d", process.Start),
				"-acodec", "copy",
				"-to", fmt.Sprintf("%d", process.End),
				process.Output,
				"-hide_banner", "-loglevel", "error")

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		} else {

			cmd := exec.Command("ffmpeg",
				"-ss", fmt.Sprintf("%d", process.Start), "-i",
				process.Source, "-acodec", "copy",
				process.Output, "-hide_banner", "-loglevel", "error")

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		}

		m3u = append(m3u, fmt.Sprintf("#EXTINF:,%s\n%s", process.Title, file))

	}

	content := strings.Join(m3u, "\n")
	err = os.WriteFile(path.Join(outputPath, book.Title.Main+".m3u"), []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

}
