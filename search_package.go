package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var INTERNAL_MCP_ERROR = errors.New("Internal MCP Error")
var EXTERNAL_ERROR = errors.New("External Error")

type NixHubPkgSummary struct {
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	LastUpdated time.Time `json:"last_updated"`
}

type NixHubSearchPkgsResponse struct {
	Query        string             `json:"query"`
	TotalResults int                `json:"total_results"`
	Results      []NixHubPkgSummary `json:"results"`
}

func search_package_core(ctx context.Context, name string) (SearchPackageResult, error, bool) {
	result := SearchPackageResult{
		Packages: []NixHubPkgSummary{},
	}

	parseUrl, err := url.Parse("https://search.devbox.sh/v2/search")
	if err != nil {
		return result, err, true
	}
	queryValues := parseUrl.Query()
	queryValues.Add("q", name)
	parseUrl.RawQuery = queryValues.Encode()

	finalUrl := parseUrl.String()
	log.Println("Making request to:", finalUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", finalUrl, nil)
	if err != nil {
		return result, errors.Join(
			INTERNAL_MCP_ERROR,
			errors.New("Failed to construct request with context"),
			err,
		), true
	}
	req.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, errors.Join(EXTERNAL_ERROR, err), false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, errors.Join(EXTERNAL_ERROR, err), false
	}

	log.Println("Received:", string(body))
	var nhResponse NixHubSearchPkgsResponse
	err = json.Unmarshal(body, &nhResponse)
	if err != nil {
		return result, errors.Join(INTERNAL_MCP_ERROR, err), true
	}

	result.Packages = nhResponse.Results
	return result, nil, false
}
