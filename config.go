package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type configuration struct {
	GithubToken string `json:"GithubToken"`
	Manifest    struct {
		Server   string `json:"Server"`
		Owner    string `json:"Owner"`
		Repo     string `json:"Repo"`
		Filename string `json:"Filename"`
	} `json:"Manifest"`
	ExcludeList []string `json:"ExcludeList"`
	Branches    []struct {
		Name                 string `json:"Name"`
		AllowForcePushes     bool   `json:"AllowForcePushes"`
		AllowDeletions       bool   `json:"AllowDeletions"`
		RequireLinearHistory bool   `json:"RequireLinearHistory"`
		EnforceAdmins        bool   `json:"EnforceAdmins"`
		RequiredStatusChecks *struct {
			Strict   bool     `json:"Strict"`
			Contexts []string `json:"Contexts"`
		} `json:"RequiredStatusChecks"`
		RequiredPullRequestReviews *struct {
			DismissStaleReviews          bool `json:"DismissStaleReviews"`
			RequireCodeOwnerReviews      bool `json:"RequireCodeOwnerReviews"`
			RequiredApprovingReviewCount int  `json:"RequiredApprovingReviewCount"`
		} `json:"RequiredPullRequestReviews"`
	} `json:"Branches"`
}

func readConfigFile(path string) *configuration {
	fmt.Println("Reading configuration file " + path)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Configuration file not found at " + path + ".\nPlease check that the path is correct.\n")
	}
	decoder := json.NewDecoder(file)
	configuration := configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error during json configuration file decoding. Please check file format.")
	}
	return &configuration
}
