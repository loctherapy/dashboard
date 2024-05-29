package repository

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/loctherapy/dashboard/internal/model"
)

type ToDoRepository struct {
	ToDoPattern   *regexp.Regexp
	FrontMatterRE *regexp.Regexp
	ContextRE     *regexp.Regexp
	FileFetcher   *FileFetcher
}

func NewToDoRepository(fileFetcher *FileFetcher) *ToDoRepository {
	contextRE := regexp.MustCompile(`^(?P<context_name>[a-zA-Z0-9_-]+)-(?P<context_gravity>\d+)$`)
	todoPattern := regexp.MustCompile(`^\s*- \[ \] `)
	frontMatterRE := regexp.MustCompile(`(?m)^---\s*$`)
	return &ToDoRepository{
		ToDoPattern:   todoPattern,
		FrontMatterRE: frontMatterRE,
		ContextRE:     contextRE,
		FileFetcher:   fileFetcher,
	}
}

func (r *ToDoRepository) GetAll() ([]model.FileToDos, error) {
	var results []model.FileToDos

	files, err := r.FileFetcher.Fetch()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		context, gravity, err := r.parseFrontMatter(file)
		if err != nil {
			return nil, err
		}

		todos, err := r.extractToDos(file)
		if err != nil {
			return nil, err
		}

		if len(todos) > 0 {
			results = append(results, model.FileToDos{
				FilePath: file,
				ToDos:    todos,
				Context:  context,
				Gravity:  gravity,
			})
		}
	}

	return results, nil
}

func (r *ToDoRepository) parseFrontMatter(filePath string) (string, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var context string
	var gravity int
	inFrontMatter := false

	for scanner.Scan() {
		line := scanner.Text()
		if r.FrontMatterRE.MatchString(line) {
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
			} else if strings.HasPrefix(line, "gravity:") {
				gravityStr := strings.TrimSpace(strings.TrimPrefix(line, "gravity:"))
				gravity, _ = strconv.Atoi(gravityStr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	return context, gravity, nil
}

func (r *ToDoRepository) extractToDos(filePath string) ([]model.ToDo, error) {
	var todos []model.ToDo

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if r.ToDoPattern.MatchString(line) {
			todos = append(todos, model.ToDo{Line: line})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
