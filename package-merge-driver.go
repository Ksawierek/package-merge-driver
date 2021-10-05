package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
	"gitlab.com/c0b/go-ordered-json"
	"golang.org/x/mod/semver"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: ", os.Args[0], "%O", "%A", "%B")
		fmt.Println("\t%O - ancestor’s version of the conflicting file")
		fmt.Println("\t%A - current version of the conflicting file")
		fmt.Println("\t%B - other branch's version of the conflicting file")
		os.Exit(1)
	}

	// This is the information we pass through in the driver config via
	// the placeholders `%O %A %B` where:
	// %O -> ancestor’s version of the conflicting file
	// %A -> current version of the conflicting file
	// %B -> other branch's version of the conflicting file
	ancestorFilePath := os.Args[1]
	currentFilePath := os.Args[2]
	otherFilePath := os.Args[3]

	ancestorVersion := resolveVersion(ancestorFilePath)
	currentVersion := resolveVersion(currentFilePath)
	otherVersion := resolveVersion(otherFilePath)

	maxVersion := maxVersion([]string{otherVersion, ancestorVersion, currentVersion})

	// first replace to eliminates version conflicts
	replaceVersion(currentFilePath, otherVersion)

	output, err := exec.Command("git", "merge-file", "-p", "-L mine", "-L base", "-L theirs", currentFilePath, ancestorFilePath, otherFilePath).CombinedOutput()
	if err != nil {
		log.Fatal(string(output))
	}

	// then replace to max_version
	replaceVersion(currentFilePath, maxVersion)

	err = ioutil.WriteFile(currentFilePath, output, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func resolveVersion(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var om = ordered.NewOrderedMap()
	err = json.Unmarshal(file, om)
	if err != nil {
		log.Fatal(err)
	}

	version := "v" + om.Get("version").(string)

	if !semver.IsValid(version) {
		log.Fatal("Wrong semantic version: " + version)
	}

	return version[1:]
}

func replaceVersion(path string, version string) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	result, err := sjson.Set(string(fileContent), "version", version)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, []byte(result), 0)
	if err != nil {
		log.Fatal(err)
	}
}

func maxVersion(versions []string) string {
	max := "0.0.0"
	for i := range versions {
		if semver.Compare("v"+max, "v"+versions[i]) < 0 {
			max = versions[i]
		}
	}
	return max
}
