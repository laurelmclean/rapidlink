[![Go Report Card](https://goreportcard.com/badge/github.com/laurelmclean/rapidlink)](https://goreportcard.com/report/github.com/laurelmclean/rapidlink)

# RapidLink

### Description

RapidLink is a simple URL shortener service and QR code generator written in Go. It provides a minimalistic web interface for shortening URLs and includes QR code generation for shortened links.

### Key Features

1. **URL Shortening**: Rapidly generate short URLs 
2. **QR Code Generator**: Generate basic QR codes to complement every shortened URL
3. **Simple Interface**: Provide a straightforward user interface for generating short URLs
4. **URLs Persistence** Original URLs and their corresponding shortened URLs are saved to the urls.json file to ensure your shortened links remain persistent across server restarts

### Installation

To install RapidLink, you need to have Go installed on your system.

1. Clone the repository:

    ```
    git clone https://github.com/laurelmclean/rapidlink.git
    ```

2. Navigate to the project directory

3. Install dependencies:

    ```
    go mod download
    ```

4. Build the project:

    ```
    go build
    ```

### Usage

1. Start the RapidLink server:

    ```
    go run main.go
    ```

2. Access RapidLink in your web browser at `http://localhost:3000`.

3. Enter the URL you want to shorten in the provided form and click on the "Generate RapidLink" button.

4. Once the URL is shortened, you will be provided with a shortened link and QR code.

5. Original URLs and their corresponding shortened URLs are saved to the urls.json file.

### Running Tests

To run the table-driven tests, run the following command

      go test


To run the benchmark test, run the following command

      go test -bench=.

### Third-Party Packages Used

1. [github.com/go-chi/chi](https://pkg.go.dev/github.com/go-chi/chi)

**Description:** `chi` is a lightweight, idiomatic and composable router for building Go HTTP services. In my project, it is used to define routes for handling form submissions, URL shortening, redirection, and QR code generation.

2. [github.com/skip2/go-qrcode](https://pkg.go.dev/github.com/skip2/go-qrcode)

**Description:** `go-qrcode` is a QR Code encoder library for Go. It provides functionalities to generate QR codes from textual data. This package is utilized in the project to generate QR codes for the shortened URLs.

### Screenshots

![RapidLink Form](screenshot_form.png)
![RapidLink Result](screenshot_result.png)
