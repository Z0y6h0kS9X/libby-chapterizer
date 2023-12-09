package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Playlist struct {
	Title     string
	Author    string
	BookTitle string
	Tracks    []Track
}

type Track struct {
	Title    string
	Length   int
	LengthMS int
	File     string
}

func ReadExtM3U(fileName string) (Playlist, error) {

	var f io.ReadCloser
	var play Playlist

	file, err := os.Open(fileName)
	if err != nil {
		return play,
			fmt.Errorf("unable to open playlist file: %v", err)
	}
	f = file

	defer f.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		// Extract the Playlist Title
		if strings.HasPrefix(scanner.Text(), "#PLAYLIST") {
			item := strings.Split(scanner.Text(), ": ")
			play.Title = item[1]
		}

		if strings.HasPrefix(scanner.Text(), "#EXTART") {
			item := strings.Split(scanner.Text(), ": ")
			play.Author = item[1]
		}

		if strings.HasPrefix(scanner.Text(), "#EXTALB") {
			item := strings.Split(scanner.Text(), ": ")
			play.BookTitle = item[1]
		}

		if strings.HasPrefix(scanner.Text(), "#EXTINF") {
			data := strings.Replace(scanner.Text(), "#EXTINF:", "", -1)
			items := strings.Split(data, ",")
			title := items[1]
			times := strings.Split(items[0], " ")
			length, err := strconv.Atoi(times[0])
			if err != nil {
				return play, err
			}
			lengthMS, err := strconv.Atoi(strings.Replace(times[1], "ms=", "", -1))
			if err != nil {
				return play, err
			}

			var file string
			if scanner.Scan() {
				file = scanner.Text()
			}

			track := Track{Title: title, Length: length, LengthMS: lengthMS, File: file}
			play.Tracks = append(play.Tracks, track)

		}
	}

	return play, nil
}

func (p Playlist) WriteM3UFile(fileName string) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString("#EXTM3U\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("#PLAYLIST: " + p.Title + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("#EXTART: " + p.Author + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("#EXTALB: " + p.BookTitle + "\n")
	if err != nil {
		return err
	}

	for _, track := range p.Tracks {

		_, err = file.WriteString(fmt.Sprintf("#EXTINF:%d ms=%d,%s\n", track.Length, track.LengthMS, track.Title) + track.File + "\n")
		if err != nil {
			return err
		}
	}

	file.Close()
	return nil

}

func (p Playlist) WriteFFMPEGMetadataFile(fileName string) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(";FFMETADATA1\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("title=" + p.Title + "\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("artist=" + p.Author + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	duration := 0

	for _, track := range p.Tracks {

		_, err = file.WriteString("[CHAPTER]\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("TIMEBASE=1/1000\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString(fmt.Sprintf("START=%d\n", duration))
		if err != nil {
			return err
		}
		_, err = file.WriteString(fmt.Sprintf("END=%d\n", duration+track.LengthMS))
		if err != nil {
			return err
		}
		_, err = file.WriteString("title=" + track.Title + "\n")
		if err != nil {
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}

		duration += track.LengthMS

	}

	file.Close()

	return nil
}
