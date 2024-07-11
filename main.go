package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Define structs based on the RSS XML structure
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type APIResponse struct {
	Horoscope *Item `json:"horoscope"`
	Panchang  *Item `json:"panchang"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func fetchRSSFeed(url string) (RSS, error) {
	var rss RSS

	// Make HTTP GET request to fetch the RSS feed
	resp, err := http.Get(url)
	if err != nil {
		return rss, fmt.Errorf("error fetching RSS feed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rss, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the XML into RSS struct
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return rss, fmt.Errorf("error parsing XML: %v", err)
	}

	return rss, nil
}

func findTodayHoroscope(items []Item) *Item {
	todayDate := time.Now().Format("02-January-2006")
	todayDate = strings.ToLower(todayDate) // Format today's date as "11-July-2024"
	fmt.Println(todayDate)

	for _, item := range items {
		compareLink := strings.ToLower(item.Link)
		if strings.Contains(compareLink, "horoscope") && strings.Contains(compareLink, todayDate) {
			return &item
		}
	}
	return nil
}

func findTodayPanchang(items []Item) *Item {
	todayDate := time.Now().Format("02-January-2006")
	todayDate = strings.ToLower(todayDate) // Format today's date as "11-July-2024"
	fmt.Println(todayDate)

	for _, item := range items {
		compareLink := strings.ToLower(item.Link)
		if strings.Contains(compareLink, "panchang") && strings.Contains(compareLink, todayDate) {
			return &item
		}
	}
	return nil
}

func rashifalHandler(w http.ResponseWriter, r *http.Request) {
	// URL of the RSS feed
	url := "https://www.gujaratsamachar.com/rss/category/astro"
	fmt.Println("Project started...rashifalHandler")

	// Fetch the RSS feed data
	rss, err := fetchRSSFeed(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch RSS feed: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare JSON response
	type RashifalResponse struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Items       []Item `json:"items"`
	}

	response := RashifalResponse{
		Title:       rss.Channel.Title,
		Description: rss.Channel.Description,
		Items:       rss.Channel.Items,
	}

	// Encode JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	horoscope := findTodayHoroscope(response.Items)
	panchang := findTodayPanchang(response.Items)
	jsonResponse := APIResponse{
		Horoscope: horoscope,
		Panchang:  panchang,
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rashifal/divya", rashifalHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
	fmt.Println("Project started...")
}
