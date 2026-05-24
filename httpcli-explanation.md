# httpcli Code Review Notes

This document explains the `httpcli` project in this repository. The goal is to help you explain what was implemented during a code review.

## 1. What `httpcli` Is

`httpcli` is a small command-line HTTP client written in Go.

It lets a user send HTTP requests from the terminal, similar to a very small version of tools like `curl` or HTTPie.

The implementation is in:

- `httpcli/main.go`
- `httpcli/go.mod`
- `httpcli/go.sum`

There is also a compiled executable at:

- `httpcli/httpcli`

## 2. Main Features Implemented

The CLI supports these HTTP methods:

- `GET`
- `POST`
- `PUT`
- `DELETE`

It also supports:

- Adding query parameters with `--query`
- Adding request headers with `--header`
- Sending a JSON request body for `POST` and `PUT` with `--json`
- Automatically adding `https://` when the user enters a URL without a protocol
- Validating JSON before sending a request body
- Printing the response body to the terminal

## 3. Dependency Used

The project uses this external package:

```go
github.com/urfave/cli/v2
```

This package is used to build the command-line interface. It handles:

- App name and usage text
- Commands like `get`, `post`, `put`, and `delete`
- Flags like `--query`, `--header`, and `--json`
- Argument parsing
- Help output
- CLI exit errors

The dependency is declared in `httpcli/go.mod`:

```go
require github.com/urfave/cli/v2 v2.27.7
```

The other packages in `go.mod` are indirect dependencies used by `urfave/cli`.

## 4. File Structure

```text
httpcli/
  go.mod
  go.sum
  main.go
  httpcli
```

### `main.go`

This is the whole application source code.

It contains:

- Imports
- The `sendRequest` function
- The `main` function
- CLI command definitions

### `go.mod`

This defines the Go module name and dependencies.

The module name is:

```go
module httpcli
```

### `go.sum`

This stores checksums for dependencies. It helps Go verify that downloaded dependencies match the expected versions.

### `httpcli`

This is a compiled binary executable. It is the built version of the CLI app.

## 5. High-Level Program Flow

The program flow is:

1. The user runs a command in the terminal.
2. `main()` creates a CLI app using `urfave/cli`.
3. The CLI app parses the command, flags, and URL argument.
4. The matching command action calls `sendRequest`.
5. `sendRequest` builds an HTTP request.
6. The request is sent using Go's `http.DefaultClient`.
7. The response body is read.
8. The response body is printed to the terminal.

Example:

```bash
./httpcli get https://example.com
```

That command calls:

```go
sendRequest("GET", rawUrl, headers, queries, "")
```

## 6. Important Function: `sendRequest`

The main logic is inside this function:

```go
func sendRequest(method string, rawUrl string, headers []string, queries []string, jsonBody string)
```

This function receives everything needed to send a request:

- `method`: HTTP method, such as `GET`, `POST`, `PUT`, or `DELETE`
- `rawUrl`: URL entered by the user
- `headers`: header values from the `--header` flag
- `queries`: query values from the `--query` flag
- `jsonBody`: JSON body from the `--json` flag

### 6.1 URL Protocol Handling

The first part checks whether the URL starts with `http://` or `https://`.

```go
if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
	rawUrl = "https://" + rawUrl
}
```

This means the user can type:

```bash
./httpcli get example.com
```

The program changes it internally to:

```text
https://example.com
```

This is useful because Go's HTTP client needs a URL with a protocol scheme.

### 6.2 URL Parsing

After adding the protocol if needed, the code parses the URL:

```go
parsedURL, err := url.Parse(rawUrl)
```

If parsing fails, the program prints an error and exits:

```go
fmt.Println("Invalid url :", err)
os.Exit(1)
```

This prevents the app from trying to send a request to a broken URL.

### 6.3 Query Parameters

The code reads existing query parameters from the URL:

```go
queryParams := parsedURL.Query()
```

Then it loops through the query flags:

```go
for _, q := range queries {
	parts := strings.SplitN(q, "=", 2)
	if len(parts) == 2 {
		queryParams.Set(parts[0], parts[1])
	}
}
```

The expected query format is:

```text
key=value
```

Example:

```bash
./httpcli --query page=1 --query limit=10 get https://api.example.com/users
```

