package main

import (
	p "Z0y6h0kS9X/libby-chapterizer/pkg"
	prov "Z0y6h0kS9X/libby-chapterizer/provider"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
var audibleChapters bool

func init() {
	rootCmd.Flags().StringVarP(&jsonPath, "json", "j", "", "The path to the openbook.json file")
	rootCmd.Flags().StringVarP(&outPath, "out", "o", "", "The path to the directory you want to output the files to")
	rootCmd.Flags().BoolVarP(&test, "test", "t", false, "Test mode")
	rootCmd.Flags().BoolVarP(&audibleChapters, "use-audible-chapters", "c", false, "Specifies to override default breaks and use audible markers instead")
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

	// Reads the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshals JSON
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

	fmt.Println("Looking up Book ASIN...")

	// Gets the ASIN
	asin, err := prov.GetBook(book.Title.Main, authorString, narratorString)
	if err != nil {
		fmt.Println("Error getting book:", err)
		return
	}

	if asin == "" {
		fmt.Println("No ASIN found")
		return
	} else {
		fmt.Println("ASIN:", asin)
	}

	// Gets the book details
	details, err := prov.GetBookDetails(asin)
	if err != nil {
		fmt.Println("Error getting book details:", err)
		return
	}

	// Checks if the user specified an output path
	var outDir string
	if outPath != "" {
		outDir = strings.Replace(outPath, "\\", "/", -1)
	} else {
		outDir = fileDir
	}

	// Gets the output dir path
	outputPath, err := p.GetOutputDirPath(details, asin, outDir)
	if err != nil {
		fmt.Println("Error getting output dir path:", err)
		return
	}

	// Checks if the folder exists and creates it if it does not
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Gets a list of all the .mp3 files in the fileDir
	files, err := os.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	mp3FileMap := make(map[string]string)
	totalDuration := 0.000
	var mp3Files []string

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

			// If use-audible-chapters flag is set, add the full path to the array for later processing
			if audibleChapters {
				mp3Files = append(mp3Files, fullPath)
			}
		}
	}

	duration := p.FormatDuration(totalDuration)

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

	if test {
		fmt.Println("Test mode enabled. Exiting...")
		os.Exit(0)

	}

	var ProcessBlock []p.Process

	var tempFile string

	if audibleChapters {

		// Set the output path for the combined MP3 file
		fileName := fmt.Sprintf("%s (%s)", p.NormalizeName(book.Title.Main), asin)
		tempFile = path.Join(outputPath, fileName+".mp3")

		// Call the MakeCombinedMP3 function to create the combined MP3 file
		err = p.MakeCombinedMP3(mp3Files, tempFile)
		if err != nil {
			fmt.Println("Error making combined MP3:", err)
			return
		}

		chapters, err := prov.GetChapters(asin)
		if err != nil {
			fmt.Println("Error getting chapters:", err)
			return
		}

		for i, chapter := range chapters.Chapters {

			title := chapter.Title
			start, _ := strconv.ParseFloat(strconv.Itoa(chapter.StartOffsetMs/1000), 64)
			dur, _ := strconv.ParseFloat(strconv.Itoa(chapter.LengthMs/1000), 64)
			end := start + dur

			process := p.Process{}

			cmd, err := p.GetSimpleSplit(tempFile, start, end)
			if err != nil {
				fmt.Println("Error getting simple split:", err)
				return
			}

			// Normalizes the title
			outputFileNormal := p.NormalizeName(title)
			iteration := fmt.Sprintf("%02d", i)

			process.Source = tempFile
			process.Start = start
			process.End = end
			process.Duration = p.FormatDuration(dur)
			process.Title = title
			process.Command = cmd
			process.Output = path.Join(outputPath, "["+iteration+"]. "+outputFileNormal+".mp3")

			ProcessBlock = append(ProcessBlock, process)

		}

	} else {

		for i, toc := range book.Nav.Toc {

			process := p.Process{}
			part, seconds := p.GetFileNameAndSeconds(toc.Path)

			switch mp3File := mp3FileMap[part]; {
			case mp3File != "":
				process.Source = mp3File
			default:
				log.Fatalf("Part not found: %s", part)
			}

			var cmd *exec.Cmd
			var dur string

			if i < len(book.Nav.Toc)-1 {
				toc2 := book.Nav.Toc[i+1]
				part2, seconds2 := p.GetFileNameAndSeconds(toc2.Path)

				if part != part2 {
					if seconds2 == 0 {
						seconds2, err = p.GetFileDuration(process.Source)
						if err != nil {
							log.Fatalf("Error getting duration: %v", err)
						}
						dur, err = p.GetSimpleDuration(seconds, seconds2)
						if err != nil {
							log.Fatalf("Error getting duration: %v", err)
						}
						cmd, err = p.GetSimpleSplit(process.Source, seconds, seconds2)
						if err != nil {
							log.Fatalf("Error getting command: %v", err)
						}
					} else {
						dur, err = p.GetComplexDuration(process.Source, seconds, seconds2)
						if err != nil {
							log.Fatalf("Error getting duration: %v", err)
						}
						cmd, err = p.GetComplexSplit(process.Source, mp3FileMap[part2], seconds, seconds2)
						if err != nil {
							log.Fatalf("Error getting command: %v", err)
						}
					}
				} else {
					cmd, err = p.GetSimpleSplit(process.Source, seconds, seconds2)
					if err != nil {
						log.Fatalf("Error getting command: %v", err)
					}
					dur, err = p.GetSimpleDuration(seconds, seconds2)
					if err != nil {
						log.Fatalf("Error getting duration: %v", err)
					}
				}

				process.End = seconds2
			} else if i == len(book.Nav.Toc)-1 {
				process.End, err = p.GetFileDuration(process.Source)
				if err != nil {
					fmt.Println("Error getting file duration:", err)
					os.Exit(1)
				}

				cmd, err = p.GetSimpleSplit(process.Source, seconds, process.End)
				if err != nil {
					fmt.Println("Error getting command:", err)
					os.Exit(1)
				}

				dur, err = p.GetSimpleDuration(seconds, process.End)
				if err != nil {
					fmt.Println("Error getting duration:", err)
					os.Exit(1)
				}
			}

			// Normalizes the title
			outputFileNormal := p.NormalizeName(toc.Title)
			iteration := fmt.Sprintf("%02d", i)

			// Adds the Process object properties
			process.Title = toc.Title
			process.Start = seconds
			process.Duration = dur
			process.Command = cmd
			process.Output = path.Join(outputPath, "["+iteration+"]. "+outputFileNormal+".mp3")

			// Adds the process to the ProcessBlock
			ProcessBlock = append(ProcessBlock, process)
		}

	}

	// Runs the commands to generate the output files
	var m3u []string
	m3u = append(m3u, "#EXTM3U")
	m3u = append(m3u, fmt.Sprintf("#PLAYLIST: %s", book.Title.Main))

	for _, process := range ProcessBlock {

		_, file := path.Split(process.Output)

		fmt.Printf("Processing Item: %s (%s)\n", process.Title, process.Duration)
		newCmd := process.Command
		newCmd.Args = append(newCmd.Args, process.Output)

		output, err := newCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error running command:", err)
			fmt.Println("Command Output: ", string(output))
			return
		}

		m3u = append(m3u, fmt.Sprintf("#EXTINF:,%s\n%s", process.Title, file))

	}

	content := strings.Join(m3u, "\n")
	err = os.WriteFile(path.Join(outputPath, p.NormalizeName(book.Title.Main)+".m3u"), []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	// removes the tempFile if it exists
	if _, err := os.Stat(tempFile); err == nil {
		os.Remove(tempFile)
	}

}
