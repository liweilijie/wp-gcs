package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"wp-gcs/config"
	"wp-gcs/database"
	"wp-gcs/gcs"
)

var PTotalFiles atomic.Uint64
var CTotalFiles atomic.Uint64
var BUCKET_PREFIX = "wp-content/uploads/"

func consumer(idx int, rep string, buff chan string, wpUploadsHandle database.WpUploadsHandle, gcsHandle gcs.Gcs, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		p, ok := <-buff
		if ok {
			// make the origin path and bucket path
			bucketObj := strings.Replace(p, rep, BUCKET_PREFIX, 1)
			//log.Infof("Consumer %d received from channel: %s=>%s", idx, p, bucketObj)

			// use p or bucketObj to find one in mysql
			//contents, err := wpUploadsHandle.SelectByNames(p, bucketObj)
			contents, err := wpUploadsHandle.SelectByLocalPath(p)
			if err != nil {
				log.Errorf("%s or %s select from mysql error:%v\n", p, bucketObj, err)
				continue
			}

			// if we can find in mysql and nothing to do.
			if len(contents) > 0 {
				//log.Warnf("%s or %s already in mysql\n", p, bucketObj)
				continue
			}
			// otherwise we need to upload to GCS

			ctx := context.Background()
			err = gcsHandle.UploadFile(ctx, p, bucketObj)
			if err != nil {
				log.Errorf("local %s to %s bucket error:%v\n", p, bucketObj, err)
				continue
			}

			// and insert into mysql
			err = wpUploadsHandle.Insert(database.WpUploads{OriginPath: p, BucketPath: bucketObj})
			if err != nil {
				log.Errorf("insert %s to mysql error:%v\n", bucketObj, err)
				continue
			}
			CTotalFiles.Add(1)
			log.Infof("done local %s to %s bucket", p, bucketObj)
		} else {
			log.Infof("Consumer #%d: no more values to process, exiting\n", idx)
			return
		}
	}
}

func run(app *config.AppConfig, wpUploadsHandle database.WpUploadsHandle, gcsHandle gcs.Gcs) {
	var wg sync.WaitGroup
	var workersCount = app.WORK_COUNT

	var buffer = make(chan string, 100000)
	go producer(buffer, app.LOCAL_PATH)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go consumer(i, app.REPLACE_PREFIX_PATH, buffer, wpUploadsHandle, gcsHandle, &wg)
	}

	wg.Wait()
	log.Infof("Total files: %d, consumer total files: %d", PTotalFiles.Load(), CTotalFiles.Load())
}

func producer(buff chan string, dir string) {
	defer close(buff)
	err := dirWalk(dir, buff)
	if err != nil {
		log.Error("producer dir walk error:", err)
	}
}

func dirWalk(dir string, buff chan string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			// Recursively calls Dirwalk in the case of a directory
			err := dirWalk(filepath.Join(dir, file.Name()), buff)
			if err != nil {
				return fmt.Errorf("dirwalk %s: %w", filepath.Join(dir, file.Name()), err)
			}
			continue
		}

		// debug
		//if PTotalFiles.Load() >= 10 {
		//	return nil
		//}

		buff <- filepath.Join(dir, file.Name())
		PTotalFiles.Add(1)
	}

	return nil
}
