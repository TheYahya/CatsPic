package main

import (
	"catspic/infra/rest"
	"catspic/internal/twitter"
	"catspic/internal/unsplash"
	"catspic/usecase"
	"fmt"
	"os"

	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("PORT can't be empty")
	}

	unsplashAccessKey := os.Getenv("UNSPLASH_ACCESS_KEY")
	if unsplashAccessKey == "" {
		logger.Fatal("UNSPLASH_ACCESS_KEY can't be empty")
	}

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	if consumerKey == "" {
		logger.Fatal("TWITTER_CONSUMER_KEY can't be empty")
	}

	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET_KEY")
	if consumerSecret == "" {
		logger.Fatal("TWITTER_CONSUMER_SECRET_KEY can't be empty")
	}

	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	if accessToken == "" {
		logger.Fatal("TWITTER_ACCESS_TOKEN can't be empty")
	}

	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		logger.Fatal("TWITTER_ACCESS_TOKEN_SECRET can't be empty")
	}

	un := unsplash.New(unsplashAccessKey)
	tw := twitter.New(twitter.TwitterDI{
		ConsumerKey:       consumerKey,
		ConsumerSecretKey: consumerSecret,
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
	})

	tweetUsecase := usecase.NewTweet(usecase.TweetDI{
		Unsplash: un,
		Twitter:  tw,
	})

	tweetRest := rest.NewTweet(rest.TweetDI{
		Logger:  logger,
		Usecase: tweetUsecase,
	})

	http.Handle("/tweet/send", http.HandlerFunc(tweetRest.Tweet))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal(fmt.Errorf("ListenAndServe, err: %w", err).Error())
	}

	logger.Info("Server started")
}
