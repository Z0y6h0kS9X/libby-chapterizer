package pkg

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func NormalizeFileName(filename string) string {

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
