package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const api string = "http://localhost:8080/api/"
const session_name = "session_1"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CsvFormat struct {
	Nama  string
	NoWa  string
	Lomba string
}

func StartWorker(username string, number string, message string, session_name string) {
	var res Response
	bodyForm := map[string]string{
		"sessions": session_name,
		"target":   number,
		"message":  message,
	}
	body, _ := json.Marshal(bodyForm)
	request, err := http.NewRequest("POST", api+"sendtext", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	readBody, _ := io.ReadAll(response.Body)
	if err := json.Unmarshal(readBody, &res); err != nil {
		panic(err)
	}
	val := fmt.Sprintf("Number : %s -> Status -> %v ", number, res.Status)
	fmt.Println(val)
}

func ReadCsv(filename string) []CsvFormat {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	var csvFormat []CsvFormat
	for i, line := range data {
		if i > 0 { // omit header line
			var rec CsvFormat
			for j, field := range line {
				if j == 1 {
					rec.Nama = field
				} else if j == 4 {
					rec.NoWa = field
				} else if j == 5 {
					rec.Lomba = field
				}
			}
			csvFormat = append(csvFormat, rec)
		}
	}
	return csvFormat
}

func StartJob(csvinput []CsvFormat) {
	var gmeet_link string = ""
	var kahoot_link string = ""
	var message string
	for _, x := range csvinput {
		message = fmt.Sprintf(`Halo %s, Selamat kamu telah terpilih menjadi peserta lomba %s,berikut adalah link untuk memasuki lomba tersebut : %s , dan jangan lupa juga untuk memasuki link live gmeet berikut ini %s untuk melihat soal yang tersedia`, x.Nama, x.Lomba, kahoot_link, gmeet_link)
		StartWorker(x.Nama, x.NoWa, message, session_name)
	}
}

func main() {
	var csvInput string
	fmt.Print("csv: ")
	fmt.Scanln(&csvInput)
	res := ReadCsv(csvInput)
	StartJob(res)
}
