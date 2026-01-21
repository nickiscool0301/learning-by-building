package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var words = []string{"hello", "world", "foo", "bar", "distributed", "system", "map", "reduce", "golang", "worker", "coordinator", "task", "count", "word"}

func main() {
	dir := "test-data"
	os.MkdirAll(dir, 0755)

	numFiles := 20
	wordsPerFile := 1000

	for i := 0; i < numFiles; i++ {
		content := generateContent(wordsPerFile)
		filename := filepath.Join(dir, fmt.Sprintf("file-%03d.txt", i))
		os.WriteFile(filename, []byte(content), 0644)
		fmt.Printf("Created %s\n", filename)
	}
}

func generateContent(numWord int) string {
	var sb strings.Builder
	for i := 0; i < numWord; i++ {
		sb.WriteString(words[rand.Intn(len(words))])
		sb.WriteString(" ")
	}
	return sb.String()

}
