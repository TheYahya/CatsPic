package usecase

import (
	"catspic/internal/twitter"
	"catspic/internal/unsplash"
)

type Tweet interface {
	SendRandomPhoto() error
}

type (
	TweetDI struct {
		Unsplash *unsplash.Unsplash
		Twitter  *twitter.Twitter
	}

	tweet struct {
		u *unsplash.Unsplash
		t *twitter.Twitter
	}
)

var queries []string = []string{"cat", "cats", "kitty", "kitten", "kittens"}

func NewTweet(di TweetDI) Tweet {
	return &tweet{
		u: di.Unsplash,
		t: di.Twitter,
	}
}

func (t *tweet) SendRandomPhoto() error {
	photo, err := t.u.RandImage(queries)
	if err != nil {
		return err
	}

	t.t.TweetPhoto("#Cat #Cats #Kitty #Kitten", photo)

	return nil
}
