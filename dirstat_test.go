package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func setupTestFileWithTimestamp(dir string, file string, ts time.Time) *os.File {
	f, err := ioutil.TempFile(dir, file)
	if err != nil {
		log.Fatal("could not create temp file")
	}
	err = os.Chtimes(f.Name(), ts, ts)
	if err != nil {
		log.Fatal("could not change timestamp of file")
	}
	return f
}

func TestOldestFileInDir(t *testing.T) {
	tmpDir := setupTmpDir()
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			log.Printf("could not remove tmpDir")
		}
	}()

	setupTestFileWithTimestamp(tmpDir, "test", time.Now().Add(time.Second*time.Duration(-20)))

	t.Run("given a dir with one file with 20 seconds of age when analysed non recursively then return 20", func(t *testing.T) {
		age := time.Now().Unix() - getOldestFileModTimestamp(tmpDir, false)
		if age != 20 {
			t.Fail()
			t.Errorf("the file age is not 20, it's %v\n", age)
		}
	})

	t.Run("given a dir with one file with 20 seconds of age when analysed recursively then return 20", func(t *testing.T) {
		age := time.Now().Unix() - getOldestFileModTimestamp(tmpDir, true)
		if age != 20 {
			t.Fail()
			t.Errorf("the file age is not 20, it's %v\n", age)
		}
	})
}
