package fun

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func GetQuote() string {
	csvfile, err := os.Open("/bin/fun/quotes.csv")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer csvfile.Close()

	var lines []map[string]string
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if record[1] == "" || record[1] == " " {
			record[1] = "Anonymous"
		}

		lines = append(lines, map[string]string{"author": record[0], "quote": record[1]})
	}

	rand.Seed(time.Now().UnixNano())

	min := 0
	max := len(lines) - 1
	stt := rand.Intn(max-min+1) + min

	if lines != nil {
		return lines[stt]["quote"] + " By " + lines[stt]["author"]
	} else {
		return ""
	}
}
