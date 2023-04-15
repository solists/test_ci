package tgservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/solists/test_ci/pkg/logger"
)

type Downloader struct {
	token string

	server      string
	downloadURL string
	httpClient  Doer
}

func NewDownloader(token string) *Downloader {
	return &Downloader{
		token:       token,
		server:      "https://api.telegram.org",
		downloadURL: "%s/file/bot%s/%s",
		httpClient:  http.DefaultClient,
	}
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

func (d *Downloader) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	url := d.buildDownloadURL(d.token, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	res, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		return nil, fmt.Errorf("while download: %v, %v", res.StatusCode, res.Status)
	}

	return res.Body, nil
}

func (d *Downloader) DownloadVoiceMP3(ctx context.Context, path string) (filePath *string, err error) {
	res, err := d.Download(ctx, path)
	if err != nil {
		return nil, err
	}
	fileOutput, err := d.createDownloadFile()
	if err != nil {
		return nil, err
	}
	name := fileOutput.Name()
	defer func() {
		fileOutput.Close()
	}()

	r, err := d.convertToMP3(res)
	if err != nil {
		return &name, err
	}

	_, err = io.Copy(fileOutput, r)
	return &name, nil
}

func (d *Downloader) buildDownloadURL(token, path string) string {
	return fmt.Sprintf(d.downloadURL, d.server, token, path)
}

func (d *Downloader) convertToMP3(res io.ReadCloser) (*io.PipeReader, error) {
	r, w := io.Pipe()
	cmd := exec.Command("ffmpeg", "-i", "-", "-vn", "-acodec", "libmp3lame", "-f", "mp3", "pipe:1")
	cmd.Stdin = res
	cmd.Stdout = w

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		defer w.Close()
		defer res.Close()
		err = cmd.Wait()
		if err != nil {
			logger.Infof("Error in conversion: %v", err)
		}
	}()

	return r, nil
}

func (d *Downloader) createDownloadFile() (*os.File, error) {
	tempDir := "tempVoice"
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return nil, err
	}

	return os.CreateTemp(tempDir, "voice-*.mp3")
}
