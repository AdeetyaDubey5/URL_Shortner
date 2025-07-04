package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel,omitempty"` // "rel" describes the relationship (e.g., "self", "redirect")
	Type string `json:"type,omitempty"` // Optional: "application/json" or "text/html" etc.
}

// ShortenResponse is the structure for the API response
type ShortenResponse struct {
	ShortCode string `json:"short_code"` // Renamed from ShortURL for clarity in response
	Links     []Link `json:"_links,omitempty"` // HATEOAS links
}
var urlDB = make(map[string]URL)

/*
	"0001" --> {
				ID : "0001"
				OriginalURL : "https://www.google.com"
				ShortURL : "https://www.shorturl.com/0001"
				CreationDate : "2022-01-01 00:00:00"
				}
*/
func URL_shortner(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // Convert the Original URL into byte slice
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]		// 8 characters is a common length seen in many short URL services.
}

func createShortURL(OriginalURL string) string {
	shortURL := URL_shortner(OriginalURL)
	id := shortURL // Using the shortened url as our ID
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getOriginalURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return  url , nil
}

func handler(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintln(w,"Handling request")
}

func shortURL_handler(w http.ResponseWriter,r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortCode := createShortURL(data.URL) // Get the generated short code

	// Construct the full redirect URL
	redirectHref := fmt.Sprintf("http://localhost:3000/redirect/%s", shortCode)

	// Create the HATEOAS response structure
	response := ShortenResponse{
		ShortCode: shortCode,
		Links: []Link{
			{Href: redirectHref, Rel: "redirect_to_original", Type: "text/html"},
			{Href: "http://localhost:3000/shorten", Rel: "self", Type: "application/json"}, // Link to the current resource/action
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler (w http.ResponseWriter,r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url,err := getOriginalURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
		}
	http.Redirect(w,r,url.OriginalURL,http.StatusFound)	
}
 
func main() {
	fmt.Println("Starting URL Shortner")
	OriginalURL := "https://gemini.google.com/app/0f81ef95d7aacd96"
	ShortURL := URL_shortner(OriginalURL)
	fmt.Println("Original URL: ", OriginalURL)
	fmt.Println("Shortend URL: ", ShortURL)

	// Register the `Handler` func which will handle the requests to the root URL("/")
	http.HandleFunc("/",handler)
	http.HandleFunc("/shorten",shortURL_handler)
	http.HandleFunc("/redirect/",redirectURLHandler)


	// Start the http server on port 5000
	fmt.Println("Starting serveon Port:3000")
	err := http.ListenAndServe(":3000",nil) // Starting the server 
	if err!=nil {
		fmt.Println("Error starting server:",err)
	}

	
	
}
