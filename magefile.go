//go:build mage

package main

import (
	"errors"
	"fmt"

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
	if title == "" {
		return errors.New("title is missing")
	}

	fileName := fmt.Sprintf("%s.md", title)
	err := sh.RunV("hugo", "new", "content", "posts/"+fileName)
	if err != nil {
		return fmt.Errorf("creating post: %w", err)
	}

	return nil
}
