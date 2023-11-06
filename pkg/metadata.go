package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Book struct {
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

type Process struct {
	Title  string
	Source string
	Output string
	Start  int
	End    int
}

func (p Process) ToString() string {
	// return a string representation of the Process struct
	return fmt.Sprintf("Source: %s, Title: %s, Start: %d, End: %d", p.Source, p.Title, p.Start, p.End)
}

func GetFileNameAndSeconds(path string) (string, int) {

	var seconds int

	part := strings.Split(path, "Fmt425-")[1]
	// Splits the path on "#" to enumerate the seconds (if applicable) and match to the file paths above
	if strings.Contains(part, "#") {
		result := strings.Split(part, "#")
		part = result[0]
		seconds, _ = strconv.Atoi(result[1])
	}

	return part, seconds

}
