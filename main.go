package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Artist represents the structure of an artist's data
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	FirstAlbum   string   `json:"firstAlbum"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
	LocationsURL string   `json:"locations"`
	DatesURL     string   `json:"concertDates"`
	RelationsURL string   `json:"relations"`
}

// ConcertDetails represents the structure of concert details data
type ConcertDetails struct {
	LocationsDates map[string][]string `json:"datesLocations"`
}

// HomePageVars contains variables to pass to the homepage template
type HomePageVars struct {
	Artists []Artist
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist/", ArtistHandler)

	log.Println("Listening on : http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// fetchArtists makes an HTTP request to the API and returns a slice of Artists
func fetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// HomeHandler handles the homepage requests and displays a list of artists
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := HomePageVars{Artists: artists}
	err = t.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fetchArtistDetails makes an HTTP request to the API and returns an Artist struct and ConcertDetails for the given ID
func fetchArtistDetails(id int) (Artist, ConcertDetails, error) {
	// Fetch artist details
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}
	defer resp.Body.Close()

	var artist Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}

	// Fetch relation details
	resp, err = http.Get(artist.RelationsURL)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}
	defer resp.Body.Close()

	var relations ConcertDetails
	err = json.NewDecoder(resp.Body).Decode(&relations)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}

	return artist, relations, nil
}

// ArtistHandler handles requests for individual artist details and their concert details
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Artist ID", http.StatusBadRequest)
		return
	}

	artist, concertDetails, err := fetchArtistDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a struct to hold both artist and concert details
	data := struct {
		Artist         Artist
		ConcertDetails ConcertDetails
	}{
		Artist:         artist,
		ConcertDetails: concertDetails,
	}

	// Logging data for debugging


	// Execute the template with the combined data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
