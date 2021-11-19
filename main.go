package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

const (
	filePerm = 0600
	dirPerm  = 0700
)

func getClient(logger *logrus.Logger, config *oauth2.Config, tokenFile string) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(logger, config)
		saveToken(logger, tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func main() {
	var (
		outputFile string
		verbosity  bool
	)

	flag.StringVar(&outputFile, "output-file", "-", "Output file, default to stdout")
	flag.BoolVar(&verbosity, "v", false, "Verbose output")
	flag.Parse()

	logger := logrus.New()
	ctx := context.Background()

	if verbosity {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.ErrorLevel)
	}

	configRootDir, err := os.UserConfigDir()
	if err != nil {
		logger.Fatalf("Unable to determine UserConigDir: %v\n", err)
	}

	configDir := path.Join(configRootDir, "google-workspace-contacts")
	credentialsFile := path.Join(configDir, "credentials.json")
	tokenFile := path.Join(configDir, "token.json")

	os.MkdirAll(configDir, dirPerm)

	logger.Info("Reading credentials file...")
	credentials, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		logger.Fatalf("Unable to read client secret file: %v\n", err)
	}

	logger.Info("Generating configuration...")
	config, err := google.ConfigFromJSON(
		credentials,
		"https://www.googleapis.com/auth/directory.readonly",
		"https://www.googleapis.com/auth/contacts.readonly",
	)
	if err != nil {
		logger.Fatalf("Unable to parse client secret file to config: %v\n", err)
	}

	client := getClient(logger, config, tokenFile)

	logger.Info("Creating new People service client...")
	svc, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.Fatalf("Unable to retrieve People client: %v\n", err)
	}

	logger.Info("Requesting directory...")
	response, err := svc.People.ListDirectoryPeople().
		MergeSources("DIRECTORY_MERGE_SOURCE_TYPE_CONTACT").
		PageSize(1000).
		Sources("DIRECTORY_SOURCE_TYPE_DOMAIN_PROFILE", "DIRECTORY_SOURCE_TYPE_DOMAIN_CONTACT").
		ReadMask(strings.Join(queryFields(), ",")).
		Do()
	if err != nil {
		logger.Fatalf("Unable to retrieve directory list: %v\n", err)
	}

	contacts := []Contact{}

	logger.Infoln("Parsing response...")
	for _, item := range response.People {
		contact := Contact{}

		for _, email := range item.EmailAddresses {
			contact.Email = append(contact.Email, email.Value)
		}

		for _, name := range item.Names {
			// contact.Name = append(contact.Name, name.DisplayName)
			if contact.Name == "" || len(contact.Name) < len(name.DisplayName) {
				contact.Name = name.DisplayName
			}
		}

		for _, nickname := range item.Nicknames {
			// contact.Nickname = append(contact.Nickname, nickname.Value)
			if contact.Nickname == "" || len(contact.Nickname) < len(nickname.Value) {
				contact.Nickname = nickname.Value
			}
		}

		contacts = append(contacts, contact)
	}

	var output io.WriteCloser

	if outputFile == "-" {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
		if err != nil {
			logger.Fatalf("unable to open output file: %v\n", err)
		}
		defer output.Close()
	}

	logger.Info("Generating aliases list...")
	for _, contact := range contacts {
		for _, email := range contact.Email {
			fmt.Fprintf(
				output,
				"alias %s %s <%s>\n",
				contact.DisplayNickname(email),
				contact.DisplayName(email),
				email,
			)
		}
	}

	logger.Info("Done. I hope it's enough.")
}

type Contact struct {
	Name     string
	Nickname string
	Email    []string
}

func (c Contact) DisplayName(email string) string {
	if c.Name != "" {
		return c.Name
	}

	name := email[:strings.Index(email, "@")]
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.Title(name)

	return name
}

func (c Contact) DisplayNickname(email string) string {
	if c.Nickname != "" {
		return c.Nickname
	}

	return email[:strings.Index(email, "@")]
}

func queryFields() []string {
	return []string{
		"addresses",
		"ageRanges",
		"biographies",
		"birthdays",
		"calendarUrls",
		"clientData",
		"coverPhotos",
		"emailAddresses",
		"events",
		"externalIds",
		"genders",
		"imClients",
		"interests",
		"locales",
		"locations",
		"memberships",
		"metadata",
		"miscKeywords",
		"names",
		"nicknames",
		"occupations",
		"organizations",
		"phoneNumbers",
		"photos",
		"relations",
		"sipAddresses",
		"skills",
		"urls",
		"userDefined",
	}
}
