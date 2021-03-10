package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	//"strings"
	"log"
	"net/url"
)

type remote struct {
	Fetch string `xml:"fetch,attr"`
	Name  string `xml:"name,attr"`
}

type project struct {
	Name   string `xml:"name,attr"`
	Path   string `xml:"path,attr"`
	Remote string `xml:"remote,attr"`
}

type Manifest struct {
	Remote  []remote  `xml:"remote"`
	Project []project `xml:"project"`
}

func readManifest(data []byte, manifest *Manifest) {
	err := xml.Unmarshal(data, &manifest)
	if err != nil {
		fmt.Println("Error during manifest parsing: ")
		fmt.Println(err)
	}
}

func listRemotes(manifest *Manifest) map[string](url.URL) {
	remotes := make(map[string](url.URL))
	for _, r := range manifest.Remote {
		remotes[r.Name] = getURL(r.Fetch)
	}
	return remotes
}

func getURL(remoteURL string) url.URL {
	//Convert an ssh URL if needed
	u, err := url.Parse(remoteURL)
	if err == nil {
		if u.Scheme == "ssh" {
			fmt.Println("Remote " + remoteURL + " is an SSH URL")
			u.Scheme = "http"
			//Remove port information is any in case of ssh url
			//and replace with global sshPort value
			u.Host = strings.Split(u.Host, ":")[0] + ":" + strconv.Itoa(sshPort)
		}
		//Remove leading '/' character
		u.Path = strings.Replace(u.Path, "/", "", -1)
	} else {
		fmt.Println("Error during remote url processing:")
		log.Fatal(err)
	}
	return *u
}
