package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
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

func TestModTimeIfFileDoesNotExist(t *testing.T) {
	tmpDir := setupTmpDir()
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			log.Printf("could not remove tmpDir")
		}
	}()

	// short explain why: it it would return -1, then -1 will always be the oldest file in a dir. it should not change the oldest file
	// todo: make better / more clear test for this.
	t.Run("given a dir with no file in it when analysed non recursively then return current timestamp", func(t *testing.T) {
		age := getModTime(path.Join(tmpDir, "a-file-that-does-not.exist"))
		current := time.Now().Unix()
		if age < (current-5) || age > (current+1) {
			t.Fail()
			t.Errorf("the result should be %v, it was %v\n", current, age)
		}
	})
}
