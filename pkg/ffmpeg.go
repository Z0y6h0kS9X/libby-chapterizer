// This file is responsible for executing ffmp... commands.

package pkg

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// GetFileDurationMS calculates the duration of a file in milliseconds.
// It takes the filepath as input and returns the duration in milliseconds and any error encountered.
func GetFileDurationMS(filepath string) (int, error) {

	// Create a new exec.Command with the ffprobe command and arguments
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filepath)

	// Run the command and capture the stdout
	stdout, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error running ffprobe command: %w", err)
	}

	// Trim any leading or trailing whitespace from the stdout and convert it to a float
	durationInSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(stdout)), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}

	// Converts the duration to milliseconds
	durationInMilliseconds := int(durationInSeconds * 1000)

	// Return the duration
	return durationInMilliseconds, nil
}

// MakeCombinedMP3 concatenates multiple MP3 files into a single output file.
// It takes a slice of file paths and the path of the output file as input.
// It returns an error if the operation fails.
func MakeCombinedMP3(files []string, meta Metadata, outputFile string) error {

	fmt.Println("Making Combined MP3...")

	// Create a slice to store the command line arguments
	var args []string

	// Append the input file paths to the arguments slice using the "concat" format
	args = append(args, "-i", "concat:"+strings.Join(files, "|"))

	// Adds simple metadata to the output file
	args = append(args, "-metadata", "title="+meta.Title)
	args = append(args, "-metadata", "artist="+meta.Author)
	args = append(args, "-metadata", "album="+meta.Title)
	args = append(args, "-metadata", "ASIN="+meta.ASIN)

	// Set the audio codec to "copy" to preserve the original audio codecs
	args = append(args, "-acodec", "copy", outputFile)

	// Create a new command using the "ffmpeg" executable and the arguments
	cmd := exec.Command("ffmpeg", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
	}

	// Return nil if the operation is successful
	return nil
}

// MakeCombinedM4B combines multiple audio files into a single M4B file
// and adds metadata to it.
//
// Parameters:
//   - files: a slice of input file paths
//   - metadataFile: the path to the metadata file
//   - outputFile: the path to the output file
//
// Returns:
//   - an error if the operation fails, or nil if successful
func MakeCombinedM4B(files []string, metadataFile, outputFile string) error {
	// Print a message to indicate that the function is starting
	fmt.Println("Making Combined M4B...")

	// Create a slice to store the command line arguments
	var args []string

	// Append the input file paths to the arguments slice using the "concat" format
	args = append(args, "-i", "concat:"+strings.Join(files, "|"))

	// Adds the metadata file to the output file
	args = append(args, "-i", metadataFile)

	// Sets the metadata map
	args = append(args, "-map_metadata", "1")

	// Set the audio codec to aac and bitrate to 64k
	args = append(args, "-c:a", "aac", "-b:a", "64k")

	// Set the output file path
	args = append(args, outputFile)

	// Create a new command using the "ffmpeg" executable and the arguments
	cmd := exec.Command("ffmpeg", args...)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Print the error and output if the command fails
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
	}

	// Return nil if the operation is successful
	return nil
}

// MakeSplitMP3Files splits an audiobook into MP3 files based on chapters.
func MakeSplitMP3Files(files []string, chapters []Chapter, meta Metadata, outputDir string) error {

	// Print a message indicating that the audiobook is being split into MP3 files
	fmt.Println("Splitting Audiobook into MP3 files based on chapters...")

	// Create a slice to store the command line arguments
	var args []string

	// Append the input file paths to the arguments slice using the "concat" format
	args = append(args, "-i", "concat:"+strings.Join(files, "|"))

	for i, chap := range chapters {
		// Adds the ffmpeg arguments for each chapter
		count := i + 1
		args = append(args, "-ss", fmt.Sprintf("%dms", chap.StartOffsetMs), "-t", fmt.Sprintf("%dms", chap.LengthMs))
		args = append(args, "-metadata", "title="+chap.Title, "-metadata", "artist="+meta.Author, "-metadata", "album="+meta.Title, "-metadata", fmt.Sprintf("track=%d", count))
		args = append(args, "-acodec", "copy")
		args = append(args, path.Join(outputDir, fmt.Sprintf("[%d]. ", count)+chap.Title+".mp3"))
	}

	// Create a new command using the "ffmpeg" executable and the arguments
	cmd := exec.Command("ffmpeg", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Print the error and the command output if an error occurs
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
	}

	// Return nil if the operation is successful
	return nil
}

// MakeSplitM4BFiles splits an audiobook into M4B files based on chapters.
//
// Parameters:
// - files: a list of input file paths.
// - chapters: a list of Chapter structs representing the chapters of the audiobook.
// - meta: a Metadata struct containing metadata about the audiobook.
// - outputDir: the directory where the output M4B files will be saved.
//
// Returns:
// - error: an error if any occurred during the splitting process.
func MakeSplitM4BFiles(files []string, chapters []Chapter, meta Metadata, outputDir string) error {
	// Print a message indicating that the audiobook is being split into M4B files
	fmt.Println("Splitting Audiobook into M4B files based on chapters...")

	// Iterate over the chapters
	for i, chap := range chapters {
		// Calculate the chapter count
		count := i + 1

		// Create a slice to store the command line arguments
		var args []string

		// Append the input file paths to the arguments slice using the "concat" format
		args = append(args, "-i", "concat:"+strings.Join(files, "|"))

		// Adds the ffmpeg arguments for each chapter
		args = append(args, "-ss", fmt.Sprintf("%dms", chap.StartOffsetMs), "-t", fmt.Sprintf("%dms", chap.LengthMs))
		args = append(args, "-metadata", "title="+chap.Title, "-metadata", "artist="+meta.Author, "-metadata", "album="+meta.Title, "-metadata", fmt.Sprintf("track=%d", count))
		args = append(args, "-acodec", "aac")
		args = append(args, path.Join(outputDir, fmt.Sprintf("[%d]. ", count)+chap.Title+".m4b"))

		// Create a new command using the "ffmpeg" executable and the arguments
		cmd := exec.Command("ffmpeg", args...)

		// Execute the command and capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Output:", string(output))
		}
	}

	// Return nil if the operation is successful
	return nil
}
