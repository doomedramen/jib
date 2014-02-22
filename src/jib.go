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
	baseURL         = "http://repository.sonatype.org/service/local/artifact/maven/redirect?r=central-proxy"

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile() []dep {
	abspath, _ := os.Getwd()
	filename := abspath + "/" + packageFileName

	filebyte, err := ioutil.ReadFile(filename)
	check(err)

	var deps []dep

	var depsMap map[string]interface{}
	json.Unmarshal(filebyte, &depsMap)

	for key, value := range depsMap {

		splitName := strings.Split(key, "#")

		str, _ := value.(string)
		deppy := dep{splitName[1], splitName[0], str}

		deps = append(deps, deppy)
	}

	return deps

}

func downloadFile(depURL string) io.ReadCloser {

	resp, err := http.Get(depURL)
	check(err)
	return resp.Body
}

func runChecks() bool {

	_, err := os.Stat(packageFileName)
	if err != nil {
		fmt.Println("There is not a package.json in the current directory")
		return true
	}

	_, err = os.Stat(libFolderName)
	if err != nil {
		fmt.Println("There is not a 'lib' folder in the current directory")
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
		return //exit
	}

	deps := readFile()

	for _, value := range deps {
		file := value.ArtifactId + "-" + value.Version + jar

		finalURL := baseURL + "&g=" + value.GroupId + "&a=" + value.ArtifactId + "&v=" + value.Version

		fmt.Println("Going to dowload", finalURL)

		fileDownload := downloadFile(finalURL)
		defer fileDownload.Close()
		moveFile(fileDownload, file)
	}

}
