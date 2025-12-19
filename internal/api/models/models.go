package models


type (
	Courses struct {
		Courseid int `json:"courseid"`
		Title string `json:"title"`
		Description string  `json:"description"`
	}

	Videos struct {
		Title string
		Description string
		Videourl string

	}
)