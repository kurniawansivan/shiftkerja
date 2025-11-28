package entity

// Shift represents a job posting with a location
type Shift struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	PayRate  float64 `json:"pay_rate"`
}