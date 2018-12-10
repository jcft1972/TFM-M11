package main

import (
	"os"
)

func saveFile(nameFile string, data []byte) {
	file, err := os.Create(nameFile)
	failOnError(err, "Fail in create file...")
	defer file.Close()

	_, err = file.Write(data)
	failOnError(err, "Fail in write file...")

}

func removeFile(nameFile string) error {
	return(os.Remove(nameFile))
}
