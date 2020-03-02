package govisual

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func guessModuleName(path string) string {
	if name := getModuleNameFromGoMod(path); name != "" {
		return name
	}

	if name := getModuleNameFromGOPATH(path); name != "" {
		return name
	}

	return ""
}

func getModuleNameFromGoMod(path string) string {
	data, err := ioutil.ReadFile(pathjoin(path, "go.mod"))
	if err != nil {
		return ""
	}

	r := regexp.MustCompile("module (.*?)\n")
	res := r.FindStringSubmatch(string(data))
	if len(res) > 1 {
		return res[1]
	}

	return ""
}

func getModuleNameFromGOPATH(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return ""
	}

	return abs[len(gopath)+5:]
}
