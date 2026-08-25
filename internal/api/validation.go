package api

import (
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	Limit  int
	Offset int
}

func parsePage(r *http.Request) Page {
	page := Page{Limit: 50}
	if value, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && value > 0 && value <= 200 {
		page.Limit = value
	}
	if value, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && value >= 0 {
		page.Offset = value
	}
	return page
}

func paginate[T any](items []T, page Page) []T {
	if page.Offset >= len(items) {
		return []T{}
	}
	end := page.Offset + page.Limit
	if end > len(items) {
		end = len(items)
	}
	return items[page.Offset:end]
}

func cleanPath(value string) string {
	value = strings.TrimSpace(value)
	return strings.Trim(value, "/")
}

func validMethod(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}
