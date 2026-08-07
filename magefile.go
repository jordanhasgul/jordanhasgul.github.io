//go:build mage

package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	texttemplate "text/template"
	"time"
	"unicode"

	"github.com/magefile/mage/mg"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	goldmarkext "github.com/yuin/goldmark/extension"
	goldmarkparser "github.com/yuin/goldmark/parser"
)

//go:embed templates/*.tmpl
var templates embed.FS

type Site mg.Namespace

func (_ Site) NewPost(title string) error {
	date := time.Now().
		Format(time.DateOnly)

	newPostFileName := fmt.Sprintf("%s-%s.md", date, slugify(title))
	newPostFilePath := filepath.Join("posts", newPostFileName)
	_, err := os.Stat(newPostFilePath)
	if err == nil {
		return fmt.Errorf("post with title '%s' already exists", title)
	}

	const newPostMarkdownTemplate = "post.md.tmpl"
	err = renderTemplate(
		newPostMarkdownTemplate,
		newPostFilePath,
		struct {
			Title string
			Date  string
		}{
			Title: title,
			Date:  date,
		},
	)
	if err != nil {
		return fmt.Errorf("rendering new post: %w", err)
	}

	return nil
}

func slugify(s string) string {
	slugFields := strings.FieldsFunc(
		strings.ToLower(s),
		func(r rune) bool {
			return !isAlphaNumeric(r)
		},
	)

	slug := strings.Join(slugFields, "-")
	return slug
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

type Post struct {
	Slug string

	Title string
	Date  time.Time

	Body htmltemplate.HTML
}

func (_ Site) Build() error {
	markdownPostFilePathsGlob := filepath.Join("posts", "*.md")
	markdownPostFilePaths, err := filepath.Glob(markdownPostFilePathsGlob)
	if err != nil {
		return fmt.Errorf("finding markdown posts: %w", err)
	}

	markdownRenderer := goldmark.New(
		goldmark.WithExtensions(
			goldmarkext.GFM,
			goldmarkmeta.Meta,
		),
		goldmark.WithParserOptions(
			goldmarkparser.WithAutoHeadingID(),
		),
	)
	posts := make([]Post, 0, len(markdownPostFilePaths))
	for _, markdownPostFilePath := range markdownPostFilePaths {
		post, err := renderMarkdownPostAsHTML(markdownRenderer, markdownPostFilePath)
		if err != nil {
			return fmt.Errorf("rendering markdown post '%s' as html: %w", markdownPostFilePath, err)
		}

		posts = append(posts, post)
	}

	sort.Slice(
		posts,
		func(i, j int) bool {
			var (
				dateI = posts[i].Date
				dateJ = posts[j].Date
			)
			return dateI.After(dateJ)
		},
	)

	const indexHTMLTemplate = "index.html.tmpl"
	err = renderTemplate(indexHTMLTemplate, "index.html", posts)
	if err != nil {
		return fmt.Errorf("rendering index: %w", err)
	}

	const postHTMLTemplate = "post.html.tmpl"
	for _, post := range posts {
		postFilePath := filepath.Join("posts", post.Slug+".html")

		err = renderTemplate(postHTMLTemplate, postFilePath, post)
		if err != nil {
			return fmt.Errorf("rendering post %s: %w", postFilePath, err)
		}
	}

	return nil
}

func renderMarkdownPostAsHTML(markdownRenderer goldmark.Markdown, markdownFilePath string) (Post, error) {
	markdownFileContent, err := os.ReadFile(markdownFilePath)
	if err != nil {
		var zero Post
		return zero, fmt.Errorf("reading markdown file: %w", err)
	}

	var htmlFileContent bytes.Buffer
	ctx := goldmarkparser.NewContext()
	err = markdownRenderer.Convert(
		markdownFileContent,
		&htmlFileContent,
		goldmarkparser.WithContext(ctx),
	)
	if err != nil {
		var zero Post
		return zero, fmt.Errorf("converting markdown to html: %w", err)
	}

	metadata := goldmarkmeta.Get(ctx)
	title, _ := metadata["title"].(string)
	if title == "" {
		var zero Post
		return zero, errors.New("missing 'title' in metadata")
	}

	dateString, _ := metadata["date"].(string)
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		var zero Post
		return zero, fmt.Errorf("invalid 'date' in metadata: %w", err)
	}

	fileName := filepath.Base(markdownFilePath)
	slug := strings.TrimSuffix(fileName, ".md")
	return Post{
		Slug: slug,

		Title: title,
		Date:  date,

		Body: htmltemplate.HTML(htmlFileContent.String()),
	}, nil
}

func renderTemplate(templateName, filePath string, data any) error {
	tmpl, err := texttemplate.ParseFS(templates, "templates/"+templateName)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", filePath, err)
	}

	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
