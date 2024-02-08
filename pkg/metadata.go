package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Openbook struct {
	Cover struct {
		Front struct {
			MediaType              string    `json:"media-type,omitempty"`
			OdreadAspectRatio      float64   `json:"-odread-aspect-ratio,omitempty"`
			OdreadColor            []int     `json:"-odread-color,omitempty"`
			OdreadFileBytes        int       `json:"-odread-file-bytes,omitempty"`
			OdreadFileLastModified time.Time `json:"-odread-file-last-modified,omitempty"`
			OdreadHeight           int       `json:"-odread-height,omitempty"`
			OdreadWidth            int       `json:"-odread-width,omitempty"`
		} `json:"front,omitempty"`
	} `json:"cover,omitempty"`
	Creator []struct {
		Bio  string `json:"bio,omitempty"`
		Name string `json:"name,omitempty"`
		Role string `json:"role,omitempty"`
	} `json:"creator,omitempty"`
	Description struct {
		Full  string `json:"full,omitempty"`
		Short string `json:"short,omitempty"`
	} `json:"description,omitempty"`
	Language string `json:"language,omitempty"`
	Nav      struct {
		Toc []struct {
			Path  string `json:"path,omitempty"`
			Title string `json:"title,omitempty"`
		} `json:"toc,omitempty"`
	} `json:"nav,omitempty"`
	OdreadAnchor                any      `json:"-odread-anchor,omitempty"`
	OdreadBankScope             string   `json:"-odread-bank-scope,omitempty"`
	OdreadBankVerificationToken string   `json:"-odread-bank-verification-token,omitempty"`
	OdreadBonafidesD            string   `json:"-odread-bonafides-d,omitempty"`
	OdreadBonafidesM            any      `json:"-odread-bonafides-m,omitempty"`
	OdreadBonafidesP            string   `json:"-odread-bonafides-p,omitempty"`
	OdreadBonafidesS            any      `json:"-odread-bonafides-s,omitempty"`
	OdreadBuid                  string   `json:"-odread-buid,omitempty"`
	OdreadCoverColor            []int    `json:"-odread-cover-color,omitempty"`
	OdreadCoverRatio            float64  `json:"-odread-cover-ratio,omitempty"`
	OdreadCrid                  []string `json:"-odread-crid,omitempty"`
	OdreadFurbishURI            string   `json:"-odread-furbish-uri,omitempty"`
	OdreadMsgAccess             string   `json:"-odread-msg-access,omitempty"`
	OdreadMsgExpires            int      `json:"-odread-msg-expires,omitempty"`
	OdreadMsgSync               bool     `json:"-odread-msg-sync,omitempty"`
	OdreadUilinks               struct {
		Dictionarycom string `json:"DICTIONARYCOM,omitempty"`
		Helpaud       string `json:"HELPAUD,omitempty"`
		Helpcompat    string `json:"HELPCOMPAT,omitempty"`
		Helpolm       string `json:"HELPOLM,omitempty"`
		Helpwip       string `json:"HELPWIP,omitempty"`
		Lexisnexis    string `json:"LEXISNEXIS,omitempty"`
		Odhome        string `json:"ODHOME,omitempty"`
		Odreadinfo    string `json:"ODREADINFO,omitempty"`
		Odsearch      string `json:"ODSEARCH,omitempty"`
		Sspulseit     string `json:"SSPULSEIT,omitempty"`
		Ssxoxo        string `json:"SSXOXO,omitempty"`
	} `json:"-odread-uilinks,omitempty"`
	RenditionFormat string `json:"rendition-format,omitempty"`
	Spine           []struct {
		AudioBitrate        int     `json:"audio-bitrate,omitempty"`
		AudioDuration       float64 `json:"audio-duration,omitempty"`
		MediaType           string  `json:"media-type,omitempty"`
		OdreadFileBytes     int     `json:"-odread-file-bytes,omitempty"`
		OdreadOriginalPath  string  `json:"-odread-original-path,omitempty"`
		OdreadSpinePosition int     `json:"-odread-spine-position,omitempty"`
		Path                string  `json:"path,omitempty"`
	} `json:"spine,omitempty"`
	Title struct {
		Collection string `json:"collection,omitempty"`
		Main       string `json:"main,omitempty"`
		Subtitle   string `json:"subtitle,omitempty"`
	} `json:"title,omitempty"`
}

