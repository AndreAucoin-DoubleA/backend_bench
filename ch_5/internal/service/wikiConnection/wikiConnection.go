package wikiconnection

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
)

func ConnectToWikiStream(ctx context.Context, url string) (scanner *bufio.Scanner, resp *http.Response) {
	fmt.Println("Connecting to:", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, nil
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", "backend-bench-dev (youremail@example.com)")

	response, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil, nil
	}

	if response.StatusCode != 200 {
		fmt.Println("Bad status code:", response.Status)
		_ = response.Body.Close()
		return nil, nil
	}

	fmt.Println("Connected! Reading events...")

	return bufio.NewScanner(response.Body), response

}
