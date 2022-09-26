package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func main() {
	source := ""
	destination := ""
	release := ""
	flavor := ""

	cwd, err := os.Getwd()
	to := "/usr/local"
	if err != nil {
		panic(fmt.Errorf("Getwd Error: %w", err))
	}

	flag.StringVar(&source, "source", cwd, "the source directory of nvim.")
	flag.StringVar(&destination, "destination", to, "the source directory of nvim.")
	flag.StringVar(&release, "release", "nightly", "the nightly or stable release.")
	flag.StringVar(&flavor, "flavor", "macos", "the flavor to install (unix64, macos)")

	flag.Parse()
	baseSource := source
	source = path.Join(source, fmt.Sprintf("nvim-%s/bin/nvim", flavor))
	destination = path.Join(destination, "bin/nvim")

	template := "/neovim/neovim/releases/download/%s/nvim-%s.tar.gz"
	release = fmt.Sprintf(template, release, flavor)

	fmt.Printf("baseSource: %s, destination: %s, release: %s\n", baseSource, destination, release)

	if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("Remove Symlink Error: %w", err))
	}
	ctx := context.Background()
	cli := NewClient(
		"https://github.com",
		release,
	)
	sr, err := cli.Fetch(ctx)
	if err != nil {
		panic(fmt.Errorf("Fetch Error: %w", err))
	}
	if err := Untar(baseSource, sr); err != nil {
		panic(fmt.Errorf("Untar Error: %w", err))
	}

	if err := os.Symlink(source, destination); err != nil {
		panic(fmt.Errorf("Symlink Error: %w", err))
	}
}

// Client wraps an http.Client.
type Client struct {
	URL        string
	Path       string
	HTTPClient *http.Client
}

// NewClient returns a new client to fetch neovim tars.
func NewClient(rawurl, path string) *Client {
	return &Client{
		URL:        rawurl,
		Path:       path,
		HTTPClient: http.DefaultClient,
	}
}

// Fetch returns the neovim tar file, and wrapps it in a TarReader.
func (c *Client) Fetch(ctx context.Context) (*TarReader, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid client URL: %w", err)
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme:")
	} else if u.Host == "" {
		return nil, fmt.Errorf("URL host required")
	}

	*u = url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   c.Path,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("invalid response: code=%d", resp.StatusCode)
	}

	return &TarReader{
		rc: resp.Body,
	}, nil
}

// TarReader implements ReadCloser.
type TarReader struct {
	rc io.ReadCloser
}

// Close implements io.ReadCloser
func (r *TarReader) Close() (err error) {
	if e := r.rc.Close(); err == nil {
		err = e
	}
	return err
}

// Read implements io.Reader.
func (r TarReader) Read(p []byte) (n int, err error) {
	return r.rc.Read(p)
}

// Untar performs the untar operation with a ReadCloser and a destination.
func Untar(dst string, r io.ReadCloser) error {
	fmt.Println("Untar destination", dst)
	if _, err := os.Stat(dst); err != nil {
		return fmt.Errorf("Unable to tar files - %w", err)
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return fmt.Errorf("error creating directory: %w", err)
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("error copying file: %w", err)
			}
		}
	}
}
