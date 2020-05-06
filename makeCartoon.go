package main

import (
	"log"
	"gocv.io/x/gocv"
	"image"
)

func convComic(srcPath string) gocv.Mat {

	srcRGB := gocv.IMRead(srcPath, gocv.IMReadColor)
	srcEdit := srcRGB
	
	nDownSamples := 0
	nBilateralFilters := 10

	tmpImg := gocv.NewMat()
	defer tmpImg.Close()
	/*Downsample using Gaussian Pyramid*/

	for i:=0; i < nDownSamples; i++ {
		imgSize := srcEdit.Size()	
		startPt := image.Pt(imgSize[0]/2, imgSize[1]/2)
		gocv.PyrDown(srcEdit, &srcEdit, startPt, gocv.BorderDefault)
	}

	/*Apply the Bilateral Filter*/
	for i:=0; i < nBilateralFilters; i++ {
		gocv.BilateralFilter(srcEdit, &tmpImg, 30, 30, 24)
		tmpImg.CopyTo(&srcEdit)
	}

	/*Upsample using Gaussian Pyramid*/
	
	for i:=0; i < nDownSamples; i++ {
		imgSize := srcEdit.Size()	
		endPt := image.Pt(imgSize[0]*2, imgSize[1]*2)
		gocv.PyrUp(srcEdit, &srcEdit, endPt, gocv.BorderDefault)
	}

	/*Covert to Grayscale and Apply Median Blur*/
	srcGray := gocv.NewMat()
	srcBlur := gocv.NewMat()
	defer srcGray.Close()
	defer srcBlur.Close()
	
	gocv.CvtColor(srcRGB, &srcGray, gocv.ColorRGBToGray) // 7
	gocv.MedianBlur(srcGray, &srcBlur, 3)
	
	/*Detect and Enhance Edges*/
	srcEdge := gocv.NewMat()
	defer srcEdge.Close()
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
	
