package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/mod/semver"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

	ancestorContent := readContent(ancestorFilePath)
	currentContent := readContent(currentFilePath)
	otherContent := readContent(otherFilePath)

	ancestorVersion := getVersion(ancestorContent)
	currentVersion := getVersion(currentContent)
	otherVersion := getVersion(otherContent)

	maxVersion := maxVersion([]string{otherVersion, ancestorVersion, currentVersion})

	// first replace to eliminates version conflicts
	if currentVersion != "" && ancestorVersion != "" && otherVersion != "" && currentVersion != otherVersion && otherVersion != ancestorVersion {
		currentContent = setVersion(currentContent, otherVersion)
		writeContent(currentFilePath, currentContent)
	}

	output, err := exec.Command("git", "merge-file", "-L mine", "-L base", "-L theirs", "-p", currentFilePath, ancestorFilePath, otherFilePath).CombinedOutput()
	if err != nil {
		log.Fatal(string(output))
	}

	currentContent = string(output)

	if currentVersion != "" {
		branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
		if err != nil {
			log.Fatal(string(branch))
		}
		println("Merging version " + otherVersion + " into " + strings.ReplaceAll(string(branch), "\n", "") + ". Calculated version is: " + maxVersion)
		currentContent = setVersion(string(output), maxVersion)
	}
	writeContent(currentFilePath, currentContent)
}

func readContent(path string) string {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileContent)
}

func writeContent(path string, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0)
	if err != nil {
		log.Fatal(err)
	}
}

func getVersion(content string) string {
	result := gjson.Get(content, "version")
	if !result.Exists() {
		return ""
	}

	if !semver.IsValid("v" + result.String()) {
		log.Fatal("Wrong semantic version: " + result.String())
	}

	return result.String()
}

func setVersion(content string, version string) string {
	result, err := sjson.Set(content, "version", version)
	if err != nil {
		log.Fatal(err)
	}

	result, err = sjson.Set(result, "packages..version", version)
	if err != nil {
		log.Fatal(err)
	}

	return result
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
