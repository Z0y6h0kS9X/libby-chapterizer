package provider

import (
	"encoding/json"
	"fmt"
	"math"
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
func GetBook(title, author, narrator string, duration int) (string, error) {

	// Create the ASIN
	var asin string

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
		return asin, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	var rsp Response
	if err := json.NewDecoder(response.Body).Decode(&rsp); err != nil {
		return asin, fmt.Errorf("error decoding response: %w", err)
	}

	// Check if any books were found
	if len(rsp.Products) == 0 {
		fmt.Println("No books found")
		return asin, nil

	} else if len(rsp.Products) > 1 {

		// Goes through each, performing an Audnexus API call for each to match the duration
		m := make(map[string]int)
		for _, item := range rsp.Products {
			details, err := GetBookDetailsASIN(item.ASIN)
			if err != nil {
				return asin, fmt.Errorf("error getting book details: %w", err)
			}

			// If there is an exact match, return the ASIN - else store it in a map
			if details.RuntimeLengthMin == duration {
				fmt.Println("Exact match found!")
				return item.ASIN, nil
			} else {
				m[item.ASIN] = details.RuntimeLengthMin
			}

		}

		for a, runtime := range m {
			// Gets difference between duration and runtime
			diff := math.Abs(float64(duration - runtime))
			if diff <= 2 {
				fmt.Println("Close match found!, duration differs by less than 2 minutes")
				fmt.Println("ASIN: ", a)
				fmt.Println("Runtime: ", runtime)
				fmt.Println("Duration: ", duration)
				fmt.Println("Difference: ", diff)
				asin = a
				break
			}
		}

	} else {
		// Only 1 book was found
		asin = rsp.Products[0].ASIN
	}

	return asin, nil

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
