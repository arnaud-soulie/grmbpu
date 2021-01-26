package main

import (
	"testing"
)

func TestCheckProjectExcluded(t *testing.T) {
	if !checkProjectExcluded("project2", []string{"project1", "project2", "project3"}) {
		t.Fatal("Project not correctly excluded")
	}
	if checkProjectExcluded("project4", []string{"project1", "project2", "project3"}) {
		t.Fatal("Project excluded whereas it shouldn't be")
	}

}
