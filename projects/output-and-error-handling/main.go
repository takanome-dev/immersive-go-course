package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const Server = "http://localhost:8080"
var ErrRetry =  errors.New("should retry")

func request() ([]byte, error) {
	fmt.Printf("ℹ making the request... ℹ")
	resp, err := http.Get(Server)

	if err != nil {
		fmt.Printf("ℹ request failed... ℹ")
		return []byte{}, fmt.Errorf("error getting weather: %v\n\n", err.Error())
	}

	defer resp.Body.Close()

    switch resp.StatusCode {
    case http.StatusOK:

        fmt.Printf("ℹ success, reading body... ℹ")
        body, err := io.ReadAll(resp.Body)

        if err != nil {
            fmt.Printf("ℹ error when reading response body... ℹ")
            return []byte{}, fmt.Errorf("an error occurred while reading response body: %v\n\n", err.Error())
        }

        fmt.Printf("ℹ returning response... ℹ")
        return body, nil
    case http.StatusTooManyRequests:
       if err := handleRetry(resp.Header.Get("Retry-After")); err != nil {
            return []byte{}, fmt.Errorf("error handling 429: %w", err.Error())
        }
      return []byte{}, ErrRetry
    default:
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            body = []byte("<error reading body>")
        }
        return body, nil
    }
}

func parseDelay(retryHeader string) (time.Duration, error) {
    waitFor, err := strconv.Atoi(retryHeader)
    if err == nil {
        fmt.Printf("\nℹ retry header is a whole number... ℹ\n");
        return time.Duration(waitFor) / time.Nanosecond * time.Second, nil
    }

    waitUntil, err := http.ParseTime(retryHeader)
    if err == nil {
        fmt.Printf("\nℹ retry header is a timesptamp... ℹ\n")
        return time.Until(waitUntil), nil
    }

    fmt.Printf("\nℹ failed converting string to int... ℹ\n")
    return -1, fmt.Errorf("couldn't parse retry header, value was: %q", retryHeader)
}

func handleRetry(retryHeader string) error {
    fmt.Printf("\nℹ failed with retry... ℹ\n")
    duration, err := parseDelay(retryHeader)
    fmt.Printf("\nℹ num of seconds to retry after: %d ℹ\n", duration)

    if err != nil {
        return err
    }

    if duration > 1 * time.Second {
        fmt.Fprintf(os.Stderr, "\nℹ server receiving too many requests - waiting for %s before retrying ℹ\n", duration)
    }

    if duration > 5 * time.Second {
        return fmt.Errorf("\nserver receiving too many requests and this one is taking for ever, giving up\n")
    }

    fmt.Printf("\nℹ sleeping before retry... ℹ\n")
    time.Sleep(duration)
    // fmt.Printf("ℹ retying request... ℹ")
    // request()
    return nil
}

func main() {
	// result, err := request()

	// if err != nil {
    //		fmt.Fprintf(os.Stderr, err.Error())
	//	os.Exit(1)
	// }

	// fmt.Fprintf(os.Stdout, "\nresponse: %s\n\n", result)
	// os.Exit(0)

    result, err := request()
    for {
        if err != nil {
            if errors.Is(err, ErrRetry) {
                continue
            }

            fmt.Fprint(os.Stderr, "error retrieving weather: %v\n", err)
            os.Exit(1)
        } else {
            fmt.Fprintf(os.Stdout, string(result))
             break
        }
    }
}
