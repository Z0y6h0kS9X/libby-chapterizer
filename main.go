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

	// Gets the ASIN
	asin, err := prov.GetBook(book.Title.Main, author, narrator, book.CalculateRuntime())
	if err != nil {
		fmt.Println("Error getting book:", err)
		os.Exit(1)
	}

	// var details p.BookDetails
	var metadata p.Metadata
	// If there is no ASIN, create details manually using openbook
	if asin == "" {

		metadata, err = p.GetMetadataLocal(book)
		if err != nil {
			fmt.Println("Error getting metadata (No ASIN):", err)
			os.Exit(1)
		}

	} else {
		metadata, err = p.GetMetadataFromASIN(asin)
		if err != nil {
			fmt.Println("Error getting metadata (ASIN):", err)
			os.Exit(1)
		}
	}

	outputPath, err := p.GetOutputDirPath(metadata, asin, outPath)
	if err != nil {
		fmt.Println("Error getting output dir path:", err)
		return
	}

	// Prints the book details
	fmt.Println("=================== Book Details ====================")
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
	fmt.Println("=====================================================")

	// ------------ Starts Destructive Code ------------

	// Gets a list of all the .mp3 files in the jsonDir
	files, err := p.GetAllMp3Files(jsonDir)
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

	// Check if the user wants to use audible chapters or not
	var chapters []p.Chapter
	if audibleChapters {
		chapters, err = prov.GetAudibleChapters(metadata.ASIN)
		if err != nil {
			fmt.Println("Error getting audible chapters:", err)
			os.Exit(1)
		}

		if chapters == nil {
			fmt.Println("No audible chapters found, using local chapters")
			chapters, err = p.GetChaptersLocal(book, files)
			if err != nil {
				fmt.Println("Error getting local chapters:", err)
				os.Exit(1)
			}
		}
	} else {
		chapters, err = p.GetChaptersLocal(book, files)
		if err != nil {
			fmt.Println("Error getting local chapters:", err)
			os.Exit(1)
		}
	}

	metadata.Chapters = chapters

	// Check if the output will be a single file or not
	if single {

		var outputFile string
		if asin == "" {
			outputFile = path.Join(outputPath, fmt.Sprintf("%s.%s", metadata.Title, format))
		} else {
			outputFile = path.Join(outputPath, fmt.Sprintf("%s (%s).%s", metadata.Title, asin, format))
		}

		if format == "mp3" {

			// Output single mp3, with limited metadata
			p.MakeCombinedMP3(files, metadata, outputFile)

		} else {
			fmt.Println("Making single m4b file")

			ffmetadata := metadata.ToFFMPEGMetadata()

			// Create the file
			file, err := os.Create(outputPath + "/ffmetadata.txt")
			if err != nil {
				fmt.Println("Error creating ffmetadata file:", err)
				os.Exit(1)
			}
			defer file.Close()

			// Writes the contents of ffmetadata out to the file
			_, err = file.WriteString(ffmetadata)
			if err != nil {
				fmt.Println("Error writing ffmetadata file:", err)
				os.Exit(1)
			}

			// Output single m4b with metadata
			err = p.MakeCombinedM4B(files, outputPath+"/ffmetadata.txt", outputFile)
			if err != nil {
				fmt.Println("Error making single m4b file:", err)
				os.Exit(1)
			}

		}

	} else {

		//
		if format == "mp3" {
			// Output split mp3s
			err = p.MakeSplitMP3Files(files, chapters, metadata, outputPath)
			if err != nil {
				fmt.Println("Error making split mp3 files:\n", err)
				os.Exit(1)
			}
		} else {
			// Output split m4bs
			err = p.MakeSplitM4BFiles(files, chapters, metadata, outputPath)
			if err != nil {
				fmt.Println("Error making split m4b files:\n", err)
				os.Exit(1)
			}
		}

	}
}
