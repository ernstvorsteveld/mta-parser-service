package sta

import (
	"log"
	"os"
	"testing"
)

func Test_read_complete_file(t *testing.T) {
	file, err := os.Open("bunq.sta")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	handleFile(file)
}
