package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mycodesmells/golang-examples/nats/pubsub/proto"
)

// Data to be rendered into post page template.
type postPageData struct {
	Title   string
	Content string
	Date    string
}

// URL-friendly version of post title.
func (d postPageData) Slug() string {
	t := strings.TrimSpace(d.Title)
	t = strings.ToLower(t)
	return strings.Replace(t, " ", "-", -1)
}

type pageGenerator struct {
	basePath string
}

func newPageGenerator(basePath string) (pageGenerator, error) {
	if err := os.MkdirAll(basePath, os.ModeDir); err != nil {
		return pageGenerator{}, errors.Wrapf(err, "failed to create '%s' directory", basePath)
	}

	return pageGenerator{
		basePath: basePath,
	}, nil
}

func (g pageGenerator) Generate(pubMsg pb.PublishPostMessage) error {
	data := postPageData{
		Title:   pubMsg.Title,
		Content: pubMsg.Content,
		Date:    time.Now().Format("01-02-2006"),
	}
	fileName := filepath.Join(g.basePath, fmt.Sprintf("%s.html", data.Slug()))

	f, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to create post file")
	}

	t := template.New("post-page")
	buff, err := Asset("post-page.tpl")
	if err != nil {
		return errors.Wrap(err, "failed to load template")
	}
	t, err = t.Parse(string(buff))
	if err != nil {
		return errors.Wrap(err, "failed to parse template file")
	}

	if err := t.Execute(f, data); err != nil {
		if rErr := os.Remove(fileName); rErr != nil {
			log.Errorf("Failed to remove corrupt file: %v", rErr)
		}

		return errors.Wrap(err, "failed to render post page")
	}
	log.Infof("Generated new post page: '%s'", fileName)

	return nil
}
