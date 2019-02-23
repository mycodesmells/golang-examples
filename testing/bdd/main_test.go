package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
)

type apiClient struct{}

func (cli *apiClient) Create(quote, author string) (string, error) {
	url := "http://localhost:8888/quotes"

	payloadStr := fmt.Sprintf(`{"quote": "%s", "author": "%s"}`, quote, author)
	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resp map[string]string
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}
	return resp["id"], nil
}

func (cli apiClient) Get(id string) (quote, error) {
	url := fmt.Sprintf("http://localhost:8888/quotes/%s", id)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var q quote
	err := json.NewDecoder(res.Body).Decode(&q)
	return q, err
}

type testContext struct {
	id    string
	quote quote
	cli   *apiClient
}

func (ctx *testContext) createQuote(quote, author string) error {
	id, err := ctx.cli.Create(quote, author)
	if err != nil {
		return err
	}
	ctx.id = id
	return nil
}

func (ctx *testContext) askForLastCreatedQuote() error {
	q, err := ctx.cli.Get(ctx.id)
	if err != nil {
		return err
	}
	ctx.quote = q
	return nil
}

func (ctx *testContext) shouldGetAsQuote(quote string) error {
	if want, got := quote, ctx.quote.Quote; want != got {
		return fmt.Errorf("expected quote '%s', got '%s'", want, got)
	}
	return nil
}

func (ctx *testContext) shouldGetAsAuthor(author string) error {
	if want, got := author, ctx.quote.Author; want != got {
		return fmt.Errorf("expected author '%s', got '%s'", want, got)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	ctx := &testContext{
		cli: &apiClient{},
	}

	s.Step(`^I create quote "([^"]*)" by "([^"]*)"$`, ctx.createQuote)
	s.Step(`^I ask for last created quote$`, ctx.askForLastCreatedQuote)
	s.Step(`^I should get "([^"]*)" as quote$`, ctx.shouldGetAsQuote)
	s.Step(`^I should get "([^"]*)" as author$`, ctx.shouldGetAsAuthor)
}