This creates a final URL like:

```text
https://api.example.com/users?limit=10&page=1
```

Important details:

- `strings.SplitN(q, "=", 2)` splits only on the first `=`.
- `queryParams.Set(key, value)` adds or replaces the query value.
- If the query is missing `=`, it is ignored.
- Existing query parameters in the URL are preserved unless a flag uses the same key.
- `queryParams.Encode()` safely URL-encodes the final query string.

### 6.4 JSON Body Handling

The request body is only created when `jsonBody` is not empty:

```go
if jsonBody != "" {
	var js interface{}
	if err := json.Unmarshal([]byte(jsonBody), &js); err != nil {
		fmt.Println("Invalid json:", err)
		os.Exit(1)
	}
	bodyReader = strings.NewReader(jsonBody)
}
```

This does two things:

1. It validates that the body is valid JSON.
2. It converts the JSON string into an `io.Reader` so it can be sent as the request body.

The parsed variable `js` is only used for validation. The actual request sends the original JSON string.

Example:

```bash
./httpcli post https://api.example.com/users --json '{"name":"Alice"}'
```

If the JSON is invalid, the request is not sent.

Example invalid JSON:

```bash
./httpcli post https://api.example.com/users --json '{"name":}'
```

The app prints:

```text
Invalid json: ...
```

Then it exits with status code `1`.

### 6.5 Creating the HTTP Request

The request is created here:

```go
req, err := http.NewRequest(method, parsedURL.String(), bodyReader)
```

This creates a request using:

- The HTTP method
- The final URL after query parameters were added
- The optional body reader

If request creation fails, the app prints an error and exits.

### 6.6 Headers

Headers are added from the `--header` flag:

```go
for _, h := range headers {
	parts := strings.SplitN(h, "=", 2)
	if len(parts) == 2 {
		req.Header.Set(parts[0], parts[1])
	}
}
```

The expected header format is:

```text
Header-Name=value
```

Example:

```bash
./httpcli --header Authorization=Bearer-token get https://api.example.com/profile
```

Important details:

- The code uses `SplitN` so values can contain `=`.
- `req.Header.Set` adds or replaces the header.
- If the header is missing `=`, it is ignored.

### 6.7 Automatic Content-Type for JSON

If a JSON body is provided, the code automatically sets:

```go
Content-Type: application/json
```

That happens here:

```go
if jsonBody != "" {
	req.Header.Set("Content-Type", "application/json")
}
```

This is important because APIs usually need this header to understand that the body is JSON.

One detail to know for review:

- If the user also passes a custom `Content-Type` header, this code overwrites it with `application/json` whenever `--json` is used.

### 6.8 Sending the Request

The request is sent here:

```go
resp, err := http.DefaultClient.Do(req)
```

`http.DefaultClient` is Go's built-in default HTTP client.

If the request fails because of network problems, DNS problems, or an invalid destination, the app prints:

```text
Error sending request : ...
```

Then it exits.

### 6.9 Closing the Response Body

The code closes the response body after reading:

```go
defer resp.Body.Close()
```

This is important because response bodies use resources. Closing them prevents resource leaks.

### 6.10 Reading and Printing the Response

The response body is read here:

