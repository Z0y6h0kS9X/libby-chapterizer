package main

import (
	p "Z0y6h0kS9X/libby-chapterizer/pkg"
	prov "Z0y6h0kS9X/libby-chapterizer/provider"
	"fmt"
	"os"
	"path"
	"path/filepath"

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
var single bool
var format string

func init() {
	rootCmd.Flags().StringVarP(&jsonPath, "json", "j", "", "The path to the openbook.json file")
	rootCmd.Flags().StringVarP(&outPath, "out", "o", "", "The path to the directory you want to output the files to")
	rootCmd.Flags().BoolVarP(&test, "test", "t", false, "Test mode")
	rootCmd.Flags().BoolVarP(&audibleChapters, "use-audible-chapters", "c", false, "Specifies to override default breaks and use audible markers instead")
	rootCmd.Flags().BoolVarP(&single, "single", "s", false, "Indicates if you want the output as a single file, or sepearate files for each chapter")
	rootCmd.Flags().StringVarP(&format, "format", "f", "mp3", "What format you want the output in (mp3|m4b)")
}

func main() {

	// Parses the flags
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if test {
		meta, err := p.GetMetadataFromASIN("B002UZKI96")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(meta.ToString())

		// fmt.Println("Test mode enabled. Exiting...")

		// // Creates a generic playlist object
		// playlist := p.Playlist{
		// 	Title:     "TestTitle",
		// 	Author:    "TestAuthor",
		// 	BookTitle: "TestBook",
		// 	Tracks: []p.Track{
		// 		{
		// 			Title:    "Chapter 1",
		// 			Length:   27,
		// 			LengthMS: 27000,
		// 			File:     "[01]. Chapter 1.mp3",
		// 		},
		// 		{
		// 			Title:    "Chapter 2",
		// 			Length:   54,
		// 			LengthMS: 54380,
		// 			File:     "[02]. Chapter 2.mp3",
		// 		},
		// 	},
		// }

		// playlist.WriteM3UFile("F:/Music/Audiobook/Output/Gary Paulsen/Hatchet/[04.0]. Brian's Return (B00771TZ92)/generated.m3u")
		// playlist.WriteFFMPEGMetadataFile("F:/Music/Audiobook/Output/Gary Paulsen/Hatchet/[04.0]. Brian's Return (B00771TZ92)/metadata-generated.txt")

		// Parses the m3u file
		// playlist, err := p.ReadExtM3U("F:/Music/Audiobook/Output/Gary Paulsen/Hatchet/[04.0]. Brian's Return (B00771TZ92)/Brian's Return-mod.m3u")
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// fmt.Println("----------")
		// fmt.Println("Playlist:")
		// fmt.Println("----------")
		// fmt.Println("Playlist Title:", playlist.Title)
		// for _, track := range playlist.Tracks {
		// 	fmt.Printf("Track Title: %s | Track Length: %d | Track Length MS: %d\nFile: %s\n", track.Title, track.Length, track.LengthMS, track.File)
		// }

		// Ends program execution
		os.Exit(0)

	}

	// Checks to see if the jsonPath was specified
	if jsonPath == "" {
		fmt.Println("Error: path to openbook.json was not specified")
		os.Exit(1)
	}

	// Checks to see if the jsonPath is valid
	_, err := os.Stat(jsonPath)
	if os.IsNotExist(err) {
		fmt.Println("Error: path to openbook.json is not valid")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Error!:", err)
		os.Exit(1)
	}

	// Gets the directory path from the json path, converst to *nix path (if windows)
	jsonPath = filepath.ToSlash(jsonPath)
	jsonDir := path.Dir(jsonPath)

	// Checks to see if the outPath was specified
	if outPath == "" {
		outPath = jsonDir
	}

	// Checks to see if the outPath is valid
	_, err = os.Stat(outPath)
	if os.IsNotExist(err) {
		// Create path if it doesn't exist
		err = os.MkdirAll(outPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating output path:", err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Gets the directory path from the json path, converst to *nix path (if windows)
	outPath = filepath.ToSlash(outPath)

	// Checks the output format specified is valid
	if format != "mp3" && format != "m4b" {
		fmt.Println("Error: output format must be 'mp3' or 'm4b'")
		os.Exit(1)
	}

	// Converts the JSON file to an Openbook
	book, err := p.JSONFileToOpenBook(jsonPath)
	if err != nil {
		fmt.Println("Error: Unable to convert JSON file to Openbook!\n", err)
		os.Exit(1)
	}

	// Gets the primary author and narrator
	author := p.GetPrimaryAuthor(book)
	narrator := p.GetPrimaryNarrator(book)

	fmt.Println("Looking up Book ASIN...")

	// Gets the ASIN
	asin, err := prov.GetBook(book.Title.Main, author, narrator, book.CalculateRuntime())
	if err != nil {
		fmt.Println("Error getting book:", err)
		os.Exit(1)
	}

	var details p.BookDetails
	// If there is no ASIN, create details manually using openbook
	if asin == "" {
		fmt.Println("No ASIN found")

		details, err = p.GetBookDetailsNoASIN(book)
		if err != nil {
			fmt.Println("Error getting book details:", err)
			return
		}

	} else {
		fmt.Println("ASIN:", asin)
		// Gets the book details
		details, err = prov.GetBookDetailsASIN(asin)
		if err != nil {
			fmt.Println("Error getting book details:", err)
			return
		}
	}

	// Gets the output dir path
	outputPath, err := p.GetOutputDirPath(details, asin, outPath)
	if err != nil {
		fmt.Println("Error getting output dir path:", err)
		return
	}

	// Prints the book details
	fmt.Println("Author:", author)
	fmt.Println("Narrator:", narrator)
	fmt.Println("Directory:", jsonDir)
	fmt.Println("Output Directory:", outPath)
	fmt.Println("Output Path:", outputPath)
	fmt.Println("Format:", format)
	if asin != "" {
		fmt.Println("ASIN:", asin)
	} else {
		fmt.Println("ASIN: Book does not have an ASIN")
	}
	if single {
		fmt.Println("Output Type: Single File")
	} else {
		fmt.Println("Output Type: Multiple Files")
	}
	if audibleChapters {
		fmt.Println("Audible Chapters: Enabled")
	} else {
		fmt.Println("Audible Chapters: Disabled")
	}

	// ------------ Starts Destructive Code ------------

	// Gets a list of all the .mp3 files in the jsonDir
	files1, err := p.GetAllMp3Files(jsonDir)
	if err != nil {
		fmt.Println("Error getting list of .mp3 files:", err)
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

	var outputName string
	// var makeTemp bool
	var tempFile string

	// if audibleChapters && format == "mp3" && single {
	// 	fmt.Println("Single, MP3, and Use Audible Chapters enabled - will ignore chapters as mp3 does not support chapters")
	// 	makeTemp = true
	// } else if audibleChapters && format == "m4b" && single {
	// 	fmt.Println("Single, M4B, and Use Audible Chapters enabled - will output a single chapterized m4b using audible chapters")
	// 	// Generate a ffmpeg metadata file using
	// }

	// If we are using audible chapters, we will create a single temp file - regardless of format
	if audibleChapters {

		fmt.Println("Creating single temp file for audible chapters")
		tempFile = path.Join(outputPath, "temp.mp3")
		p.MakeCombinedMP3(files1, details, tempFile)
	}

	// Does the simplest check, if single is true and format is mp3, then create a single file with limited metdata
	if single && format == "mp3" {
		fmt.Println("Creating single MP3 file with limited metadata, as it does not support chapters")
		outputName = path.Join(outputPath, fmt.Sprintf("%s (%s).mp3", p.NormalizeName(details.Title), details.Asin))

		// If the tempfile exists, rename it - else create it
		if _, err := os.Stat(tempFile); err == nil {
			// Temp file exists, rename it to output name
			err := os.Rename(tempFile, outputName)
			if err != nil {
				fmt.Println("Error renaming temp file:", err)
				return
			}
		} else {
			// Temp file does not exist, create it
			p.MakeCombinedMP3(files1, details, outputName)
		}

		fmt.Println("File created:", outputName)

		// Successfully exits the program
		os.Exit(0)
	}

	if single && format == "m4b" {

	}

	// In any other situation, we will need to get the chapters
	// If not single and format == mp3, we will create multiple files
	// if !single && format == "mp3" {

	// }

	// Abort before doing anything
	os.Exit(0)
}

// // Gets a list of all the .mp3 files in the jsonDir
// files, err := os.ReadDir(jsonDir)
// if err != nil {
// 	fmt.Println("Error reading directory:", err)
// 	return
// }

// mp3FileMap := make(map[string]string)
// totalDuration := 0.000
// var mp3Files []string

// for _, file := range files {
// 	if path.Ext(file.Name()) == ".mp3" {

// 		fullPath := path.Join(jsonDir, file.Name())
// 		// mp3Files = append(mp3Files, fullPath)
// 		part := p.GetPartFromMp3File(fullPath)
// 		mp3FileMap[part] = fullPath

// 		duration, err := p.GetFileDuration(fullPath)
// 		if err != nil {
// 			fmt.Println("Error getting file duration:", err)
// 			return
// 		}
// 		totalDuration += duration

// 		// If use-audible-chapters flag is set, add the full path to the array for later processing
// 		if audibleChapters {
// 			mp3Files = append(mp3Files, fullPath)
// 		}
// 	}
// }

// duration := p.FormatDuration(totalDuration)

// fmt.Println("============================")
// fmt.Println("Audiobook Information")
// fmt.Println("----------------------------")
// fmt.Println("Title:", book.Title.Main)
// fmt.Println("ASIN:", asin)
// if details.SeriesPrimary.Name != "" {
// 	fmt.Println("Series:", details.SeriesPrimary.Name)
// 	fmt.Println("Book Number:", details.SeriesPrimary.Position)
// }
// fmt.Println("Author:", author)
// fmt.Println("Narrator:", narrator)
// fmt.Println("Duration:", duration)
// fmt.Println("============================")

// var ProcessBlock []p.Process

// var tempFile string

// // If either use-audible-chapters or single flag is set, create the combined MP3 file
// if audibleChapters || single {

// 	// Set the output path for the combined MP3 file
// 	fileName := fmt.Sprintf("%s (%s)", p.NormalizeName(book.Title.Main), asin)
// 	tempFile = path.Join(outputPath, fileName+"_temp.mp3")

// 	// Call the MakeCombinedMP3 function to create the combined MP3 file
// 	err = p.MakeCombinedMP3(mp3Files, tempFile)
// 	if err != nil {
// 		fmt.Println("Error making combined MP3:", err)
// 		return
// 	}

// }

// 	if audibleChapters {

// 		// Set the output path for the combined MP3 file
// 		fileName := fmt.Sprintf("%s (%s)", p.NormalizeName(book.Title.Main), asin)
// 		tempFile = path.Join(outputPath, fileName+"_temp.mp3")

// 		// Call the MakeCombinedMP3 function to create the combined MP3 file
// 		err = p.MakeCombinedMP3(mp3Files, details, tempFile)
// 		if err != nil {
// 			fmt.Println("Error making combined MP3:", err)
// 			return
// 		}

// 		chapters, err := prov.GetChapters(asin)
// 		if err != nil {
// 			fmt.Println("Error getting chapters:", err)
// 			return
// 		}

// 		for i, chapter := range chapters.Chapters {

// 			title := chapter.Title
// 			start, _ := strconv.ParseFloat(strconv.Itoa(chapter.StartOffsetMs/1000), 64)
// 			dur, _ := strconv.ParseFloat(strconv.Itoa(chapter.LengthMs/1000), 64)
// 			end := start + dur

// 			process := p.Process{}

// 			cmd, err := p.GetSimpleSplit(tempFile, start, end)
// 			if err != nil {
// 				fmt.Println("Error getting simple split:", err)
// 				return
// 			}

// 			// Normalizes the title
// 			outputFileNormal := p.NormalizeName(title)
// 			iteration := fmt.Sprintf("%02d", i)

// 			process.Source = tempFile
// 			process.Start = start
// 			process.End = end
// 			process.DurationStr = p.FormatDuration(dur)
// 			process.Title = title
// 			process.Command = cmd
// 			process.Output = path.Join(outputPath, "["+iteration+"]. "+outputFileNormal+".mp3")

// 			ProcessBlock = append(ProcessBlock, process)

// 		}

// 	} else {

// 		for i, toc := range book.Nav.Toc {

// 			process := p.Process{}
// 			part, seconds := p.GetFileNameAndSeconds(toc.Path)

// 			switch mp3File := mp3FileMap[part]; {
// 			case mp3File != "":
// 				process.Source = mp3File
// 			default:
// 				log.Fatalf("Part not found: %s", part)
// 			}

// 			var cmd *exec.Cmd
// 			var dur string

// 			if i < len(book.Nav.Toc)-1 {
// 				toc2 := book.Nav.Toc[i+1]
// 				part2, seconds2 := p.GetFileNameAndSeconds(toc2.Path)

// 				if part != part2 {
// 					if seconds2 == 0 {
// 						seconds2, err = p.GetFileDuration(process.Source)
// 						if err != nil {
// 							log.Fatalf("Error getting duration: %v", err)
// 						}
// 						dur, err = p.GetSimpleDuration(seconds, seconds2)
// 						if err != nil {
// 							log.Fatalf("Error getting duration: %v", err)
// 						}
// 						cmd, err = p.GetSimpleSplit(process.Source, seconds, seconds2)
// 						if err != nil {
// 							log.Fatalf("Error getting command: %v", err)
// 						}
// 					} else {
// 						dur, err = p.GetComplexDuration(process.Source, seconds, seconds2)
// 						if err != nil {
// 							log.Fatalf("Error getting duration: %v", err)
// 						}
// 						cmd, err = p.GetComplexSplit(process.Source, mp3FileMap[part2], seconds, seconds2)
// 						if err != nil {
// 							log.Fatalf("Error getting command: %v", err)
// 						}
// 					}
// 				} else {
// 					cmd, err = p.GetSimpleSplit(process.Source, seconds, seconds2)
// 					if err != nil {
// 						log.Fatalf("Error getting command: %v", err)
// 					}
// 					dur, err = p.GetSimpleDuration(seconds, seconds2)
// 					if err != nil {
// 						log.Fatalf("Error getting duration: %v", err)
// 					}
// 				}

// 				process.End = seconds2
// 			} else if i == len(book.Nav.Toc)-1 {
// 				process.End, err = p.GetFileDuration(process.Source)
// 				if err != nil {
// 					fmt.Println("Error getting file duration:", err)
// 					os.Exit(1)
// 				}

// 				cmd, err = p.GetSimpleSplit(process.Source, seconds, process.End)
// 				if err != nil {
// 					fmt.Println("Error getting command:", err)
// 					os.Exit(1)
// 				}

// 				dur, err = p.GetSimpleDuration(seconds, process.End)
// 				if err != nil {
// 					fmt.Println("Error getting duration:", err)
// 					os.Exit(1)
// 				}
// 			}

// 			// Normalizes the title
// 			outputFileNormal := p.NormalizeName(toc.Title)
// 			iteration := fmt.Sprintf("%02d", i)

// 			// Adds the Process object properties
// 			process.Title = toc.Title
// 			process.Start = seconds
// 			process.DurationStr = dur
// 			process.Command = cmd
// 			process.Output = path.Join(outputPath, "["+iteration+"]. "+outputFileNormal+".mp3")

// 			// Adds the process to the ProcessBlock
// 			ProcessBlock = append(ProcessBlock, process)
// 		}

// 	}

// 	// Runs the commands to generate the output files
// 	// Starts building the m3u file
// 	// var m3u []string
// 	// m3u = append(m3u, "#EXTM3U")
// 	// m3u = append(m3u, fmt.Sprintf("#PLAYLIST: %s", book.Title.Main))

// 	// Starts building the metadata file
// 	// var metadata []string
// 	// metadata = append(metadata, ";FFMETADATA1")
// 	// metadata = append(metadata, fmt.Sprintf("title=%s", book.Title.Main))
// 	// metadata = append(metadata, fmt.Sprintf("artist=%s", author))
// 	// metadata = append(metadata, fmt.Sprintf("\n"))

// 	// durMS := 0
// 	playlist := p.Playlist{}

// 	// Sets the top level playlist properties
// 	playlist.Title = book.Title.Main
// 	playlist.Author = author
// 	playlist.BookTitle = book.Title.Main

// 	for _, process := range ProcessBlock {

// 		_, file := path.Split(process.Output)

// 		fmt.Printf("Processing Item: %s (%s)\n", process.Title, process.DurationStr)
// 		newCmd := process.Command
// 		newCmd.Args = append(newCmd.Args, process.Output)

// 		output, err := newCmd.CombinedOutput()
// 		if err != nil {
// 			fmt.Println("Error running command:", err)
// 			fmt.Println("Command Output: ", string(output))
// 			return
// 		}

// 		// Gets duration of the output file
// 		durationMS, err := p.GetFileDurationMS(process.Output)
// 		if err != nil {
// 			fmt.Println("Error getting duration:", err)
// 			return
// 		}

// 		length := int(durationMS / 1000)

// 		playlist.Tracks = append(playlist.Tracks, p.Track{Title: process.Title, File: file, Length: length, LengthMS: durationMS})

// 		// Adds entry to the m3u file
// 		// m3u = append(m3u, fmt.Sprintf("#EXTINF:,%s\n%s", process.Title, file))

// 		// Generates Process Block metadata
// 		// dur, err := p.GetFileDurationMS(process.Output)
// 		// if err != nil {
// 		// 	fmt.Println("Error getting duration:", err)
// 		// 	return
// 		// }
// 		// metadata = append(metadata, p.GenerateChapterBlock(file, process.Title, dur, durMS))
// 		// durMS += dur

// 	}

// 	// Writes the playlist out in M3u and FFMPEG formats
// 	playlist.WriteM3UFile(path.Join(outputPath, p.NormalizeName(book.Title.Main)+".m3u"))
// 	playlist.WriteFFMPEGMetadataFile(path.Join(outputPath, "metadata.txt"))

// 	// // Gets the list of files in the output directory
// 	// mp3s, err := p.GetAllMp3Files(outputPath)
// 	// if err != nil {
// 	// 	fmt.Println("Error getting mp3 files:", err)
// 	// 	return
// 	// }

// 	// content := strings.Join(m3u, "\n")
// 	// err = os.WriteFile(path.Join(outputPath, p.NormalizeName(book.Title.Main)+".m3u"), []byte(content), 0644)
// 	// if err != nil {
// 	// 	fmt.Println("Error writing M3U file:", err)
// 	// 	return
// 	// }

// 	// content = strings.Join(metadata, "\n")
// 	// err = os.WriteFile(path.Join(outputPath, "metadata.txt"), []byte(content), 0644)
// 	// if err != nil {
// 	// 	fmt.Println("Error writing metadata file:", err)
// 	// 	return
// 	// }

// 	if single {

// 		fmt.Println("Generating combined chapterized file...")

// 		fileName := fmt.Sprintf("%s (%s)", p.NormalizeName(book.Title.Main), asin)
// 		outCombined := path.Join(outputPath, fileName+"."+format)

// 		var input string

// 		// Converts the temp mp3 into an m4b file with chapters
// 		if _, err := os.Stat(tempFile); err == nil {

// 			// sets the input to the temp file
// 			input = tempFile

// 		} else {

// 			// Converts with concat commands
// 			mp3s, err := p.GetAllMp3Files(outputPath)
// 			if err != nil {
// 				fmt.Println("Error getting mp3 files:", err)
// 				return
// 			}

// 			// Gets the list of files in the output directory
// 			var concat []string
// 			concat = append(concat, mp3s...)

// 			// Sets the input to the concat command
// 			input = "concat:" + strings.Join(concat, "|")

// 		}

// 		if format == "m4b" {
// 			cmd := exec.Command("ffmpeg", "-i", input, "-i", path.Join(outputPath, "metadata.txt"), "-acodec", "aac", "-strict", "experimental", "-ac", "1", "-vn", outCombined)

// 			output, err := cmd.CombinedOutput()
// 			if err != nil {
// 				fmt.Println("Error running command:", err)
// 				fmt.Println("Command Output: ", string(output))
// 				return
// 			}
// 		} else if format == "mp3" {
// 			fmt.Println("MP3's do not support chapters, only Top Level Metadata will be added")
// 			err = p.MakeCombinedMP3(mp3Files, details, tempFile)
// 			if err != nil {
// 				fmt.Println("Error making combined MP3:", err)
// 				return
// 			}
// 		} else {
// 			fmt.Println("Invalid format:", format)
// 			return
// 		}

// 		// Preforms cleanup, removing everything except the outCombined file
// 		files, err := os.ReadDir(outputPath)
// 		if err != nil {
// 			fmt.Println("Error reading directory:", err)
// 			return
// 		}

// 		for _, file := range files {

// 			fullPath := path.Join(outputPath, file.Name())
// 			if fullPath != outCombined {
// 				err := os.Remove(fullPath)
// 				if err != nil {
// 					fmt.Println("Error removing file:", err)
// 				}
// 			}

// 		}

// 	} else {

// 		// removes the tempFile if it exists
// 		if _, err := os.Stat(tempFile); err == nil {
// 			os.Remove(tempFile)
// 		}

// 	}

// }

// ffmpeg  -i "concat:chapter1.mp3|chapter2.mp3|chapter3.mp3|chapter4.mp3" -i .\metadata.txt -acodec aac -strict experimental -ac 1 -vn output_with_chapters.m4b
