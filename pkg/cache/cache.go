package cache

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dadjeibaah/production-tips-server/pkg/search"
	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type TipCache struct {
	suggestions []*search.Tip
	searchers []search.VideoSearch

}

func (h *TipCache) Search(q string) []*search.Tip{
	if len(h.suggestions) > 0 && h.suggestions != nil {
		return h.suggestions
	}

	for _, s := range h.searchers {
		tips, err := s.Search(q)
		log.Print(err)
		h.suggestions = append(h.suggestions, tips...)
	}
	return h.suggestions
}

func (h *TipCache) ClearSuggestionsOnTimeout() {
	for {
		<-time.After(1 * time.Hour)
		h.suggestions = nil
	}
}

func (h *TipCache) WithYoutubeSearcher(apiKey string) *TipCache{
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("unable to create Youtube Service: %s", err)
		return h
	}
	h.searchers = append(h.searchers, search.NewYTSearcher(service))
	return h
}

func (h *TipCache) WithVimeoSearcher(bearerToken string) *TipCache{
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: bearerToken,
	}))

	vmoClient := vimeo.NewClient(client, nil)
	h.searchers = append(h.searchers, search.NewVimeoSearcher(vmoClient, bearerToken))
	return h
}

func NewTipCache() (*TipCache, error) {
	tips := make([]*search.Tip, 0)
	return &TipCache{suggestions: tips}, nil
}
