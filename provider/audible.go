package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	meta "Z0y6h0kS9X/libby-chapterizer/pkg"
)

type Response struct {
	ProductFilters []interface{} `json:"product_filters"`
	Products       []Product     `json:"products"`
	ResponseGroups []string      `json:"response_groups"`
	TotalResults   int           `json:"total_results"`
}

type Product struct {
	ASIN string `json:"asin"`
}

// GetBook queries the Audible API to retrieve a book's ASIN based on its title, author, and narrator.
func GetBook(title, author, narrator string) (string, error) {
	// Create URL parameters
	params := url.Values{
		"num_results":      {"10"},
		"products_sort_by": {"Relevance"},
		"title":            {title},
		"author":           {author},
		"narrator":         {narrator},
	}

	// Encode parameters into query string
	queryString := params.Encode()

	// Create request URL
	requestURL := fmt.Sprintf("https://api.audible.com/1.0/catalog/products?%s", queryString)

	// Send HTTP GET request
	response, err := http.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	var rsp Response
	if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	// Check if any books were found
	if len(rsp.Products) == 0 {
		return "", nil
	} else {
		return rsp.Products[0].ASIN, nil
	}

}

// GetBookDetailsASIN retrieves the details of a book with the given ASIN.
func GetBookDetailsASIN(asin string) (meta.BookDetails, error) {

	// Construct the request URL
	requestURL := fmt.Sprintf("https://api.audnex.us/books/%s", asin)

	// Send an HTTP GET request to the API
	response, err := http.Get(requestURL)
	if err != nil {
		return meta.BookDetails{}, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	// Decode the JSON response into a BookDetails struct
	var rsp meta.BookDetails
	if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
		return meta.BookDetails{}, fmt.Errorf("error decoding response: %w", err)
	}

	return rsp, nil
}

// GetChapters retrieves the chapters for a given ASIN.
// It makes an HTTP GET request to the audnex API and decodes the response into a Chapters struct.
// The ASIN is used to construct the request URL.
func GetChapters(asin string) (meta.Chapters, error) {
	// Construct the request URL
	requestURL := fmt.Sprintf("https://api.audnex.us/books/%s/chapters", asin)

	// Send an HTTP GET request to the API
	response, err := http.Get(requestURL)
	if err != nil {
		return meta.Chapters{}, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	// Decode the JSON response into a Chapters struct
	var rsp meta.Chapters
	if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
		return meta.Chapters{}, fmt.Errorf("error decoding response: %w", err)
	}

	// Return the chapters
	return rsp, nil
}
