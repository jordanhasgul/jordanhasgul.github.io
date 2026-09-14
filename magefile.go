//go:build mage

package main

import (
	"errors"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Posts mg.Namespace

func (Posts) Create(title string) error {
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
