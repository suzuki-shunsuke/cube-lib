package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type (
	Option struct {
		Headers map[string]string
	}
)

func Download(
	ctx context.Context, client *http.Client, uri string, option Option,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	for k, v := range option.Headers {
		header.Add(k, v)
	}
	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("status code = %d >= 400", resp.StatusCode)
	}
	return resp.Body, nil
}
