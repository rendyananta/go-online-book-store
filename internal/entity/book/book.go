package book

import "time"

type PaginationParam struct {
	PerPage int
	LastID  string
}

type PaginationResult struct {
	Data    []Book
	PerPage int
	LastID  string
}

type Publisher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Price            float64    `json:"price"`
	ISBN             string     `json:"isbn"`
	Language         string     `json:"language"`
	Edition          string     `json:"edition"`
	Pages            int        `json:"pages"`
	PublishedAt      *time.Time `json:"published_at"`
	FirstPublishedAt *time.Time `json:"first_published_at"`
	CoverImg         string     `json:"cover_img"`
	Rating           float64    `json:"rating"`
	Publisher        Publisher  `json:"publisher"`
	Authors          []Author   `json:"authors"`
	Genres           []Genre    `json:"genres"`
}
