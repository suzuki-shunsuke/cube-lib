package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v28/github"
	"github.com/suzuki-shunsuke/cube-lib/cube/download"
)

type (
	ClientParams struct {
		Token      string
		BaseURL    string
		UploadURL  string
		Enterprise bool
	}

	GetReleasesByTagClient interface {
		GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	}

	DownloadReleaseAssetClient interface {
		DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64) (rc io.ReadCloser, redirectURL string, err error)
	}

	GetReleaseByTagParams struct {
		Owner string
		Repo  string
		Tag   string
	}

	DonwnloadAssetParams struct {
		ID    int64
		Repo  string
		Owner string
	}
)

func NewClient(ctx context.Context, params ClientParams) (*github.Client, error) {
	httpClient := oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.Token}))
	if params.Enterprise {
		return github.NewEnterpriseClient(params.BaseURL, params.UploadURL, httpClient)
	}
	return github.NewClient(httpClient), nil
}

func GetReleaseByTag(
	ctx context.Context, client GetReleasesByTagClient, params GetReleaseByTagParams,
) (*github.RepositoryRelease, error) {
	if client == nil {
		return nil, errors.New("*github.Client is nil")
	}
	release, resp, err := client.GetReleaseByTag(ctx, params.Owner, params.Repo, params.Tag)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get GitHub release by tag %s %s %s: %w",
			params.Owner, params.Repo, params.Tag, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to get GitHub release by tag %s %s %s: status code %d >= 400",
			params.Owner, params.Repo, params.Tag, resp.StatusCode)
	}

	return release, nil
}

func DownloadAsset(
	ctx context.Context, client DownloadReleaseAssetClient, params DonwnloadAssetParams,
) (io.ReadCloser, error) {

	r, redirectURL, err := client.DownloadReleaseAsset(
		ctx, params.Owner, params.Repo, params.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to download GitHub release %s %s: %w", params.Owner, params.Repo, err)
	}
	if redirectURL != "" {
		return download.Download(ctx, &http.Client{}, redirectURL, download.Option{})
	}

	return r, nil
}
