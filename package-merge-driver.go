package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/c0b/go-ordered-json"
	"golang.org/x/mod/semver"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {

	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: ", os.Args[0], "%O", "%A", "%B")
		fmt.Println("\t%O - ancestor’s version of the conflicting file")
		fmt.Println("\t%A - current version of the conflicting file")
		fmt.Println("\t%B - other branch's version of the conflicting file")
		return
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

	maxVersion := semver.Max(otherVersion, semver.Max(ancestorVersion, currentVersion))[1:]

	if len(maxVersion) > 0 {
		replaceVersion(currentFilePath, maxVersion)
	}
}

func resolveVersion(filePath string) string {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var om = ordered.NewOrderedMap()
	json.Unmarshal(byteValue, om)

	version := "v" + om.Get("version").(string)

	if !semver.IsValid(version) {
		fmt.Println("Wrong semantic version: " + version)
	}

	return version
}

func replaceVersion(path string, version string) {

	read, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	versionRegexp := regexp.MustCompile(`"version"\s*:.*,`)

	result := versionRegexp.ReplaceAllString(string(read), fmt.Sprintf(`"version": "%s",`, version))

	err = ioutil.WriteFile(path, []byte(result), 0)
	if err != nil {
		panic(err)
	}
}
