package main

// #cgo darwin pkg-config: zbar
// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"
import (
	"bytes"
	"fmt"
	"image/png"
	"log"

	"github.com/clsung/grcode"
	"gocv.io/x/gocv"
)

func getDataFromImage(img gocv.Mat) (results []string, err error) {
	imgBytes, _ := gocv.IMEncode("img.png", img)
	m, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Printf("decode file error: %v", err)
		return results, err
	}
	scanner := grcode.NewScanner()
	defer scanner.Close()

	scanner.SetConfig(0, C.ZBAR_CFG_ENABLE, 1)
	zImg := grcode.NewZbarImage(m)
	defer zImg.Close()

	scanner.Scan(zImg)
	symbol := zImg.GetSymbol()
	for ; symbol != nil; symbol = symbol.Next() {
		results = append(results, symbol.Data())
	}
	return results, nil
}

func main() {
	deviceID := 0

	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		log.Printf("error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("QR Code Reader")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	log.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			log.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		results, err := getDataFromImage(img)
		if err != nil {
			log.Fatal(err)
		}
		if len(results) == 0 {
			log.Printf("No qrcode detected")
		}
		for _, result := range results {
			fmt.Printf("%s\n", result)
		}

		window.IMShow(img)
		window.WaitKey(1)
	}
}
