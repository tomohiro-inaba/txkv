package txkv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type CSVStorage struct {
	path string
	file *os.File
}

func NewCSVStorage(path string) (*CSVStorage, error) {
	file, err := openFile(path)
	if err != nil {
		return nil, err
	}
	return &CSVStorage{path, file}, nil
}

func openFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil

}

// TODO: fix
func (csv *CSVStorage) close() {
	if err := csv.file.Close(); err != nil {
		log.Fatal(err)
	}
}

func (csv *CSVStorage) Read(key Key) (Value, error) {
	exists := false
	ret := Value("")

	csv.file.Seek(0, 0)
	scanner := bufio.NewScanner(csv.file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if Key(line[0]) == key {
			exists = true
			ret = Value(line[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	if exists {
		return ret, nil
	}
	return "", errors.New(fmt.Sprintf("Not found: key=%s\n", key))
}

func (csv *CSVStorage) Write(key Key, value Value) error {
	line := fmt.Sprintf("%s,%s\n", key, value)
	if _, err := csv.file.WriteString(line); err != nil {
		return err
	}
	return nil
}

func (csv *CSVStorage) GC() (bool, error) {
	exists := false
	tmp := make(map[Key]Value)

	csv.file.Seek(0, 0)
	scanner := bufio.NewScanner(csv.file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		exists = true
		tmp[Key(line[0])] = Value(line[1])
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// TODO: fix
	os.Remove(csv.path)

	file, err := openFile(csv.path)
	if err != nil {
		return false, err
	}

	for k, v := range tmp {
		line := fmt.Sprintf("%s,%s\n", k, v)
		if _, err := file.WriteString(line); err != nil {
			return false, err
		}
	}
	csv.file = file

	return true, nil
}
