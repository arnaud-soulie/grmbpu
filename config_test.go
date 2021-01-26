package main

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

const configText = `{
	"GithubToken":"1234567890",
	"Manifest":{
	   "Server":"https://domain.com/",
	   "Owner":"foo",
	   "Repo":"manifestrepo",
	   "Filename":"default.xml"
	},
	"ExcludeList":[
	   "repo1",
	   "repo2"
	],
	"Branches":[
	   {
		  "Name": "branch1",
		  "AllowForcePushes":true,
		  "AllowDeletions":true,
		  "RequireLinearHistory":false,
		  "EnforceAdmins":false,
		  "RequiredStatusChecks":{
			 "Strict":true,
			 "Contexts":[
				"test",
				"test2"
			 ]
		  },
		  "RequiredPullRequestReviews":{
			 "DismissStaleReviews":true,
			 "RequireCodeOwnerReviews":false,
			 "RequiredApprovingReviewCount":42
		  }
	   },
	   {
		"Name": "branch2",
		"AllowForcePushes":false,
		"AllowDeletions":false,
		"RequireLinearHistory":true,
		"EnforceAdmins":true
	 }
	]
 }`

func TestReadConfigFile(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	file.Write([]byte(configText))
	config := readConfigFile(file.Name())
	wanted := configuration{}
	wanted.GithubToken = "1234567890"
	wanted.Manifest = struct {
		Server   string `json:"Server"`
		Owner    string `json:"Owner"`
		Repo     string `json:"Repo"`
		Filename string `json:"Filename"`
	}{Server: "https://domain.com/", Owner: "foo", Repo: "manifestrepo", Filename: "default.xml"}
	wanted.ExcludeList = []string{"repo1", "repo2"}
	requiredStatusChecks := struct {
		Strict   bool     `json:"Strict"`
		Contexts []string `json:"Contexts"`
	}{Strict: true, Contexts: []string{"test", "test2"}}
	requiredPullRequestReviews := struct {
		DismissStaleReviews          bool `json:"DismissStaleReviews"`
		RequireCodeOwnerReviews      bool `json:"RequireCodeOwnerReviews"`
		RequiredApprovingReviewCount int  `json:"RequiredApprovingReviewCount"`
	}{DismissStaleReviews: true, RequireCodeOwnerReviews: false, RequiredApprovingReviewCount: 42}
	wanted.Branches = []struct {
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
	}{{Name: "branch1", AllowForcePushes: true, AllowDeletions: true, RequireLinearHistory: false, EnforceAdmins: false,
		RequiredStatusChecks: &requiredStatusChecks, RequiredPullRequestReviews: &requiredPullRequestReviews},
		{Name: "branch2", AllowForcePushes: false, AllowDeletions: false, RequireLinearHistory: true, EnforceAdmins: true}}
	if !reflect.DeepEqual(*config, wanted) {
		t.Fatalf("Error: result %v different from expected %v\n", config, wanted)
	}
	file.Close()
	os.Remove(file.Name())
}
