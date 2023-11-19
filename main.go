package main

import (
	p "Z0y6h0kS9X/libby-chapterizer/pkg"
	prov "Z0y6h0kS9X/libby-chapterizer/provider"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "libby-chapterizer",
	Short: "A brief description of your application",
	Long:  "A longer description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		// This function will be executed when your application is run
		// fmt.Println("Hello, World!")
	},
}

var jsonPath string
var outPath string
var test bool

func init() {
	rootCmd.Flags().StringVarP(&jsonPath, "json", "j", "", "The path to the openbook.json file")
	rootCmd.Flags().StringVarP(&outPath, "out", "o", "", "The path to the directory you want to output the files to")
	rootCmd.Flags().BoolVarP(&test, "test", "t", false, "Test mode")
}

func main() {
	// fmt.Println("Hello, World!")

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

	var book p.Openbook
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

	// Gets the ASIN
	asin, err := prov.GetBook(book.Title.Main, authorString, narratorString)
	if err != nil {
		fmt.Println("Error getting book:", err)
		return
	}

	// Gets the book details
	details, err := prov.GetBookDetails(asin)
	if err != nil {
		fmt.Println("Error getting book details:", err)
		return
	}

	// Gets the path
	var outDir string

	if outPath != "" {
		outDir = strings.Replace(outPath, "\\", "/", -1)

	} else {
		outDir = fileDir
	}

	// Adds the first author, series name, and title to the path
	outName := book.Title.Main
	if asin != "" {
		outName = outName + " (" + asin + ")"
	}

	if details.SeriesPrimary.Name != "" {
		floatNumber, err := strconv.ParseFloat(details.SeriesPrimary.Position, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
			return
		}
		padded := fmt.Sprintf("%04.1f", floatNumber)
		outName = "[" + padded + "]. " + outName
	}

	outputPath := path.Join(outDir, authors[0], book.Title.Collection, outName)

	// Gets a list of all the .mp3 files in the fileDir
	files, err := os.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	mp3FileMap := make(map[string]string)
	// var mp3Files []string
	totalDuration := 0.000

	for _, file := range files {
		if path.Ext(file.Name()) == ".mp3" {

			fullPath := path.Join(fileDir, file.Name())
			// mp3Files = append(mp3Files, fullPath)
			part := p.GetPartFromMp3File(fullPath)
			mp3FileMap[part] = fullPath

			duration, err := p.GetFileDuration(fullPath)
			if err != nil {
				fmt.Println("Error getting file duration:", err)
				return
			}
			totalDuration += duration
		}
	}

	duration := p.GetDurationFormatted(totalDuration)

	fmt.Println("============================")
	fmt.Println("Audiobook Information")
	fmt.Println("----------------------------")
	fmt.Println("Title:", book.Title.Main)
	fmt.Println("ASIN:", asin)
	if details.SeriesPrimary.Name != "" {
		fmt.Println("Series:", details.SeriesPrimary.Name)
		fmt.Println("Book Number:", details.SeriesPrimary.Position)
	}
	fmt.Println("Author:", authorString)
	fmt.Println("Narrator:", narratorString)
	fmt.Println("Duration:", duration)
	fmt.Println("============================")

	// if test {

	// 	fmt.Println("Testing...")
	// 	// "F:\Music\Audiobook\Testing\origin_full.mp3"
	// 	// testFile := "F:\\Music\\Audiobook\\Testing\\origin_full.mp3"

	// 	asin, err := prov.GetBook(book.Title.Main, authorString, narratorString)
	// 	if err != nil {
	// 		fmt.Println("Error getting book:", err)
	// 		return
	// 	}

	// 	// fmt.Println("ASIN:", asin)

	// 	// details, err := prov.GetBookDetails(asin)
	// 	// if err != nil {
	// 	// 	fmt.Println("Error getting book details:", err)
	// 	// 	return
	// 	// }

	// 	// fmt.Printf("Book %s in the %s series\n", details.SeriesPrimary.Position, details.SeriesPrimary.Name)

	// 	chapters, err := prov.GetChapters(asin)
	// 	if err != nil {
	// 		fmt.Println("Error getting chapters:", err)
	// 		return
	// 	}

	// 	fullFile := p.GetTestFile()

	// 	for i, chapter := range chapters.Chapters {

	// 		padded := fmt.Sprintf("%02d", i)

	// 		fileName := fmt.Sprintf("[%s]. %s.mp3", padded, p.NormalizeFileName(chapter.Title))
	// 		output := path.Join("F:/Music/Audiobook/Testing", fileName)

	// 		duration := p.ConvertMilliseconds(chapter.LengthMs)
	// 		fmt.Println("Processing Chapter:", padded, " - ", chapter.Title, " ("+p.GetDurationFormatted(duration)+")")

	// 		cmd := exec.Command("ffmpeg",
	// 			"-i", fullFile,
	// 			"-ss", fmt.Sprintf("%f", p.ConvertMilliseconds(chapter.StartOffsetMs)),
	// 			"-t", fmt.Sprintf("%f", p.ConvertMilliseconds(chapter.LengthMs)),
	// 			"-acodec", "copy",
	// 			output,
	// 			"-hide_banner", "-loglevel", "error")

	// 		err := cmd.Run()
	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 			return
	// 		}

	// 	}

	// 	// fmt.Println("Chapters:", len(chapters.Chapters))

	// 	os.Exit(0)
	// }

	var ProcessBlock []p.Process

	// iterates through the book.nav.toc array and splits the path on "Fmt425-" and "#" to enumerate the seconds (if applicable) and match to the file paths above
	// for _, toc := range book.Nav.Toc {
	for i := 0; i < len(book.Nav.Toc); i++ {

		// Makes an empty Process object
		var process p.Process

		toc := book.Nav.Toc[i]
		part, seconds := p.GetFileNameAndSeconds(toc.Path)

		// Lookup the mp3File directly using a map
		if mp3File, ok := mp3FileMap[part]; ok {
			process.Source = mp3File
		} else {
			fmt.Println("Part not found:", part)
			os.Exit(1)
		}

		// Gets the next file in the mp3Files array and checks if it matches the path of the toc
		if i < len(book.Nav.Toc)-1 {
			toc2 := book.Nav.Toc[i+1]
			part2, seconds2 := p.GetFileNameAndSeconds(toc2.Path)

			if part != part2 {

				// parts don't match, but second part is 0, so we just go to the end of the first file
				if seconds2 == 0 {
					seconds2, err = p.GetFileDuration(process.Source)
					if err != nil {
						fmt.Println("Error getting duration:", err)
						os.Exit(1)
					}
				} else { // parts don't match, and second part is not 0, so we need to go to the end of the first file and partially into the second file

					cmd, err := p.GetComplexSplit(process.Source, mp3FileMap[part2], seconds, seconds2)
					if err != nil {
						fmt.Println("Error getting command:", err)
						os.Exit(1)
					}

					process.Override = cmd

				}

			}

			// if seconds2 == 0 {
			// 	seconds2, err = p.GetFileDuration(process.Source)
			// 	if err != nil {
			// 		fmt.Println("Error getting duration:", err)
			// 		os.Exit(1)
			// 	}
			// } else if part2 != part {

			// 	file1 := process.Source
			// 	file2 := mp3FileMap[part2]
			// 	cmd, err := p.GetComplexSplit(file1, file2, seconds, seconds2)
			// 	if err != nil {
			// 		fmt.Println("Error getting command:", err)
			// 		os.Exit(1)
			// 	}

			// 	// mp3File := mp3FileMap[part2]
			// 	// bitrate, err := p.GetBitRate(process.Source)
			// 	// if err != nil {
			// 	// 	fmt.Println("Error getting bitrate:", err)
			// 	// 	os.Exit(1)
			// 	// }

			// 	// bitrate2, err := p.GetBitRate(mp3File)
			// 	// if err != nil {
			// 	// 	fmt.Println("Error getting bitrate:", err)
			// 	// 	os.Exit(1)
			// 	// }

			// 	// var bitrateLarger int
			// 	// if bitrate > bitrate2 {
			// 	// 	bitrateLarger = bitrate
			// 	// } else {
			// 	// 	bitrateLarger = bitrate2
			// 	// }

			// 	// cmd := exec.Command("ffmpeg",
			// 	// 	"-i", process.Source,
			// 	// 	"-i", mp3File,
			// 	// 	"-filter_complex", fmt.Sprintf("[0:a]atrim=start=%.2f[a1];[1:a]atrim=start=0:end=%.2f[a2];[a1][a2]concat=n=2:v=0:a=1[out]", seconds, seconds2),
			// 	// 	"-map", "[out]",
			// 	// 	"-b:a", fmt.Sprintf("%d", bitrateLarger),
			// 	// )

			// 	process.Override = cmd
			// }

			process.End = seconds2

		} else if i == len(book.Nav.Toc)-1 {

			process.End, err = p.GetFileDuration(process.Source)
			if err != nil {
				fmt.Println("Error getting duration:", err)
				os.Exit(1)
			}

		}

		// Adds the title, start time, and output path
		process.Title = toc.Title
		process.Start = seconds

		outputFileNormal := p.NormalizeFileName(toc.Title)

		iteration := fmt.Sprintf("%02d", i)

		// Sets the output path
		process.Output = path.Join(outputPath, "["+iteration+"]. "+outputFileNormal+".mp3")

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

		if process.Override != nil {

			fmt.Printf("Processing Chapter: %s\n", process.Title)

			newCmd := process.Override

			// processString := fmt.Sprintf("'%s'", process.Output)
			newCmd.Args = append(newCmd.Args, process.Output)

			output, err := newCmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error running command:", err)
				fmt.Println("Command Output: ", string(output))
				return
			}

		} else {

			durationSeconds := process.End - process.Start
			duration := p.GetDurationFormatted(durationSeconds)
			fmt.Printf("Processing Chapter: %s (%s)\n", process.Title, duration)

			cmd := exec.Command("ffmpeg",
				"-i", process.Source,
				"-ss", fmt.Sprintf("%f", process.Start),
				"-to", fmt.Sprintf("%f", process.End),
				"-acodec", "copy",
				process.Output,
				"-hide_banner", "-loglevel", "error")

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
