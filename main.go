package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

func main() {
	insensitive := false
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Not enough argument")
		os.Exit(1)
	}

	if slices.Contains(args, "-i") || slices.Contains(args, "--i") {
		insensitive = true
		args = slices.DeleteFunc(args, func(s string) bool {
			return s == "-i" || s == "--i"
		})
	}
	word := args[0]
	filename := args[1]

	if word == "" {
		fmt.Println("Word cannot be empty string")
		os.Exit(1)
	}

	if insensitive {
		word = strings.ToLower(word)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if fi.IsDir() {
		fmt.Println("Cannot read a directory")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	writter := bufio.NewWriter(os.Stdout)

	lineNumber := 1

	for scanner.Scan() {
		text := scanner.Text()

		if insensitive {
			text = strings.ToLower(text)
		}

		if !strings.Contains(text, word) {
			lineNumber++
			continue
		}
		var formatedLine string

		if insensitive {
			formatedLine = caseInsensitiveColoring(scanner.Text(), word)
		} else {
			formatedLine = strings.ReplaceAll(scanner.Text(), word, fmt.Sprintf("\033[36m%s\033[0m", word))
		}

		writter.Write([]byte(formatedLine + "\n"))
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading file:", err)
		os.Exit(1)
	}

	writter.Flush()
	os.Exit(0)
}

func removeDuplicateStrings(s []string) []string {
	if len(s) <= 1 {
		return s
	}

	slices.Sort(s)
	j := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] {
			j++
			s[j] = s[i]
		}
	}
	return s[:j+1]
}

func caseInsensitiveColoring(text, s string) string {
	r, err := regexp.Compile("(?i)" + s)
	if err != nil {
		panic(err)
	}
	words := removeDuplicateStrings(r.FindAllString(text, -1))
	result := text
	for _, v := range words {
		result = strings.ReplaceAll(result, v, fmt.Sprintf("\033[36m%s\033[0m", v))
	}

	return result
}
