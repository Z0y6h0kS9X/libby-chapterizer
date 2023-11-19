package pkg

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetDurationFormatted(seconds float64) string {

	lengthRaw := time.Duration(seconds) * time.Second
	length := fmt.Sprintf("%02d:%02d:%02d.%03d",
		int(lengthRaw.Hours()),
		int(lengthRaw.Minutes())%60,
		int(lengthRaw.Seconds())%60,
		int(lengthRaw.Milliseconds()))

	length = strings.TrimRight(length, "0")

	return length

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
