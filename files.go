package gormmigrations

import (
	"io/ioutil"
)

func writeFile(name string, content string) error {
	err := ioutil.WriteFile(name+".go", []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}
