//go:build mage

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Blog mg.Namespace

func (Blog) Serve(watch *bool) error {
	args := []string{"serve"}
	if watch != nil && *watch {
		args = append(args, "-w")
	}

	err := sh.RunV("hugo", args...)
	if err != nil {
		return fmt.Errorf("serving blog: %w", err)
	}

	return nil
}

func (Blog) CreatePost(title string) error {
	title = strings.Trim(title, " ")
	if title == "" {
		return errors.New("title is missing")
	}

	file := fmt.Sprintf("posts/%s.md", title)
	err := sh.RunV("hugo", "new", "content", file)
	if err != nil {
		return fmt.Errorf("creating post: %w", err)
	}

	return nil
}
