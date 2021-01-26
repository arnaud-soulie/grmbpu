package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func checkProjectExcluded(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

func main() {
	config := readConfigFile("config.json")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GithubToken})
	tc := oauth2.NewClient(ctx, ts)

	client, err := github.NewEnterpriseClient(config.Manifest.Server, config.Manifest.Server, tc)
	if err != nil {
		fmt.Println(err)
	}

	f, _, _, err := client.Repositories.GetContents(ctx, config.Manifest.Owner, config.Manifest.Repo, config.Manifest.Filename, nil)

	var decoded []byte
	if err == nil {
		fmt.Println("------------------")
		decoded, _ = base64.StdEncoding.DecodeString((*f.Content))
	}

	//Read the manifest file, map all the remotes
	m := Manifest{}
	readManifest(decoded, &m)

	remotes := listRemotes(&m)
	fmt.Println(remotes)

	//Process all projects
	for _, project := range m.Project {
		if !checkProjectExcluded(project.Name, config.ExcludeList) {
			fmt.Printf("Project: %v ; Remote label: %v ; Remote infos: %v\n", project.Name, project.Remote, remotes[project.Remote])
			//Build github URL
			u := remotes[project.Remote]

			//Set new configuration for each branch declared in configuration file
			for _, branch := range config.Branches {
				//Build the full protection request from json configuration
				request := github.ProtectionRequest{AllowForcePushes: &branch.AllowForcePushes, AllowDeletions: &branch.AllowDeletions,
					RequireLinearHistory: &branch.RequireLinearHistory, EnforceAdmins: branch.EnforceAdmins}
				if branch.RequiredStatusChecks != nil {
					request.RequiredStatusChecks = &github.RequiredStatusChecks{Strict: branch.RequiredStatusChecks.Strict, Contexts: branch.RequiredStatusChecks.Contexts}
				}
				if branch.RequiredPullRequestReviews != nil {
					request.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{DismissStaleReviews: branch.RequiredPullRequestReviews.DismissStaleReviews,
						RequireCodeOwnerReviews: branch.RequiredPullRequestReviews.RequireCodeOwnerReviews, RequiredApprovingReviewCount: branch.RequiredPullRequestReviews.RequiredApprovingReviewCount}
				}
				fmt.Printf("%v\n", request)
				//TODO: manage different client instances if multiple remotes used
				repo, resp, e := client.Repositories.UpdateBranchProtection(ctx, u.Path, project.Name, branch.Name, &request)
				fmt.Printf("Return after request exec: rep: %v, resp: %v, e: %v", repo, resp, e)
			}
		}

	}
}
