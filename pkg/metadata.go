package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
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

type Chapters struct {
	Asin                 string `json:"asin,omitempty"`
	BrandIntroDurationMs int    `json:"brandIntroDurationMs,omitempty"`
	BrandOutroDurationMs int    `json:"brandOutroDurationMs,omitempty"`
	Chapters             []struct {
		LengthMs       int    `json:"lengthMs,omitempty"`
		StartOffsetMs  int    `json:"startOffsetMs,omitempty"`
		StartOffsetSec int    `json:"startOffsetSec,omitempty"`
		Title          string `json:"title,omitempty"`
	} `json:"chapters,omitempty"`
	IsAccurate       bool   `json:"isAccurate,omitempty"`
	Region           string `json:"region,omitempty"`
	RuntimeLengthMs  int    `json:"runtimeLengthMs,omitempty"`
	RuntimeLengthSec int    `json:"runtimeLengthSec,omitempty"`
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
	TotalMilliseconds float64
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

func (p Process) ToString() string {
	// return a string representation of the Process struct
	return fmt.Sprintf("Source: %s, Title: %s, Start: %f, End: %f", p.Source, p.Title, p.Start, p.End)
}

func (d Duration) ToString() string {
	// return a string representation of the Duration struct
	return fmt.Sprintf("%02d:%02d:%02d.%03d\n", d.Hours, d.Minutes, d.Seconds, d.Milliseconds)
}

func (o Openbook) CalculateRuntime() int {

	// Adds up the audio-duration fields of the spine and returns the total
	var totalDuration float64

	// Grabs the audio duration (in seconds) of each item in the spine
	for _, item := range o.Spine {
		totalDuration += item.AudioDuration
	}

	// Converts the seconds into minutes
	totalDuration = totalDuration / 60

	// Returns the total duration in Minutes (to match audnexus format)
	return int(totalDuration)

}

// type Author struct {
// 	Asin string `json:"asin"`
// 	Name string `json:"name"`
// }

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
}

func GetMetadataFromASIN(asin string) (Metadata, error) {

	// Starts with an empty Metadata struct
	var metadata Metadata

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

	// Gets the primary authors name and assigns it
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

		metadata.Series.Position, err = strconv.ParseFloat(fmt.Sprint(seriesObj["position"]), 64)
		if err != nil {
			return metadata, fmt.Errorf("error decoding series position")
		}

	}

	// Gets the publisher and assigns it
	var duration Duration
	if mins, ok := rsp["runtimeLengthMin"].(int); ok {
		duration = CalculateDuration(mins)
	} else if mins, ok := rsp["runtimeLengthMin"].(float64); ok {
		duration = CalculateDuration(int(mins))
	} else {
		fmt.Println("Invalid runtimeLengthMin type")
	}

	metadata.Duration = duration

	return metadata, nil

	// // Construct the request URL
	// requestURL = fmt.Sprintf("https://api.audnex.us/books/%s/chapters", asin)

	// // Send an HTTP GET request to the API
	// response, err = http.Get(requestURL)
	// if err != nil {
	// 	return meta.Chapters{}, fmt.Errorf("error making request: %w", err)
	// }
	// defer response.Body.Close()

	// // Decode the JSON response into a Chapters struct
	// var rsp meta.Chapters
	// if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
	// 	return meta.Chapters{}, fmt.Errorf("error decoding response: %w", err)
	// }

	// // Return the chapters
	// return rsp, nil

}

func (m Metadata) ToString() string {

	return fmt.Sprintf("ASIN:      %s\nTitle:     %s\nAuthor:    %s\nSeries:    %s\nPosition:  %f\nPublisher: %s\nDuration:  %s", m.ASIN, m.Title, m.Author, m.Series.Name, m.Series.Position, m.Publisher, m.Duration.ToString())

}
