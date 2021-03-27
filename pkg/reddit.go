package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	REDDIT_USER_AGENT   = "file-utils (https://github.com/StefanWin/file-utils"
	REDDIT_IMAGE_DOMAIN = "https://i.redd.it/"
)

type RedditAPI struct {
	Kind string        `json:"kind"`
	Data RedditAPIData `json:"data"`
}

type RedditAPIData struct {
	ModHash string       `json:"modHash"`
	Count   int          `json:"dist"`
	Posts   []RedditPost `json:"children"`
	Before  string       `json:"before"`
	After   string       `json:"after"`
}

type RedditPost struct {
	Kind string         `json:"kind"`
	Data RedditPostData `json:"data"`
}

type RedditPostData struct {
	SubReddit      string `json:"subreddit"`
	IsGallery      bool   `json:"is_gallery"`
	Title          string `json:"title"`
	Name           string `json:"name"`
	UpVotes        int    `json:"ups"`
	FlairText      string `json:"link_flair_text"`
	URL            string `json:"url"`
	URLDestination string `json:"url_overridden_by_dest"`
	MediaMetaData  map[string]struct {
		ID string `json:"id"`
	} `json:"media_metadata"`
}

func gsri(subreddit string, after string) ([]string, string, error) {
	url := fmt.Sprintf("https://reddit.com/r/%s/top.json?sort=top&t=all&limit=100", subreddit)
	if after != "" {
		url = fmt.Sprintf("%s&after=%s", url, after)
	}
	log.Printf("fetching %s\n", url)
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("user-agent", REDDIT_USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	rAPI := RedditAPI{}
	if err := json.Unmarshal(data, &rAPI); err != nil {
		return nil, "", err
	}

	urls := make([]string, 0)

	log.Printf("received %d posts\n", len(rAPI.Data.Posts))

	for _, post := range rAPI.Data.Posts {
		if post.Data.IsGallery {
			for _, v := range post.Data.MediaMetaData {
				urls = append(urls, fmt.Sprintf("%s%s.jpg", REDDIT_IMAGE_DOMAIN, v.ID))
			}
		}
		if strings.HasSuffix(post.Data.URL, ".jpg") || strings.HasSuffix(post.Data.URL, ".png") || strings.HasSuffix(post.Data.URL, ".gif") {
			urls = append(urls, post.Data.URL)
		}
	}

	if len(rAPI.Data.Posts) == 0 {
		return urls, "", nil
	}
	return urls, rAPI.Data.Posts[len(rAPI.Data.Posts)-1].Data.Name, nil
}

func GetSubRedditImageURLS(subreddit string) ([]string, error) {

	urls, last, err := gsri(subreddit, "")
	if err != nil {
		return nil, err
	}

	xd := last

	for {
		time.Sleep(time.Second * 3)
		u, l, err := gsri(subreddit, xd)
		if err != nil {
			return nil, err
		}
		if len(u) == 0 {
			break
		}
		urls = append(urls, u...)
		xd = l
	}

	return urls, nil
}

func gsrif(subreddit string, after string, flair string) ([]string, string, error) {
	url := fmt.Sprintf("https://reddit.com/r/%s/search.json?q=flair:%s&sort=top&t=all&limit=100&restrict_sr=on&include_over_18=on", subreddit, flair)
	if after != "" {
		url = fmt.Sprintf("%s&after=%s", url, after)
	}
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("user-agent", REDDIT_USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	rAPI := RedditAPI{}
	if err := json.Unmarshal(data, &rAPI); err != nil {
		return nil, "", err
	}

	urls := make([]string, 0)
	for _, post := range rAPI.Data.Posts {
		if post.Data.FlairText == flair {
			if post.Data.IsGallery {
				for _, v := range post.Data.MediaMetaData {
					urls = append(urls, fmt.Sprintf("%s%s.jpg", REDDIT_IMAGE_DOMAIN, v.ID))
				}
			}
			if strings.HasSuffix(post.Data.URL, ".jpg") || strings.HasSuffix(post.Data.URL, ".png") || strings.HasSuffix(post.Data.URL, ".gif") {
				urls = append(urls, post.Data.URL)
			}
		}
	}

	if len(rAPI.Data.Posts) == 0 {
		return urls, "", nil
	}
	return urls, rAPI.Data.Posts[len(rAPI.Data.Posts)-1].Data.Name, nil
}

func GetSubRedditImageURLSFlair(subreddit string, flair string) ([]string, error) {
	urls, last, err := gsrif(subreddit, "", flair)
	if err != nil {
		return nil, err
	}
	xd := last
	for {
		time.Sleep(time.Second * 3)
		u, l, err := gsrif(subreddit, xd, flair)
		if err != nil {
			return nil, err
		}
		if len(u) == 0 {
			break
		}
		urls = append(urls, u...)
		xd = l
	}

	return urls, nil
}
