package main

import (
	"log"
	"os"
	"gocv.io/x/gocv"
)

func getVideoFrames(srcPath string) ([]gocv.Mat, int, int) {

	/*Initialization*/
	srcVideo, err := gocv.OpenVideoCapture(srcPath)
	var allFrames []gocv.Mat
	if err != nil {
		log.Panic("Error opening Video File: %v\n", srcVideo)
		return allFrames, 0, 0
	}
	defer srcVideo.Close()

	/*Get Dimensions*/
	frameWidth := int(srcVideo.Get(gocv.VideoCaptureFrameWidth))
	frameHeight := int(srcVideo.Get(gocv.VideoCaptureFrameHeight))

	/*Sample Frames*/
	SamplingRate := 5
	frame := gocv.NewMat()
	x := 0
	for {
		if ok := srcVideo.Read(&frame); !ok {
			log.Println("Cannot Read Video: ", srcPath)
			break
		}

		if frame.Empty() {
			log.Println("No Image Found")
			break
		}
		if x % SamplingRate == 0 {
			allFrames = append(allFrames, frame)
			frame = gocv.NewMat()
		}
		x++
	}
	log.Println(x, "Frames Sampled")
	return allFrames, frameWidth, frameHeight
}

func makeFlipBook(allFrames []gocv.Mat, frameWidth int, frameHeight int, fps float64, dstPath string) {


	flipBook, err := gocv.VideoWriterFile(dstPath, "X264", fps, frameWidth, frameHeight, true)
	if err != nil {
		log.Panic("Can not open video writer")
		return
	}
	defer flipBook.Close()

	for i:=0; i<len(allFrames); i++ {
		err = flipBook.Write(allFrames[i])
		if err != nil {
			log.Panic("makeFlipBook: Unable to write frame")
			break
		}
	}

}

func main() {

	allFrames, FrameWidth, FrameHeight := getVideoFrames("../data/"+os.Args[1])
	allFrames = convAllFrames(allFrames)
	makeFlipBook(allFrames, FrameWidth, FrameHeight, 10, "../data/"+os.Args[2])

}
