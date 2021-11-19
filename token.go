package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func getTokenFromWeb(logger *logrus.Logger, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf(
		"Go to the following link in your browser then copy back the authorization code: \n%v\nToken: ",
		authURL,
	)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logger.Fatalf("Unable to read authorization code: %v\n", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		logger.Fatalf("Unable to retrieve token from web: %v\n", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(logger *logrus.Logger, path string, token *oauth2.Token) {
	logger.Infof("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
	defer f.Close()
	if err != nil {
		logger.Fatalf("Unable to cache OAuth token: %v\n", err)
	}
	json.NewEncoder(f).Encode(token)
}
