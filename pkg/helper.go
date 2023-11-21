package pkg

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	// prov "Z0y6h0kS9X/libby-chapterizer/provider"
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

	if details.SeriesPrimary.Name != "" {
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
