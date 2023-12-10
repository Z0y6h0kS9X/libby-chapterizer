// This file is responsible for executing ffmp... commands.

package pkg

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetSimpleSplit(file string, start, end float64) (*exec.Cmd, error) {

	cmd := exec.Command("ffmpeg",
		"-i", file,
		"-ss", fmt.Sprintf("%f", start),
		"-to", fmt.Sprintf("%f", end),
		"-acodec", "copy",
		"-hide_banner", "-loglevel", "error")

	return cmd, nil
}

// GetComplexSplit returns a *exec.Cmd that represents the command for splitting and merging audio files using ffmpeg.
// It takes in two file paths, the start and end time for the split, and returns the command for splitting and merging the files.
func GetComplexSplit(file1, file2 string, seconds, seconds2 float64) (*exec.Cmd, error) {

	// Get the bitrate of file1
	bitrate, err := GetBitRate(file1)
	if err != nil {
		log.Println("Error getting bitrate:", err)
		return nil, err
	}

	// Get the bitrate of file2
	bitrate2, err := GetBitRate(file2)
	if err != nil {
		log.Println("Error getting bitrate:", err)
		return nil, err
	}

	// Determine the larger bitrate
	bitrateLarger := bitrate
	if bitrate2 > bitrateLarger {
		bitrateLarger = bitrate2
	}

	// Create the ffmpeg command with the necessary arguments
	cmd := exec.Command(
		"ffmpeg",
		"-i", file1,
		"-i", file2,
		"-filter_complex",
		fmt.Sprintf("[0:a]atrim=start=%.2f[a1];[1:a]atrim=start=0:end=%.2f[a2];[a1][a2]concat=n=2:v=0:a=1[out]", seconds, seconds2),
		"-map", "[out]",
		"-b:a", fmt.Sprintf("%d", bitrateLarger),
	)

	// Returns the command
	return cmd, nil
}

// GetFileDuration retrieves the duration of a file using ffprobe.
// It takes the filepath as input and returns the duration (in milliseconds) as a float64 value.
// If there is an error in running the ffprobe command or parsing the duration, an error is returned.
func GetFileDuration(filepath string) (float64, error) {
	// Create a new exec.Command with the ffprobe command and arguments
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filepath)

	// Run the command and capture the stdout
	stdout, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error running ffprobe command: %w", err)
	}

	// Trim any leading or trailing whitespace from the stdout
	durationStr := strings.TrimSpace(string(stdout))

	// Parse the duration string as a float64 value
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}

	// Return the duration
	return duration, nil

}

func GetFileDurationMS(filepath string) (int, error) {

	// Create a new exec.Command with the ffprobe command and arguments
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filepath)

	// Run the command and capture the stdout
	stdout, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error running ffprobe command: %w", err)
	}

	// Trim any leading or trailing whitespace from the stdout and convert it to a float
	float, err := strconv.ParseFloat(strings.TrimSpace(string(stdout)), 64)
	if err != nil {
		fmt.Println()
	}

	// Converts the duration to milliseconds
	var duration int = int(float * 1000)

	// Return the duration
	return duration, nil

}

// MakeCombinedMP3 concatenates multiple MP3 files into a single output file.
// It takes a slice of file paths and the path of the output file as input.
// It returns an error if the operation fails.
func MakeCombinedMP3(files []string, details BookDetails, outputFile string) error {

	fmt.Println("Making Combined MP3...")

	// Create a slice to store the command line arguments
	var args []string

	// Append the input file paths to the arguments slice using the "concat" format
	args = append(args, "-i", "concat:"+strings.Join(files, "|"))

	// Adds simple metadata to the output file
	args = append(args, "-metadata", "title="+details.Title)
	args = append(args, "-metadata", "artist="+details.Authors[0].Name)
	args = append(args, "-metadata", "album="+details.Title)
	args = append(args, "-metadata", "ASIN="+details.Asin)

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
