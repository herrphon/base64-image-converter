package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"image"
	_ "image/jpeg"
	_ "image/png"
)

func decodeConfig(base64Img string) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img))
	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Width:", config.Width, "Height:", config.Height, "Format:", format)
}

func getImageFromBase64(base64Img string) (image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img))

	img, _, err := image.Decode(reader)

	return img, err
}

func getBlankImage() image.Image {
	return image.NewRGBA(image.Rect(0, 0, 16, 16))
}

func getBase64FromImage(imgLocation string) (string, error) {
	f, e := os.Open(imgLocation)
	if e != nil {
		return "", e
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}
