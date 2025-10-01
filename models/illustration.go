package models

type Illustration struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	FileName string `json:"file_name"`
}
