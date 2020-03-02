package govisual

import (
	"os"
	"path/filepath"
	"strings"
)

func unquote(s string) string {
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return ""
}

func pathjoin(s ...string) string {
	return filepath.Join(s...)
}

func sourceFilter(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}