type BookDetails struct {
	Asin    string `json:"asin,omitempty"`
	Authors []struct {
		Asin string `json:"asin,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"authors,omitempty"`
	Description string `json:"description,omitempty"`
	FormatType  string `json:"formatType,omitempty"`
	Genres      []struct {
		Asin string `json:"asin,omitempty"`
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"genres,omitempty"`
	Image     string `json:"image,omitempty"`
	IsAdult   bool   `json:"isAdult,omitempty"`
	Language  string `json:"language,omitempty"`
	Narrators []struct {
		Name string `json:"name,omitempty"`
	} `json:"narrators,omitempty"`
	PublisherName    string    `json:"publisherName,omitempty"`
	Rating           string    `json:"rating,omitempty"`
	Region           string    `json:"region,omitempty"`
	ReleaseDate      time.Time `json:"releaseDate,omitempty"`
	RuntimeLengthMin int       `json:"runtimeLengthMin,omitempty"`
	SeriesPrimary    struct {
		Asin     string `json:"asin,omitempty"`
		Name     string `json:"name,omitempty"`
		Position string `json:"position,omitempty"`
	} `json:"seriesPrimary,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Summary  string `json:"summary,omitempty"`
	Title    string `json:"title,omitempty"`
}

type Chapter struct {
	LengthMs       int    `json:"lengthMs,omitempty"`
	StartOffsetMs  int    `json:"startOffsetMs,omitempty"`
	StartOffsetSec int    `json:"startOffsetSec,omitempty"`
	Title          string `json:"title,omitempty"`
}

type Chapters struct {
	Asin                 string    `json:"asin,omitempty"`
	BrandIntroDurationMs int       `json:"brandIntroDurationMs,omitempty"`
	BrandOutroDurationMs int       `json:"brandOutroDurationMs,omitempty"`
	Chapters             []Chapter `json:"chapters,omitempty"`
	IsAccurate           bool      `json:"isAccurate,omitempty"`
	Region               string    `json:"region,omitempty"`
	RuntimeLengthMs      int       `json:"runtimeLengthMs,omitempty"`
	RuntimeLengthSec     int       `json:"runtimeLengthSec,omitempty"`
}

type Process struct {
	Title       string
	Source      string
	Output      string
	Start       float64
	End         float64
	DurationStr string
	Duration    Duration
	Command     *exec.Cmd
}

type Duration struct {
	Hours             int
	Minutes           int
	Seconds           int
	Milliseconds      int
	TotalMinutes      int
	TotalSeconds      int
	TotalMilliseconds int
}

type M3U struct {
	PlaylistTitle string
	Author        string
	Items         []struct {
		Length   int
		LengthMS int
		Title    string
		FileName string
	}
}

type Metadata struct {
	ASIN   string
	Title  string
	Series struct {
		Name     string
		Position float64
	}
	Author    string
	Narrator  string
	Publisher string
	Duration  Duration
	Summary   string
	Abridged  bool
	Chapters  []Chapter
}

// ToString returns a string representation of the Process struct.
func (p Process) ToString() string {
	return fmt.Sprintf("Source: %s, Title: %s, Start: %f, End: %f", p.Source, p.Title, p.Start, p.End)
}

// ToString returns a string representation of the Duration struct.
func (d Duration) ToString() string {
	// Format the hours, minutes, seconds, and milliseconds into a string.
	// The string is formatted as "HH:MM:SS.MMM".
	return fmt.Sprintf("%02d:%02d:%02d.%03d", d.Hours, d.Minutes, d.Seconds, d.Milliseconds)
}

// CalculateRuntime calculates the total duration of the audio in the Openbook.
// It sums up the audio duration of each item in the spine and returns the total duration in minutes.
func (o Openbook) CalculateRuntime() int {
	totalDuration := 0.0

	// Iterate over each item in the spine
	for _, item := range o.Spine {
		totalDuration += item.AudioDuration
	}

	// Convert the total duration from seconds to minutes
	totalDuration = totalDuration / 60

	// Return the total duration in minutes
	return int(totalDuration)
}

// GetMetadataFromASIN retrieves metadata for a book based on its ASIN.
func GetMetadataFromASIN(asin string) (Metadata, error) {
	metadata := Metadata{} // Starts with an empty Metadata struct

	// Construct the request URL for the top level metadata
	requestURL := fmt.Sprintf("https://api.audnex.us/books/%s", asin)

	// Send an HTTP GET request to the API
	response, err := http.Get(requestURL)
	if err != nil {
		return metadata, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	var rsp map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
		return metadata, fmt.Errorf("error decoding response: %w", err)
	}

	// Gets the ASIN and assigns it
	if asin, ok := rsp["asin"].(string); ok {
		metadata.ASIN = asin
	}

	// Gets the title and assigns it
	if title, ok := rsp["title"].(string); ok {
		metadata.Title = title
	}

	// Gets the publisher and assigns it
	if publisher, ok := rsp["publisherName"].(string); ok {
		metadata.Publisher = publisher
	}

	// Gets the summary and assigns it
	if summary, ok := rsp["summary"].(string); ok {
		metadata.Summary = summary
	}

	// Gets the abridged status and assigns it
	if abridged, ok := rsp["abridged"].(string); ok {
		switch abridged {
		case "abridged":
			metadata.Abridged = true
		case "unabridged":
			metadata.Abridged = false
		default:
			return metadata, fmt.Errorf("error decoding abridged status")
		}
	}

	// Gets the primary author's name and assigns it
	if authors, ok := rsp["authors"].([]interface{}); ok && len(authors) > 0 {
		authorObj, ok := authors[0].(map[string]interface{})
		if !ok {
			return metadata, fmt.Errorf("error decoding author object")
		}
		authorName, ok := authorObj["name"].(string)
		if !ok {
			return metadata, fmt.Errorf("error decoding author name")
		}
		metadata.Author = authorName
		// Do something with the first author's name
	}

	// Gets the series name and position and assigns it
	if seriesObj, ok := rsp["seriesPrimary"].(map[string]interface{}); ok {
		seriesName, ok := seriesObj["name"].(string)
		if !ok {
			return metadata, fmt.Errorf("error decoding series name")
		}
		metadata.Series.Name = seriesName

		// Some metadata positions are labeled with a leading word (e.g. Eragon - 'Book 1'), selects only the numbers
		posRegex := regexp.MustCompile(`\d+(\.\d+)?`)
		number := posRegex.FindAllString(fmt.Sprint(seriesObj["position"]), 1)

		// Converts the number to a float, if position is supplied (Ballad of Songbirds & Snakes has no position)
		if len(number) != 0 {

			metadata.Series.Position, err = strconv.ParseFloat(number[0], 64)
			if err != nil {
				return metadata, fmt.Errorf("error decoding series position")
			}

		}

	}

	// Gets the publisher and assigns it
	var duration Duration
	if mins, ok := rsp["runtimeLengthMin"].(int); ok {
		// Converts minutes to milliseconds and generates duration
		duration = CalculateDuration(mins * 60000)
	} else if mins, ok := rsp["runtimeLengthMin"].(float64); ok {
		// Converts minutes to milliseconds and generates duration
		duration = CalculateDuration(int(mins * 60000))
	} else {
		fmt.Println("Invalid runtimeLengthMin type")
	}

	metadata.Duration = duration

	// Returns the metadata
	return metadata, nil
}

// GetMetadataLocal converts the given openbook object to a metadata object.
// It extracts the title, primary author, primary narrator, summary, and series information from the openbook object.
// Returns the metadata object and an error if any.
func GetMetadataLocal(openbook Openbook) (Metadata, error) {

	// Create a new metadata object
	metadata := Metadata{}

	// Extract the title from the openbook and assign it to the metadata
	metadata.Title = openbook.Title.Main

	// Extract the primary author from the openbook and assign it to the metadata
	metadata.Author = GetPrimaryAuthor(openbook)

	// Extract the primary narrator from the openbook and assign it to the metadata
	metadata.Narrator = GetPrimaryNarrator(openbook)

	// Extract the summary from the openbook and assign it to the metadata
	metadata.Summary = openbook.Description.Short

	// Extract the series name from the openbook and assign it to the metadata
	metadata.Series.Name = openbook.Title.Collection

	// Return the metadata object and no error
	return metadata, nil
}

// GetChaptersLocal retrieves the chapters of a book and their durations from local mp3 files.
func GetChaptersLocal(book Openbook, mp3s []string) ([]Chapter, error) {
	// chapters will store the information about each chapter
	var chapters []ChapterInfo

	// Iterate over each item in the table of contents
	for i, item := range book.Nav.Toc {
		// Get the part and milliseconds from the item's path
		part, milliseconds := GetFileNameAndMilliseconds(item.Path)

		// Look up the part against the list of mp3s
		var path string
		var fileDuration Duration
		for _, mp3 := range mp3s {
			// Check if the mp3 contains the part
			if strings.Contains(mp3, part) {
				path = mp3

				// Get the duration of the mp3 file
				milli, err := GetFileDurationMS(path)
				if err != nil {
					log.Fatalf("Error getting duration: %v", err)
				} else {
					path = ""
					fileDuration = Duration{}
				}

				fileDuration = CalculateDuration(milli)
			}
		}

		// Create a ChapterInfo object with the retrieved information
		chp := ChapterInfo{
			ID:         (i + 1),
			Title:      item.Title,
			Start:      milliseconds,
			Duration:   -1,
			FilePath:   part,
			FileLength: fileDuration,
		}

		// Append the chapter to the chapters slice
		chapters = append(chapters, chp)
	}

	// Iterate over the chapters to calculate the durations
	for i, item := range chapters {
		// Check if there is a next chapter
		if i < len(chapters)-1 {
			next := chapters[i+1]

			// Check if the current chapter and the next chapter have the same file path
			if item.FilePath == next.FilePath {
				chapters[i].Duration = next.Start - item.Start
			} else {
				// Calculate the duration from the current chapter's start to the end of the file
				// and add the start time of the next file
				chapters[i].Duration = (item.FileLength.TotalMilliseconds - item.Start) + next.Start
			}
		} else if i == len(chapters)-1 {
			// Calculate the duration of the last chapter
			chapters[i].Duration = item.FileLength.TotalMilliseconds - item.Start
		}
	}

	// Convert the ChapterInfo objects to Chapter objects
	var chps []Chapter
	var totalDuration int
	for i, item := range chapters {
		var start int

		// Calculate the start offset for each chapter
		if i == 0 {
			start = 0
		} else {
			start = totalDuration
		}

		// Create a Chapter object with the calculated information
		chp := Chapter{
			LengthMs:       item.Duration,
			StartOffsetMs:  start,
			StartOffsetSec: start / 1000,
			Title:          strings.Replace(item.Title, `"`, "", -1),
		}

		totalDuration += item.Duration

		// Append the chapter to the chps slice
		chps = append(chps, chp)
	}

	return chps, nil
}

// ToString returns a string representation of the Metadata struct.
func (m Metadata) ToString() string {
	// Format the metadata fields into a string using fmt.Sprintf().
	// Each field is formatted with a specific format specifier.
	return fmt.Sprintf(
		"ASIN:      %s\n"+
			"Title:     %s\n"+
			"Author:    %s\n"+
			"Series:    %s\n"+
			"Position:  %f\n"+
			"Publisher: %s\n"+
			"Chapters:  %d\n"+
			"Duration:  %s\n"+
			"Abridged:  %t\n"+
			"Summary:   %s",
		m.ASIN, m.Title, m.Author, m.Series.Name, m.Series.Position,
		m.Publisher, len(m.Chapters), m.Duration.ToString(), m.Abridged, m.Summary,
	)
}

// ToFFMPEGMetadata converts the Metadata struct to a string representation of FFmpeg metadata.
func (m Metadata) ToFFMPEGMetadata() string {
	// Initialize the metadata string with the FFmpeg metadata version.
	metadata := ";FFMETADATA1\n"

	// Add the title metadata.
	metadata += "title=" + m.Title + "\n"

	// Add the series metadata.
	metadata += "series=" + m.Series.Name + "\n"

	// Add the number metadata with formatted float value.
	metadata += fmt.Sprintf("number=%f\n", m.Series.Position)

	// Add the author metadata.
	metadata += "author=" + m.Author + "\n"

	// Add the publisher metadata.
	metadata += "publisher=" + m.Publisher + "\n"

	// Add the ASIN metadata if it is not empty.
	if m.ASIN != "" {
		metadata += "asin=" + m.ASIN + "\n"
	}

	// Add a new line for separation.
	metadata += "\n"

	// Add the chapter metadata for each chapter in the list.
	for _, chapter := range m.Chapters {
		// Add the chapter header.
		metadata += "[CHAPTER]\n"

		// Add the chapter timebase.
		metadata += "TIMEBASE=1/1000\n"

		// Add the chapter start offset.
		metadata += "START=" + strconv.Itoa(chapter.StartOffsetMs) + "\n"

		// Add the chapter end offset.
		metadata += "END=" + strconv.Itoa(chapter.StartOffsetMs+chapter.LengthMs) + "\n"

		// Add the chapter title.
		metadata += "title=" + chapter.Title + "\n"

		// Add a new line for separation.
		metadata += "\n"
	}

	// Return the generated metadata string.
	return metadata
}
