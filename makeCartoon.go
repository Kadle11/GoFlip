package main

import (
	"log"
	"gocv.io/x/gocv"
	"image"
)

func convComic(srcPath string) gocv.Mat {

	srcRGB := gocv.IMRead(srcPath, gocv.IMReadColor)
	srcEdit := srcRGB
	
	nDownSamples := 2
	nBilateralFilters := 10

	var startPt image.Point
	tmpImg := gocv.NewMat()
	
	/*Downsample using Gaussian Pyramid*/

	for i:=0; i < nDownSamples; i++ {
		gocv.PyrDown(srcEdit, &srcEdit, startPt, gocv.BorderDefault)
	}

	/*Apply the Bilateral Filter*/
	for i:=0; i < nBilateralFilters; i++ {
		gocv.BilateralFilter(srcEdit, &tmpImg, 30, 30, 24)
		tmpImg.CopyTo(&srcEdit)
	}

	/*Upsample using Gaussian Pyramid*/

	for i:=0; i < nDownSamples; i++ {
		gocv.PyrUp(srcEdit, &srcEdit, startPt, gocv.BorderDefault)
	}

	/*Covert to Grayscale and Apply Median Blur*/
	srcGray := gocv.NewMat()
	srcBlur := gocv.NewMat()

	gocv.CvtColor(srcRGB, &srcGray, gocv.ColorRGBToGray) // 7
	gocv.MedianBlur(srcGray, &srcBlur, 3)
	
	/*Detect and Enhance Edges*/
	srcEdge := gocv.NewMat()
	gocv.AdaptiveThreshold(srcBlur, &srcEdge, 255.0, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 7, 2)
	gocv.CvtColor(srcEdge, &srcEdge, gocv.ColorGrayToBGR)
	
	gocv.IMWrite("Edges.png", srcEdge)

	gocv.BitwiseAnd(srcEdit, srcEdge, &srcEdit)
	return srcEdit
}

func main() {

	log.Println()
	comicImg := convComic("ImgTest.jpg")
	log.Println()
	gocv.IMWrite("Hmm.jpg", comicImg)
}
	
