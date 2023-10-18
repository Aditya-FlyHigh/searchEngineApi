package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

// Define a struct to represent the root object in data.json
type Data struct {
	WebResults []SearchResult `json:"webresults"`
}

// Define a struct to represent a search result
type SearchResult struct {
	Logo        string `json:"logo"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func searchResults(query string, results []SearchResult) []SearchResult {
	var searchResults []SearchResult

	if query != "" {
		for _, result := range results {
			if containsQueryInSearchResult(result, query) {
				searchResults = append(searchResults, result)
			}
		}

		if len(searchResults) > 0 {
			return searchResults
		}
	}

	return nil
}

func containsQueryInSearchResult(result SearchResult, query string) bool {
	// Check if the query is present in the Title and Description fields.
	if strings.Contains(strings.ToLower(result.Title), strings.ToLower(query)) ||
		strings.Contains(strings.ToLower(result.Description), strings.ToLower(query)) {
		return true
	}

	return false
}

// getSuggestions function returns search suggestions based on the query
func getSuggestions(query string, results []SearchResult) []string {
	var suggestions []string
	query = strings.ToLower(query)

	for _, result := range results {
		if strings.Contains(strings.ToLower(result.Title), query) {
			suggestions = append(suggestions, result.Title)
		}
	}

	return suggestions
}

func main() {
	// Read the JSON data from data.json
	data, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading data from data.json: %v", err)
	}
	var jsonData Data
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatalf("Error parsing data from data.json: %v", err)
	}

	results := jsonData.WebResults

	mux := http.NewServeMux()
	// API endpoint to handle search query
	// Format of API end points
	// http://localhost:8080/search/stanford
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimPrefix(r.URL.Path, "/search/") // Extract the query from the URL path

		// Get search results based on the user's query
		searchResults := searchResults(query, results)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(searchResults)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/suggestions/", func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimPrefix(r.URL.Path, "/suggestions/") // Extract the query from the URL path

		// Get search suggestions based on the user's query
		suggestions := getSuggestions(query, results)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(suggestions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	handler := cors.Default().Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
