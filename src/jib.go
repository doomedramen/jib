package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	libFolderName   = "lib"
	packageFileName = "package.json"
	baseURL         = "http://search.maven.org/remotecontent?filepath="

	pom     = ".pom"
	jar     = ".jar"
	jarDock = "-javadoc.jar"
	tests   = "-tests.jar"
	source  = "-sources.jar"
)

type dep struct {
	GroupId    string
	ArtifactId string
	Version    string
}

// type deps struct {
// 	dependencies []dep
// }

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile() []dep {
	abspath, _ := os.Getwd()
	filename := abspath + "/" + packageFileName
	// fmt.Println(filename)

	_, err := os.Stat(filename)
	check(err)

	filebyte, err := ioutil.ReadFile(filename)
	check(err)

	var deps []dep

	var depsMap map[string]interface{}
	json.Unmarshal(filebyte, &depsMap)
	// fmt.Printf("Results: %v\n", depsMap)

	for key, value := range depsMap {

		splitName := strings.Split(key, "#")

		str, _ := value.(string)
		deppy := dep{splitName[1], splitName[0], str}

		deps = append(deps, deppy)

		// fmt.Println(deps)
	}

	return deps

}

func downloadFile(depURL string) io.ReadCloser {

	resp, err := http.Get(depURL)
	check(err)
	return resp.Body
}

func runChecks() bool {
	_, err := os.Stat(libFolderName)
	if err != nil {
		return true
	}
	return false
}

func moveFile(file io.ReadCloser, name string) {

	out, err := os.Create(libFolderName + string(os.PathSeparator) + name)
	defer out.Close()

	_, err = io.Copy(out, file)
	check(err)
}

func main() {

	err := runChecks()
	if err {
		fmt.Println("There is not a 'lib' folder in the current directory")
		return //exit
	}

	deps := readFile()

	for _, value := range deps {
		groupIdURL := strings.Replace(value.GroupId, ".", "/", -1)
		file := value.ArtifactId + "-" + value.Version + jar
		finalURL := baseURL + groupIdURL + "/" + value.ArtifactId + "/" + value.Version + "/" + file

		fmt.Println("Going to dowload", finalURL)

		fileDownload := downloadFile(finalURL)
		defer fileDownload.Close()
		moveFile(fileDownload, file)
	}

}
