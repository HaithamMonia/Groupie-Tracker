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

// Locations represents the structure of locations data
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

// Dates represents the structure of dates data
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// ConcertDetails represents the structure of concert details data
type ConcertDetails struct {
	Locations []string
	Dates     []string
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

	log.Println("Listening on :8080...")
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

	// Fetch locations details
	resp, err = http.Get(artist.LocationsURL)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}
	defer resp.Body.Close()

	var locations Locations
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}

	// Fetch dates details
	resp, err = http.Get(artist.DatesURL)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}
	defer resp.Body.Close()

	var dates Dates
	err = json.NewDecoder(resp.Body).Decode(&dates)
	if err != nil {
		return Artist{}, ConcertDetails{}, err
	}

	// Remove asterisks from dates
	for i, date := range dates.Dates {
		dates.Dates[i] = strings.TrimPrefix(date, "*")
	}

	concertDetails := ConcertDetails{
		Locations: locations.Locations,
		Dates:     dates.Dates,
	}

	return artist, concertDetails, nil
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
	log.Printf("Artist: %+v\n", artist)
	log.Printf("ConcertDetails: %+v\n", concertDetails)

	// Execute the template with the combined data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