```go
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

The app prints only the response body.

It does not currently print:

- HTTP status code
- Response headers
- Request URL
- Timing information

## 7. The `main` Function

The `main` function builds the CLI app:

```go
app := &cli.App{
	Name:  "httpcli",
	Usage: "A simple command-line HTTP client",
	...
}
```

Then it runs the app:

```go
app.Run(os.Args)
```

`os.Args` contains the command-line arguments typed by the user.

## 8. Global Flags

The app defines two global flags:

```go
--query
--header
```

They are defined in the app-level `Flags` section:

```go
Flags: []cli.Flag{
	&cli.StringSliceFlag{
		Name:  "query",
		Usage: "Add request parameters",
	},
	&cli.StringSliceFlag{
		Name:  "header",
		Usage: "Add request headers ",
	},
}
```

Because these are global flags, they should be placed before the command.

Correct:

```bash
./httpcli --query page=1 --header Accept=application/json get https://api.example.com/users
```

Incorrect:

```bash
./httpcli get --query page=1 https://api.example.com/users
```

The incorrect version fails because the `get` command itself does not define a `--query` command flag.

## 9. Commands

### 9.1 Default Action

If the user runs the app with only a URL and no command, it sends a `GET` request.

Example:

```bash
./httpcli https://example.com
```

This calls:

```go
sendRequest("GET", rawUrl, c.StringSlice("header"), c.StringSlice("query"), "")
```

If the URL is missing, it returns:

```text
URL is required
```

### 9.2 `get`

The `get` command sends a `GET` request.

Example:

```bash
./httpcli get https://example.com
```

It does not send a request body.

### 9.3 `post`

The `post` command sends a `POST` request.

Example:

```bash
./httpcli post https://api.example.com/users --json '{"name":"Alice"}'
```

It supports the command flag:

```bash
--json
```

The JSON is validated before the request is sent.

### 9.4 `put`

The `put` command sends a `PUT` request.

Example:

```bash
./httpcli put https://api.example.com/users/1 --json '{"name":"Bob"}'
```

It also supports:

```bash
--json
```

### 9.5 `delete`

The `delete` command sends a `DELETE` request.

Example:

```bash
./httpcli delete https://api.example.com/users/1
```

It does not send a request body.

## 10. Example Usage

### Basic GET

```bash
cd httpcli
./httpcli get https://example.com
```

### GET Without Protocol

```bash
./httpcli get example.com
```

The app changes this to:

```text
https://example.com
```

### GET With Query Parameters

```bash
./httpcli --query page=1 --query limit=10 get https://api.example.com/users
```

### GET With Headers

```bash
./httpcli --header Accept=application/json get https://api.example.com/users
```

### POST With JSON

```bash
./httpcli post https://api.example.com/users --json '{"name":"Alice","role":"admin"}'
```

### PUT With JSON

```bash
./httpcli put https://api.example.com/users/1 --json '{"name":"Alice Updated"}'
```

### DELETE

```bash
./httpcli delete https://api.example.com/users/1
```

## 11. Help Output

Running:

```bash
./httpcli --help
```

Shows:

```text
NAME:
   httpcli - A simple command-line HTTP client

USAGE:
   httpcli [global options] command [command options]

COMMANDS:
   get      GET request
   post     POST request
   put      PUT request
   delete   DELETE request
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --query value [ --query value ]    Add request parameters
   --header value [ --header value ]  Add request headers
   --help, -h                         show help
