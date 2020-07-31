package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
	"gitlab.com/c0b/go-ordered-json"
	"golang.org/x/mod/semver"
	"io/ioutil"
	"os"
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
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var om = ordered.NewOrderedMap()
	err = json.Unmarshal(file, om)
	if err != nil {
		panic(err)
	}

	version := "v" + om.Get("version").(string)

	if !semver.IsValid(version) {
		panic("Wrong semantic version: " + version)
	}

	return version
}

func replaceVersion(path string, version string) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result, err := sjson.Set(string(fileContent), "version", version)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, []byte(result), 0)
	if err != nil {
		panic(err)
	}
}
