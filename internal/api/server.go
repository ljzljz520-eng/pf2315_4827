package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/scienceweekly/internal/flow001"
	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
)

type Server struct {
	flow  *flow001.Service
	query registry.Query
	now   time.Time
}

func New(flow *flow001.Service, query registry.Query, now time.Time) *Server {
	return &Server{flow: flow, query: query, now: now}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.recordAction)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "儿童科学实验周刊"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		filter := model.SearchFilter{Text: r.URL.Query().Get("q"), Edition: r.URL.Query().Get("edition"), IncludeArchived: r.URL.Query().Get("archived") == "true"}
		items, err := s.query.Search(filter)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, paginate(items, parsePage(r)))
	case http.MethodPost:
		var input struct {
			Title     string   `json:"title"`
			Edition   string   `json:"edition"`
			Summary   string   `json:"summary"`
			AgeRange  string   `json:"age_range"`
			Materials []string `json:"materials"`
			Steps     []string `json:"steps"`
			Editor    string   `json:"editor"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		result, err := s.flow.CreateReviewArchive(input.Title, input.Edition, input.Summary, input.AgeRange, input.Materials, input.Steps, input.Editor, "api-reviewer")
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusCreated, result)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) recordAction(w http.ResponseWriter, r *http.Request) {
	id := cleanPath(strings.TrimPrefix(r.URL.Path, "/records/"))
	parts := strings.Split(strings.Trim(id, "/"), "/")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch parts[1] {
	case "status":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var input struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		result, err := s.flow.AddStatus(parts[0], input.Content)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
