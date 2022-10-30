package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/dghubble/oauth1"
)

type (
	TwitterDI struct {
		ConsumerKey       string
		ConsumerSecretKey string
		AccessToken       string
		AccessTokenSecret string
	}

	Twitter struct {
		c *http.Client
	}
)

const StatusUpdate string = "https://api.twitter.com/1.1/statuses/update.json"
const MediaUpload string = "https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image"

type MediaInitResponse struct {
	MediaId          uint64 `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

func New(di TwitterDI) *Twitter {
	config := oauth1.NewConfig(di.ConsumerKey, di.ConsumerSecretKey)
	token := oauth1.NewToken(di.AccessToken, di.AccessTokenSecret)
	client := config.Client(context.Background(), token)
	return &Twitter{
		c: client,
	}
}

func (t *Twitter) TweetPhoto(caption string, photo []byte) error {
	mediaInitResponse, err := t.mediaInit(photo)
	if err != nil {
		return fmt.Errorf("can not init media, err: %w", err)
	}
	mediaId := mediaInitResponse.MediaId

	if t.MediaAppend(mediaId, photo) != nil {
		return fmt.Errorf("can not append media, err: %w", err)
	}

	if t.MediaFinilize(mediaId) != nil {
		return fmt.Errorf("can not fin media, err: %w", err)
	}

	if t.UpdateStatusWithMedia(caption, mediaId) != nil {
		return fmt.Errorf("can not update status, err: %w", err)
	}

	return nil
}

func (t *Twitter) mediaInit(media []byte) (*MediaInitResponse, error) {
	form := url.Values{}
	form.Add("command", "INIT")
	form.Add("media_type", "video/mp4")
	form.Add("total_bytes", fmt.Sprint(len(media)))
	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var mediaInitResponse MediaInitResponse
	err = json.Unmarshal(body, &mediaInitResponse)

	if err != nil {
		return nil, err
	}

	return &mediaInitResponse, nil
}

func (t *Twitter) MediaAppend(mediaId uint64, media []byte) error {
	step := 500 * 1024
	for s := 0; s*step < len(media); s++ {
		var body bytes.Buffer
		rangeBegining := s * step
		rangeEnd := (s + 1) * step
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}

		w := multipart.NewWriter(&body)

		w.WriteField("command", "APPEND")
		w.WriteField("media_id", fmt.Sprint(mediaId))
		w.WriteField("segment_index", fmt.Sprint(s))

		fw, err := w.CreateFormFile("media", "example.mp4")

		fw.Write(media[rangeBegining:rangeEnd])

		w.Close()

		req, err := http.NewRequest("POST", MediaUpload, &body)

		req.Header.Add("Content-Type", w.FormDataContentType())

		res, err := t.c.Do(req)
		if err != nil {
			return err
		}

		ioutil.ReadAll(res.Body)
	}

	return nil
}

func (t *Twitter) MediaFinilize(mediaId uint64) error {
	form := url.Values{}
	form.Add("command", "FINALIZE")
	form.Add("media_id", fmt.Sprint(mediaId))

	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := t.c.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(res.Body)
	return err
}

func (t *Twitter) UpdateStatusWithMedia(text string, mediaId uint64) error {
	form := url.Values{}
	form.Add("status", text)
	form.Add("media_ids", fmt.Sprint(mediaId))

	req, err := http.NewRequest("POST", StatusUpdate, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := t.c.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(res.Body)
	return err
}
