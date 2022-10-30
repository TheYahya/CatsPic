package unsplash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const photoURL string = "https://api.unsplash.com/search/photos"

type Unsplash struct {
	key string
	c   *http.Client
}

func New(key string) *Unsplash {
	return &Unsplash{
		key: key,
		c:   &http.Client{},
	}
}

func (u *Unsplash) RandImage(queries []string) ([]byte, error) {
	rand.Seed(time.Now().Unix())
	q := queries[rand.Intn(len(queries))]

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?per_page=30&query=%s", photoURL, q), nil)
	req.Header.Set("Authorization", "Client-ID "+u.key)

	resp, err := u.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	unspalshResp := unsplashResponse{}
	if err := json.Unmarshal(bodyBytes, &unspalshResp); err != nil {
		return nil, err
	}

	photo := unspalshResp.Results[rand.Intn(len(unspalshResp.Results))]

	response, err := http.Get(photo.URLs.Regular)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