```

## 12. Error Handling

The app handles several error cases.

### Missing URL

If the user does not provide a URL:

```bash
./httpcli get
```

The app returns:

```text
URL is required
```

### Invalid URL

If `url.Parse` fails, the app prints:

```text
Invalid url : ...
```

Then it exits with code `1`.

### Invalid JSON

If the user passes invalid JSON to `--json`, the app prints:

```text
Invalid json: ...
```

Then it exits with code `1`.

### Request Creation Error

If Go cannot create the request, the app prints:

```text
Error creating request : ...
```

Then it exits.

### Request Sending Error

If the actual network request fails, the app prints:

```text
Error sending request : ...
```

Then it exits.

## 13. What Was Verified

I checked that the module builds/tests with:

```bash
go test ./...
```

Result:

```text
?   	httpcli	[no test files]
```

This means:

- The code compiles.
- There are no automated test files in this module yet.

I also checked the executable help output:

```bash
./httpcli --help
```

The CLI correctly shows the app name, usage, commands, and global flags.

## 14. Important Review Talking Points

These are good points to explain to your instructor.

### Why use `urfave/cli`?

It avoids manually parsing `os.Args`. It gives a cleaner way to define commands and flags.

Instead of writing custom argument parsing logic, the code defines:

- App metadata
- Global flags
- Commands
- Command-specific flags
- Action functions

### Why use `net/url`?

`net/url` safely parses and updates URLs.

This is better than manually joining query strings because it handles encoding.

For example, a space in a query value can be encoded safely.

### Why validate JSON?

The app checks JSON before sending it so the user gets immediate feedback when the body is invalid.

Without this check, the app might send bad JSON and the server would reject it later.

### Why use `io.Reader` for the body?

`http.NewRequest` expects the body as an `io.Reader`.

The code uses:

```go
strings.NewReader(jsonBody)
```

That turns the JSON string into a readable stream.

### Why close the response body?

HTTP response bodies must be closed after use.

The code uses:

```go
defer resp.Body.Close()
```

This makes sure resources are released after the function finishes.

## 15. Current Limitations

These are not necessarily wrong for a training project, but they are useful to mention honestly.

### 15.1 No Automated Tests Yet

The module has no test files.

A good improvement would be adding tests for:

- URL protocol auto-prefixing
- Query parameter handling
- Header parsing
- JSON validation
- Request creation

### 15.2 `sendRequest` Exits Directly

`sendRequest` calls:

```go
os.Exit(1)
```

This works for a CLI, but it makes the function harder to unit test.

A more testable design would return an error instead:

```go
func sendRequest(...) error
```

Then `main` could decide whether to exit.

### 15.3 Response Status Is Not Printed

The app only prints the response body.

It does not show:

- `200 OK`
- `404 Not Found`
- `500 Internal Server Error`

For debugging APIs, printing the status code would be helpful.

### 15.4 Response Read Error Is Ignored

This line ignores a possible error:

```go
body, _ := io.ReadAll(resp.Body)
```

A stronger version would check the error:

```go
body, err := io.ReadAll(resp.Body)
if err != nil {
	return err
}
```

### 15.5 Header and Query Format Errors Are Ignored

If the user passes:

```bash
--query page
```

or:

```bash
--header Authorization
```

The code silently ignores it because there is no `=`.

A nicer user experience would print a validation error.

### 15.6 Global Flags Must Be Before Commands

Because `--query` and `--header` are global flags, this works:

```bash
./httpcli --query page=1 get https://api.example.com/users
```

But this does not:

```bash
./httpcli get --query page=1 https://api.example.com/users
```

One possible improvement would be adding `--query` and `--header` as command flags for each command, or using CLI configuration that allows global flags after commands if desired.

### 15.7 No Timeout

The app uses:

```go
http.DefaultClient
```

The default client has no custom timeout.

For a production CLI, it is safer to create a client with a timeout:

```go
client := &http.Client{
	Timeout: 10 * time.Second,
}
```

### 15.8 JSON Content-Type Overwrites User Header

The code applies user headers first, then sets:

```go
Content-Type: application/json
```

So if the user passed a different `Content-Type`, it gets replaced whenever `--json` is used.

This is acceptable for a simple JSON-focused flag, but it is worth knowing.

## 16. Possible Instructor Questions and Answers

### What did you implement?

I implemented a command-line HTTP client in Go. It supports GET, POST, PUT, and DELETE requests. It can add query parameters, set headers, validate JSON bodies, send requests using Go's `net/http` package, and print the response body.

### Why did you choose this structure?

The CLI parsing is separated from the request-sending logic. The `main` function defines commands and flags, while `sendRequest` handles URL parsing, query parameters, headers, JSON body creation, request creation, sending the request, and printing the response.

### What is the most important function?

`sendRequest` is the most important function because it contains the HTTP request logic.

### How does the app know which HTTP method to use?

Each command passes a different method string to `sendRequest`.

Examples:

```go
sendRequest("GET", ...)
sendRequest("POST", ...)
sendRequest("PUT", ...)
sendRequest("DELETE", ...)
```

### How are query parameters added?

The app reads `--query key=value` flags, splits each value into key and value, stores them in `url.Values`, then writes the encoded query string back to the parsed URL.

### How are headers added?

The app reads `--header Header-Name=value` flags, splits each value into name and value, then adds them to `req.Header`.

### How does the JSON body work?

For `POST` and `PUT`, the user can pass `--json`. The app validates the JSON with `json.Unmarshal`. If valid, it sends the original JSON string as the request body and sets `Content-Type` to `application/json`.

### What would you improve next?

I would improve testability by making `sendRequest` return errors instead of calling `os.Exit`. Then I would add unit tests and print the HTTP status code along with the response body.

## 17. Simple Summary

`httpcli` is a Go CLI app that sends HTTP requests from the terminal.

The main work is:

- Parse CLI commands and flags with `urfave/cli`
- Prepare the URL and query parameters
- Validate optional JSON body
- Create an HTTP request
- Add headers
- Send the request
- Print the response body

The project is functional and compiles, but it can be improved with tests, better error handling, response status output, and request timeouts.
