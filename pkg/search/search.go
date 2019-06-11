package search

import (
	"errors"
	"fmt"
	"log"

	vmo "github.com/silentsokolov/go-vimeo/vimeo"
	"google.golang.org/api/youtube/v3"
)

type Tip struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type youtubeSearcher struct{
	svc *youtube.Service
	maxResult int64
}
type vimeoSearcher struct{
	token string
	client *vmo.Client
}

type VideoSearch interface {
	Search(q string) ([]*Tip, error)
}

func NewYTSearcher(service *youtube.Service) *youtubeSearcher {
	return &youtubeSearcher{service, 25}
}

func NewVimeoSearcher(client *vmo.Client, bearerToken string) *vimeoSearcher {
	return &vimeoSearcher{token: bearerToken, client: client}

}

func (y *youtubeSearcher) Search(q string) ([]*Tip, error) {
	rsp, err := y.svc.
		Search.
		List("id,snippet").
		MaxResults(y.maxResult).
		Q(q).
		Do()
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New("Failed call to Google API")

	}

	var tips []*Tip
	for _, item := range rsp.Items {
		if item.Id.Kind == "youtube#video" {
			tips = append(tips, &Tip{
				Title: item.Snippet.Title,
				URL:   fmt.Sprintf("https://youtube.com/embed/%s?html5=1", item.Id.VideoId),
			})
		}
	}
	return tips, nil
}

func (v *vimeoSearcher) Search(q string) ([]*Tip, error){
	videos, _, err := v.client.Videos.List(vmo.OptQuery(q))
	if err != nil{
		log.Print(err.Error())
		return nil, errors.New("Failed call to Vimeo API")
	}
	var tips []*Tip
	for _, vid := range videos {
		tips = append(tips, &Tip{
			Title: vid.Name,
			URL: vid.Link,
		})
	}
	return tips, nil
}