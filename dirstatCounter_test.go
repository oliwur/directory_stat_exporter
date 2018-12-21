package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func setupTmpDir() string {
	tmpDir, err := ioutil.TempDir("", "dirstat_test_")
	if err != nil {
		log.Fatal("could not create temp dir", err)
	}
	return tmpDir
}

func TestGetFileCountZero(t *testing.T) {
	tmpDir := setupTmpDir()
	defer os.RemoveAll(tmpDir)

	t.Run("given a empty directory when counted not recursively then it should return 0", func(t *testing.T) {
		count := getFileCount(tmpDir, false)

		if count != 0 {
			t.Fail()
			t.Errorf("there are no documents in this folder. the count must be zero. it is %d!\n", count)
		}
	})

	t.Run("given an empty dir when counted recursively then it should return 0", func(t *testing.T) {
		count := getFileCount(tmpDir, true)
		if count != 0 {
			t.Fail()
			t.Error("the value shoud be 0 but was", count)
		}
	})
}

func TestGetFileCountThree(t *testing.T) {
	tmpDir := setupTmpDir()
	defer os.RemoveAll(tmpDir)

	file1, _ := ioutil.TempFile(tmpDir, "file1")
	file2, _ := ioutil.TempFile(tmpDir, "file2")
	file3, _ := ioutil.TempFile(tmpDir, "file3")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())
	defer os.Remove(file3.Name())

	t.Run("given a directory with 3 files when counted not recursively then it should return 3", func(t *testing.T) {
		count := getFileCount(tmpDir, false)

		if count != 3 {
			t.Fail()
			t.Errorf("there are no documents in this folder. the count must be zero. it is %d!\n", count)
		}
	})

	t.Run("given a directory with 3 files when counted recursively then it should return 3", func(t *testing.T) {
		count := getFileCount(tmpDir, true)

		if count != 3 {
			t.Fail()
			t.Errorf("there are 3 files in the directory, it reported %d\n", count)
		}
	})
}

func TestGetFileCountNotExisting(t *testing.T) {
	tmpDir := setupTmpDir()
	defer os.RemoveAll(tmpDir)

	t.Run("given the directory does not exist when counted non recursively then it should return -1 indicating an error", func(t *testing.T) {
		count := getFileCount(tmpDir+"_this_dir_does_not_exist", false)

		if count != -1 {
			t.Fail()
			t.Errorf("it should return -1, because the dir does not exist. it returned %d instead.\n", count)
		}
	})

	t.Run("given the directory does not exist when counted recursively then it should return -1 indicating an error", func(t *testing.T) {
		count := getFileCount(tmpDir+"_this_dir_does_not_exist", true)

		if count != -1 {
			t.Fail()
			t.Errorf("it should return -1, because the dir does not exist. it returned %d instead.\n", count)
		}
	})
}
