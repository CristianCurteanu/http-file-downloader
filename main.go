package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

func main() {
	filename := flag.String("file", "urls.json", "Path to files that contains the URLs")
	poolSize := flag.Int("poolsize", 3, "Number of concurrent downloading threads")
	filesFormat := flag.String("format", "mp4", "Set the downloadable files format. Note, that all files have to have same MIME type")
	flag.Parse()

	jsonFile, err := os.Open(*filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	var urls []string

	err = json.Unmarshal(byteValue, &urls)
	if err != nil {
		fmt.Println("err:", err)
	}

	httpDownloader := NewVideoDownloader()
	httpDownloader.SetBufferSize(*poolSize).SetUrls(urls).SetFilesFormat(*filesFormat)

	httpDownloader.Download()
}

type HttpDownloader struct {
	queue       []string
	urls        []string
	bufferSize  int
	downloaded  int
	cond        *sync.Cond
	filesFormat string
}

func NewVideoDownloader() *HttpDownloader {
	return &HttpDownloader{
		cond:        sync.NewCond(&sync.Mutex{}),
		bufferSize:  runtime.NumCPU(),
		filesFormat: "mp4",
	}
}

func (vd *HttpDownloader) SetBufferSize(size int) *HttpDownloader {
	vd.bufferSize = size

	return vd
}

func (vd *HttpDownloader) SetUrls(urls []string) *HttpDownloader {
	vd.urls = urls

	return vd
}

func (vd *HttpDownloader) SetFilesFormat(format string) *HttpDownloader {
	vd.filesFormat = format

	return vd
}

func (vd *HttpDownloader) Download() {
	for i, url := range vd.urls {
		vd.cond.L.Lock()

		for len(vd.queue) == vd.bufferSize {
			vd.cond.Wait()
		}
		vd.queue = append(vd.queue, url)
		fmt.Println("Adding to queue #", i)

		go vd.getFile(url, i)
		vd.cond.L.Unlock()
	}
	vd.cond.L.Lock()
	for vd.downloaded != len(vd.urls) {
		vd.cond.Wait()
	}
	vd.cond.L.Unlock()
}


func (vd *HttpDownloader) getFile(url string, id int) {
	filename := fmt.Sprintf("%d.%s", id, vd.filesFormat)
	err := DownloadFile(filename, url)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("Downloaded", filename)
	vd.cond.L.Lock()
	vd.queue = vd.queue[1:]
	vd.downloaded += 1
	vd.cond.L.Unlock()
	vd.cond.Signal()
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
