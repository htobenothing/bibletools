package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BibleSubject struct {
	Index   int    `json:"index"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Tag     string `json:"tag"`
}

// var bsList []bible_subject

func startServer(bsList []BibleSubject) {
	router := gin.Default()

	router.GET("/subject/all", func(c *gin.Context) {

		c.JSON(http.StatusOK, bsList)
	})

	r := rand.New(rand.NewSource(99))
	router.GET("/subject/random", func(c *gin.Context) {

		index := r.Intn(len(bsList))

		c.JSON(http.StatusOK, bsList[index])
	})

	router.GET("/subject/range", func(c *gin.Context) {

		start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 0)
		if err != nil {
			fmt.Printf("error %v\n", err)
		}

		end, err := strconv.ParseInt(c.DefaultQuery("end", strconv.Itoa(len(bsList))), 10, 0)
		if err != nil {
			fmt.Printf("error %v\n", err)
		}
		sublist := bsList[start:end]
		c.JSON(http.StatusOK, sublist)
	})

	router.GET("/subject", func(c *gin.Context) {
		numString := c.DefaultQuery("index", "0")

		num, err := strconv.Atoi(numString)
		if err != nil || num < 0 || num > (len(bsList)-1) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid index of scripters",
			})
			return
		}
		c.JSON(http.StatusOK, bsList[num])

	})

	// m := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("localhost:3001"),
	// 	Cache:      autocert.DirCache("./cache"),
	// }
	// s := &http.Server{
	// 	Addr:      ":https",
	// 	TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	// 	Handler:   router,
	// }

	// err := s.ListenAndServeTLS("", "")
	// if err != nil {
	// 	log.Fatal("error %v", err)
	// }

	router.RunTLS(":3001", "./cert/certificate.crt", "./cert/private.key")
}

func main() {

	bsList, err := readDataFromCSV("./data/bible_subject.csv")
	if err != nil {
		log.Fatal("error")
	}

	startServer(bsList)
}

func readDataFromCSV(path string) ([]BibleSubject, error) {

	bsList := []BibleSubject{}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("open file fail")
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	// reader.Comma = ','
	linecount := 0

	for {

		record, err := reader.Read()
		if linecount == 0 {
			linecount += 1
			continue
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New("read file fail")
		}

		// fmt.Println(record)
		bs := BibleSubject{}
		bs.Index = linecount - 1
		for i := 0; i < len(record); i++ {

			bs.Name = record[0]
			bs.Subject = record[1]
			bs.Tag = record[2]

		}

		// fmt.Printf("index %d, name %s, subject %s, tag %s\n", bs.index, bs.name, bs.subject, bs.tag)
		bsList = append(bsList, bs)
		linecount = linecount + 1
	}

	return bsList, nil

}
