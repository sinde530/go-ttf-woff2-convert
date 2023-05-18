package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/", handleConversion)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func handleConversion(c *gin.Context) {
	file, header, err := c.Request.FormFile("fontfile")
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	defer file.Close()

	// Read the TTF font file
	ttfData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	var buffer bytes.Buffer
	writer := brotli.NewWriterLevel(&buffer, brotli.BestCompression)
	if _, err = io.Copy(writer, bytes.NewReader(ttfData)); err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	writer.Close()

	filename := header.Filename[:len(header.Filename)-4] + ".woff2"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
}
