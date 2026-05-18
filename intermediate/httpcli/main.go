package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func sendRequest(method string, rawUrl string, headers []string, queries []string, jsonBody string) {
	if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
		rawUrl = "https://" + rawUrl
	}
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println("Invalid url :", err)
		os.Exit(1)
	}

	queryParams := parsedURL.Query()
	for _, q := range queries {
		parts := strings.SplitN(q, "=", 2)
		if len(parts) == 2 {
			queryParams.Set(parts[0], parts[1])
		}
	}
	parsedURL.RawQuery = queryParams.Encode()

	var bodyReader io.Reader
	if jsonBody != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(jsonBody), &js); err != nil {
			fmt.Println("Invalid json:", err)
			os.Exit(1)
		}
		bodyReader = strings.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, parsedURL.String(), bodyReader)
	if err != nil {
		fmt.Println("Error creating request :", err)
		os.Exit(1)
	}

	for _, h := range headers {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) == 2 {
			req.Header.Set(parts[0], parts[1])
		}
	}

	if jsonBody != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request :", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	app := &cli.App{
		Name:  "httpcli",
		Usage: "A simple command-line HTTP client",

		Action: func(c *cli.Context) error {
			rawUrl := c.Args().First()
			if rawUrl == "" {
				return cli.Exit("URL is required", 1)
			}
			sendRequest("GET", rawUrl, c.StringSlice("header"), c.StringSlice("query"), "")
			return nil
		},

		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "query",
				Usage: "Add request parameters",
			},
			&cli.StringSliceFlag{
				Name:  "header",
				Usage: "Add request headers ",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "GET request",
				Action: func(c *cli.Context) error {
					rawUrl := c.Args().First()
					if rawUrl == "" {
						return cli.Exit("URL is required", 1)
					}
					sendRequest("GET", rawUrl, c.StringSlice("header"), c.StringSlice("query"), "")
					return nil
				},
			},
			{
				Name:  "post",
				Usage: "POST request",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "json",
						Usage: "JSON request body",
					},
				},
				Action: func(c *cli.Context) error {
					rawUrl := c.Args().First()
					if rawUrl == "" {
						return cli.Exit("URL is required", 1)
					}
					sendRequest("POST", rawUrl, c.StringSlice("header"), c.StringSlice("query"), c.String("json"))
					return nil
				},
			},
			{
				Name:  "put",
				Usage: "PUT request",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "json",
						Usage: "JSON request body",
					},
				},
				Action: func(c *cli.Context) error {
					rawUrl := c.Args().First()
					if rawUrl == "" {
						return cli.Exit("URL is required", 1)
					}
					sendRequest("PUT", rawUrl, c.StringSlice("header"), c.StringSlice("query"), c.String("json"))
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "DELETE request",
				Action: func(c *cli.Context) error {
					rawUrl := c.Args().First()
					if rawUrl == "" {
						return cli.Exit("URL is required", 1)
					}
					sendRequest("DELETE", rawUrl, c.StringSlice("header"), c.StringSlice("query"), "")
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
