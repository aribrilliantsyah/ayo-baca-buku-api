package models

type DashboardResponse struct {
	TotalBooksReading  int64 `json:"total_books_reading"`
	TotalBooksFinished int64 `json:"total_books_finished"`
	TotalPagesRead     int64 `json:"total_pages_read"`
}
