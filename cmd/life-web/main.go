package main

import (
	"conways/input"
	"conways/life"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// The frontend and its initial example grid are kept with this executable so
// the viewer can be run as a single local binary.
//
//go:embed static/index.html static/app.js static/style.css sample.txt
var assets embed.FS

type parseResponse struct {
	Grid   life.Grid `json:"grid"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
}

type stepRequest struct {
	Grid life.Grid `json:"grid"`
	Wrap bool      `json:"wrap"`
}

type gridResponse struct {
	Grid life.Grid `json:"grid"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "HTTP listen address")
	flag.Parse()

	server := &http.Server{Addr: *address, Handler: newHandler()}
	fmt.Printf("Life viewer listening at http://%s/\n", *address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", serveAsset("static/index.html", "text/html; charset=utf-8"))
	mux.HandleFunc("GET /app.js", serveAsset("static/app.js", "text/javascript; charset=utf-8"))
	mux.HandleFunc("GET /style.css", serveAsset("static/style.css", "text/css; charset=utf-8"))
	mux.HandleFunc("GET /sample.txt", serveAsset("sample.txt", "text/plain; charset=utf-8"))
	mux.HandleFunc("POST /api/parse", handleParse)
	mux.HandleFunc("POST /api/step", handleStep)
	return mux
}

func serveAsset(name, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.ReadFile(name)
		if err != nil {
			http.Error(w, "asset unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(data)
	}
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	grid, err := input.ParseGrid(io.LimitReader(r.Body, 10<<20), "uploaded grid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, parseResponse{Grid: grid, Width: len(grid[0]), Height: len(grid)})
}

func handleStep(w http.ResponseWriter, r *http.Request) {
	var request stepRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 10<<20))
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			writeError(w, http.StatusBadRequest, errors.New("request must contain one JSON object"))
		} else {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		}
		return
	}
	if err := validateGrid(request.Grid); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	next, err := life.Step(request.Grid, request.Wrap)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid grid: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, gridResponse{Grid: next})
}

func validateGrid(grid life.Grid) error {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return errors.New("invalid grid: grid must not be empty")
	}
	width := len(grid[0])
	for row, cells := range grid {
		if len(cells) != width {
			return fmt.Errorf("invalid grid: row %d has width %d; expected %d", row+1, len(cells), width)
		}
		for column, cell := range cells {
			if cell != 0 && cell != 1 {
				return fmt.Errorf("invalid grid: cell at row %d, column %d must be 0 or 1", row+1, column+1)
			}
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	message := strings.TrimSpace(err.Error())
	writeJSON(w, status, errorResponse{Error: message})
}
