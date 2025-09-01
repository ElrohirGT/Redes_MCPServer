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

type NixHubSearchPkgsResponse struct {
	Query        string          `json:"query"`
	TotalResults int             `json:"total_results"`
	Results      []NixHubPkgInfo `json:"results"`
}

type NixHubPkgInfo struct {
	Name        string             `json:"name"`
	Summary     string             `json:"summary"`
	HomepageUrl string             `json:"homepage_url"`
	License     string             `json:"license"`
	Releases    []NixHubPkgRelease `json:"releases"`
}

type NixHubPkgRelease struct {
	Version          string    `json:"version"`
	LastUpdated      time.Time `json:"last_updated"`
	PlatformsSummary string    `json:"platforms_summary"`
	OutputsSummary   string    `json:"outputs_summary"`
}

func search_package_core(ctx context.Context, name string) (SearchPackageResult, error, bool) {
	result := SearchPackageResult{
		Packages: []NixHubPkgInfo{},
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

	log.Println("Received:\n", string(body))
	var nhResponse NixHubSearchPkgsResponse
	err = json.Unmarshal(body, &nhResponse)
	if err != nil {
		return result, errors.Join(INTERNAL_MCP_ERROR, err), true
	}

	log.Println("Getting new information for last 20 packages...")
	for _, pkg := range limit(nhResponse.Results, 20) {
		parseUrl, err = url.Parse("https://search.devbox.sh/v2/pkg")
		if err != nil {
			return result, err, true
		}
		queryValues = parseUrl.Query()
		queryValues.Add("name", pkg.Name)
		parseUrl.RawQuery = queryValues.Encode()

		finalUrl = parseUrl.String()
		log.Println("Making request to:", finalUrl)
		req, err = http.NewRequestWithContext(ctx, "GET", finalUrl, nil)
		if err != nil {
			return result, errors.Join(
				INTERNAL_MCP_ERROR,
				errors.New("Failed to construct request with context"),
				err,
			), true
		}
		req.Header.Add("Accept", "application/json")

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return result, errors.Join(EXTERNAL_ERROR, err), false
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return result, errors.Join(EXTERNAL_ERROR, err), false
		}

		log.Println("Received:\n", string(body))
		var nhResponse NixHubPkgInfo
		err = json.Unmarshal(body, &nhResponse)
		if err != nil {
			return result, errors.Join(INTERNAL_MCP_ERROR, err), true
		}

		result.Packages = append(result.Packages, nhResponse)
	}

	return result, nil, false
}

func limit[T any](arr []T, limit int) []T {
	if len(arr) > limit {
		return arr[:limit]
	}

	return arr
}
