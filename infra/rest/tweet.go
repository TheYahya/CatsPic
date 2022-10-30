package rest

import (
	"catspic/usecase"
	"net/http"

	"go.uber.org/zap"
)

type Tweeter interface {
	Tweet(w http.ResponseWriter, r *http.Request)
}

type (
	TweetDI struct {
		Logger  *zap.Logger
		Usecase usecase.Tweet
	}

	tweet struct {
		logger *zap.Logger
		uc     usecase.Tweet
	}
)

func NewTweet(di TweetDI) Tweeter {
	return &tweet{
		logger: di.Logger,
		uc:     di.Usecase,
	}
}

func (t *tweet) Tweet(w http.ResponseWriter, r *http.Request) {
	err := t.uc.SendRandomPhoto()
	if err != nil {
		t.logger.Error(err.Error())
	}

	t.logger.Info("photo sent")
}
