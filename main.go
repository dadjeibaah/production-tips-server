package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "ableton production tips", "Search term")
	maxResults = flag.Int64("max-results", 50, "Max YouTube results")
	addr       = flag.String("addr", ":8080", "address to run server on")
	apiKey     = flag.String("api-key", "", "YouTube API Key")
)

type Tip struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type TipServer struct {
	youtubeSvc  *youtube.Service
	suggestions []*Tip
}

func (h *TipServer) ClearSuggestionsOnTimeout(){
	for {
		<-time.After(1*time.Hour)
		h.suggestions = nil
	}
}

func (h *TipServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if h.suggestions == nil {
		h.suggestions = make([]*Tip, 0)
	}

	if len(h.suggestions) >= 50 {
		bytes, err := json.Marshal(h.suggestions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bytes)
	} else {
		rsp, err := h.youtubeSvc.
			Search.
			List("id,snippet").
			MaxResults(*maxResults).
			Q(*query).
			Do()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range rsp.Items {
			if item.Id.Kind == "youtube#video" {
				h.suggestions = append(h.suggestions, &Tip{
					Title: item.Snippet.Title,
					URL:   fmt.Sprintf("https://youtube.com/embed/%s?html5=1", item.Id.VideoId),
				})
			}
		}
	}
}

func NewTipServer(apiKey string) (*TipServer, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	newSuggestions := make([]*Tip, 0)

	return &TipServer{
		youtubeSvc:  service,
		suggestions: newSuggestions,
	}, nil
}

func main() {
	flag.Parse()
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Error("Unable to get ENV VAR 'API_KEY'")
		os.Exit(2)
	}

	server, err := NewTipServer(apiKey)
	if err != nil {
		log.Errorf("Failed to initialize server: %s", err)
	}

	s := &http.Server{
		Addr:      *addr,
		Handler:   server,
	}
	go server.ClearSuggestionsOnTimeout()

	log.Infof("Starting server on %s", *addr)

	s.ListenAndServe()

}
