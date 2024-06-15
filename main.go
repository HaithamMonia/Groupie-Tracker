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
	ID           int
	Name         string
	Image        string
	FirstAlbum   string
	CreationDate int
	Members      []string
}

// HomePageVars contains variables to pass to the homepage template
type HomePageVars struct {
	Artists []Artist
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist/", ArtistHandler)

	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
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
	t.Execute(w, vars)
}

// fetchArtistDetails makes an HTTP request to the API and returns an Artist struct for the given ID
func fetchArtistDetails(id int) (Artist, error) {
	// Construct the URL with the artist ID
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return Artist{}, err
	}
	defer resp.Body.Close()

	// Decode the JSON response into an Artist struct
	var artist Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		return Artist{}, err
	}

	return artist, nil
}

// ArtistHandler handles requests for individual artist details
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting the artist ID from the URL path
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

	// Fetching artist details using the ID...
	// Assuming fetchArtistDetails is a function you've implemented
	artist, err := fetchArtistDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Assuming you have an 'artist.html' template for displaying individual artist details
	t, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, artist)
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
