package main

import (
	"net/url"
	"reflect"
	"testing"
)

const manifest_text = `<?xml version="1.0" encoding="utf-8"?>
<!-- This is a comment
-->
<manifest>

	<remote fetch="ssh://user:pass@domain.com:1234/owner/" name="origin"/>
	<remote fetch="http://user:pass@domain.com:1234/owner/" name="ending"/>
  
  <default remote="origin"/>

  <project name="repo1" path="repo1dest" remote="origin" revision="1234567891234567891234567891234567891234" targetbranch="master">
    <copyfile src="file1" dest="file1"/>
    <copyfile src="file2" dest="file2"/>
  </project>

  <project name="repo2" path="repo2dest" remote="ending" revision="234567891234567891234567891234" targetbranch="master"/>
  
</manifest>`

func TestReadManifest(t *testing.T) {
	manifest := Manifest{}
	wanted := Manifest{}
	wanted.Remote = []remote{remote{Fetch: "ssh://user:pass@domain.com:1234/owner/", Name: "origin"},
		remote{Fetch: "http://user:pass@domain.com:1234/owner/", Name: "ending"}}
	wanted.Project = []project{project{Name: "repo1", Path: "repo1dest", Remote: "origin"},
		project{Name: "repo2", Path: "repo2dest", Remote: "ending"}}
	readManifest([]byte(manifest_text), &manifest)
	if !reflect.DeepEqual(manifest, wanted) {
		t.Fatalf("Error: difference between result %v and expected %v", manifest, wanted)
	}

}
func TestListRemotes(t *testing.T) {
	manifest := Manifest{}
	readManifest([]byte(manifest_text), &manifest)
	remotes := make(map[string](url.URL))
	wanted := make(map[string](url.URL))
	user := url.UserPassword("user", "pass")
	wanted["origin"] = url.URL{Scheme: "http", Host: "domain.com:80", Path: "owner", User: user}
	wanted["ending"] = url.URL{Scheme: "http", Host: "domain.com:1234", Path: "owner", User: user}
	remotes = listRemotes(&manifest)
	if !reflect.DeepEqual(remotes, wanted) {
		t.Fatalf("Error: difference between result %v and expected %v", remotes, wanted)
	}
}

func TestGetURL(t *testing.T) {
	remote := "ssh://user:password@domain.com:1234/path/"
	url := getURL(remote)
	if url.Host != "domain.com:80" || url.Path != "path" {
		t.Fatalf(`Domain: '%v' expected: 'domain.com:80' ; Path: '%v' expected 'path'`, url.Host, url.Path)
	}
	remote = "http://user:password@domain.com:1234/path/"
	url = getURL(remote)
	if url.Host != "domain.com:1234" || url.Path != "path" {
		t.Fatalf(`Domain: '%v' expected: 'domain.com:1234' ; Path: '%v' expected 'path'`, url.Host, url.Path)
	}
}
