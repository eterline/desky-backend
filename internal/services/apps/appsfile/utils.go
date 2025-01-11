package appsfile

import (
	"encoding/json"
	"os"
)

func testFile(file string) error {

	_, err := os.Stat(file)
	if err != nil {
		return ErrCannotOpen(file)
	}

	return nil
}

func genFile() error {
	file, err := os.Create(DefaultPath)
	defer file.Close()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(ExampleTable, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
