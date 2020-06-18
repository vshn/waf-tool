package file

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var regexLastIdFinder = regexp.MustCompile("id:([0-9]+)")

type RuleFile struct {
	Path string
	File billy.File
}

func (f *RuleFile) Open(worktree *git.Worktree) (billy.File, error) {
	file, err := worktree.Filesystem.OpenFile(f.Path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", f.Path, err)
	}
	f.File = file
	return file, nil
}

func (f *RuleFile) Get() (string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, f.File)
	if err != nil {
		return "", fmt.Errorf("cannot get content from file %s: %w", f.File.Name(), err)
	}
	return buffer.String(), nil
}

func (f *RuleFile) Write(content string) error {
	_, err := f.File.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("cannot write to file %s: %w", f.File.Name(), err)
	}
	return nil
}

func (f *RuleFile) Close() error {
	if err := f.File.Close(); err != nil {
		return fmt.Errorf("cannot close file %s: %w", f.File.Name(), err)
	}
	return nil
}

func (f *RuleFile) GetLastId() (int, error) {
	content, err := f.Get()
	if err != nil {
		return 0, fmt.Errorf("cannot get content from file %s: %w", f.File.Name(), err)
	}
	ids := regexLastIdFinder.FindStringSubmatch(content)
	var id = 0
	for _, stringId := range ids {
		currentId, err := strconv.Atoi(stringId)
		if err != nil {
			return id, fmt.Errorf("could not convert string %s to integer: %w", stringId, err)
		}
		if id < currentId {
			id = currentId
		}
	}
	return id, nil
}
