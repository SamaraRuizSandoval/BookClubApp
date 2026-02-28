package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type GoogleBooksAPISearch struct {
	Kind       string                `json:"kind,omitempty"`
	TotalItems int                   `json:"totalItems,omitempty"`
	Items      []GoogleBookBasicInfo `json:"items,omitempty"`
}

type GoogleBookBasicInfo struct {
	ID         string      `json:"id"`
	VolumeInfo *VolumeInfo `json:"volumeInfo,omitempty"`
}

type VolumeInfo struct {
	Title               string               `json:"title,omitempty"`
	Authors             []string             `json:"authors,omitempty"`
	PublishedDate       string               `json:"publishedDate,omitempty"`
	Description         string               `json:"description,omitempty"`
	IndustryIdentifiers []IndustryIdentifier `json:"industryIdentifiers,omitempty"`
	ImageLinks          *ImageLinks          `json:"imageLinks,omitempty"`
	PageCount           int                  `json:"pageCount,omitempty"`
}

type IndustryIdentifier struct {
	Type       string `json:"type,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail,omitempty"`
	Thumbnail      string `json:"thumbnail,omitempty"`
}

type GoogleBookAPIStore struct {
	apiKey string
}

func NewGoogleBooksStore() *GoogleBookAPIStore {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	if apiKey == "" {
		fmt.Println("WARNING: GOOGLE_BOOKS_API_KEY not set")
		return nil
	}

	return &GoogleBookAPIStore{
		apiKey: apiKey,
	}
}

type GoogleBookAPI interface {
	SearchGoogleBooks(query string) ([]GoogleBookBasicInfo, error)
}

func (s *GoogleBookAPIStore) SearchGoogleBooks(query string) ([]GoogleBookBasicInfo, error) {
	var googleBooks GoogleBooksAPISearch

	baseURL := "https://www.googleapis.com/books/v1/volumes"

	params := url.Values{}
	params.Add("q", query)
	params.Add("key", s.apiKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching Google books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleBooks); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return googleBooks.Items, nil
}
