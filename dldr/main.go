package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

// Worker struct
type Worker struct {
	URL       string
	File      *os.File
	Count     int64
	SyncWG    sync.WaitGroup
	TotalSize int64
	Progress
}

// Progress struct
type Progress struct {
	Pool *pb.Pool
	Bars []*pb.ProgressBar
}

func main() {
	var t = flag.Bool("t", false, "file name with datetime")
	var workerCount = flag.Int64("c", 5, "connection count")
	flag.Parse()

	var downloadURL string
	fmt.Print("Please enter a URL: ")
	fmt.Scanf("%s", &downloadURL)

	// Get header from the url
	log.Println("Url:", downloadURL)
	fileSize, err := getSizeAndCheckRangeSupport(downloadURL)
	handleError(err)
	log.Printf("File size: %d bytes\n", fileSize)

	var filePath string
	if *t {
		filePath = filepath.Dir(os.Args[0]) + string(filepath.Separator) + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + getFileName(downloadURL)
	} else {
		filePath = filepath.Dir(os.Args[0]) + string(filepath.Separator) + getFileName(downloadURL)
	}
	log.Printf("Local path: %s\n", filePath)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	handleError(err)
	defer f.Close()

	// New worker struct to download file
	var worker = Worker{
		URL:       downloadURL,
		File:      f,
		Count:     *workerCount,
		TotalSize: fileSize,
	}

	var start, end int64
	var partialSize = int64(fileSize / *workerCount)
	now := time.Now().UTC()
	for num := int64(0); num < worker.Count; num++ {
		// New sub progress bar (give it 0 at first for new instance and assign real size later on.)
		bar := pb.New(0).Prefix(fmt.Sprintf("Part %d  0%% ", num))
		bar.ShowSpeed = true
		bar.SetMaxWidth(100)
		bar.SetUnits(pb.U_BYTES_DEC)
		bar.SetRefreshRate(time.Second)
		bar.ShowPercent = true
		worker.Progress.Bars = append(worker.Progress.Bars, bar)

		if num == worker.Count {
			end = fileSize // last part
		} else {
			end = start + partialSize
		}

		worker.SyncWG.Add(1)
		go worker.writeRange(num, start, end-1)
		start = end
	}
	worker.Progress.Pool, err = pb.StartPool(worker.Progress.Bars...)
	handleError(err)
	worker.SyncWG.Wait()
	worker.Progress.Pool.Stop()
	log.Println("Elapsed time:", time.Since(now))
	log.Println("Done!")

}

func (w *Worker) writeRange(partNum int64, start int64, end int64) {
	var written int64
	body, size, err := w.getRangeBody(start, end)
	if err != nil {
		log.Fatalf("Part %d request error: %s\n", partNum, err.Error())
	}
	defer body.Close()
	defer w.Bars[partNum].Finish()
	defer w.SyncWG.Done()

	// Assign total size to progress bar
	w.Bars[partNum].Total = size

	// New percentage flag
	percentFlag := map[int64]bool{}

	// make a buffer to keep chunks that are read
	buf := make([]byte, 4*1024)
	for {
		nr, er := body.Read(buf)
		if nr > 0 {
			nw, err := w.File.WriteAt(buf[0:nr], start)
			if err != nil {
				log.Fatalf("Part %d occured error: %s.\n", partNum, err.Error())
			}
			if nr != nw {
				log.Fatalf("Part %d occured error of short writiing.\n", partNum)
			}

			start = int64(nw) + start
			if nw > 0 {
				written += int64(nw)
			}

			// Update written bytes on progress bar
			w.Bars[int(partNum)].Set64(written)

			// Update current percentage on progress bars
			p := int64(float32(written) / float32(size) * 100)
			_, flagged := percentFlag[p]
			if !flagged {
				percentFlag[p] = true
				w.Bars[int(partNum)].Prefix(fmt.Sprintf("Part %d  %d%% ", partNum, p))
			}
		}
		if er != nil {
			if er.Error() == "EOF" {
				if size == written {
					// Download successfully
				} else {
					handleError(fmt.Errorf("Part %d unfinished", partNum))
				}
				break
			}
			handleError(fmt.Errorf("Part %d occured error: %s", partNum, er.Error()))
		}
	}
}

func (w *Worker) getRangeBody(start int64, end int64) (io.ReadCloser, int64, error) {
	var client http.Client
	req, err := http.NewRequest("GET", w.URL, nil)
	// req.Header.Set("cookie", "")
	// log.Printf("Request header: %s\n", req.Header)
	if err != nil {
		return nil, 0, err
	}

	// Set range header
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	size, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	return resp.Body, size, err
}

func getSizeAndCheckRangeSupport(url string) (size int64, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	// req.Header.Set("cookie", "")
	// log.Printf("Request header: %s\n", req.Header)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	log.Printf("Response header: %v\n", res.Header)
	header := res.Header
	acceptRanges, supported := header["Accept-Ranges"]
	if !supported {
		return 0, errors.New("Doesn't support header `Accept-Ranges`")
	} else if supported && acceptRanges[0] != "bytes" {
		return 0, errors.New("Support `Accept-Ranges`, but value is not `bytes`")
	}
	size, err = strconv.ParseInt(header["Content-Length"][0], 10, 64)
	return
}

func getFileName(downloadURL string) string {
	urlStruct, err := url.Parse(downloadURL)
	handleError(err)
	return filepath.Base(urlStruct.Path)
}

func handleError(err error) {
	if err != nil {
		log.Println("err:", err)
		os.Exit(1)
	}
}
