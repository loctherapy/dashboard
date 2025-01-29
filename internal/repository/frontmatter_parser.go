package repository

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type FrontMatterParser struct {
	FrontMatterRE *regexp.Regexp
	ContextRE     *regexp.Regexp
}

func NewFrontMatterParser(frontMatterRE *regexp.Regexp, contextRE *regexp.Regexp) *FrontMatterParser {
	return &FrontMatterParser{
		FrontMatterRE: frontMatterRE,
		ContextRE: contextRE,
	}
}

func (f *FrontMatterParser) Parse(filePath string) (string, int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var context string
	var contextGravity int
	var gravity int
	inFrontMatter := false

	for scanner.Scan() {
		line := scanner.Text()
		if f.FrontMatterRE.MatchString(line) {
			if inFrontMatter {
				// End of front matter
				break
			} else {
				// Start of front matter
				inFrontMatter = true
				continue
			}
		}

		if inFrontMatter {
			if strings.HasPrefix(line, "context:") {
				context = strings.TrimSpace(strings.TrimPrefix(line, "context:"))
				// Extract context name and gravity using the regex
				if match := f.ContextRE.FindStringSubmatch(context); match != nil {
					for i, name := range f.ContextRE.SubexpNames() {
						if name == "context_name" && i < len(match) {
							context = match[i]
						} else if name == "context_gravity" && i < len(match) {
							contextGravity, _ = strconv.Atoi(match[i])
						}
					}
				}
			} else if strings.HasPrefix(line, "gravity:") {
				gravityStr := strings.TrimSpace(strings.TrimPrefix(line, "gravity:"))
				gravity, _ = strconv.Atoi(gravityStr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, 0, err
	}

	return context, contextGravity, gravity, nil
}