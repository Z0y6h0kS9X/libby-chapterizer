package pkg

import (
	"fmt"
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
	Title  string
	Source string
	Output string
	Start  float64
	End    float64
}

func (p Process) ToString() string {
	// return a string representation of the Process struct
	return fmt.Sprintf("Source: %s, Title: %s, Start: %f, End: %f", p.Source, p.Title, p.Start, p.End)
}
