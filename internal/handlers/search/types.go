package search

type SearchQueryEvent struct {
	UserID   string
	ClientIP string
	Query    string
	Device   string
}
