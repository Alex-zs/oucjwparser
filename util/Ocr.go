package util

import (
	"github.com/Comdex/imgo"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)


// 裁剪图片
func clipImg(srcImg image.Image) *image.Image{
	rgbImg := srcImg.(*image.RGBA)
	dstImg := rgbImg.SubImage(image.Rect(1, 1, rgbImg.Rect.Dx() - 1, rgbImg.Rect.Dy() - 1))
	return &dstImg
}

// 灰度化
func grayImg(srcImg *image.Image)  *image.Gray{
	rgbImg := (*srcImg).(*image.RGBA)
	dstImg := image.NewGray(rgbImg.Rect)
	fromsRGBLUT16 := getsRGBToLinearRGB16LUT()
	intPixel := make([]int, 1)
	for y := 0; y < rgbImg.Rect.Dy(); y++ {
		for x := 0; x < rgbImg.Rect.Dx(); x++ {
			offset := y * rgbImg.Stride + x * 4
			red, grn, blu := rgbImg.Pix[offset], rgbImg.Pix[offset+1], rgbImg.Pix[offset+2]
			newRed := uint32((*fromsRGBLUT16)[red]) & uint32(0xffff)
			newGrn := uint32((*fromsRGBLUT16)[grn]) & uint32(0xffff)
			newBlu := uint32((*fromsRGBLUT16)[blu]) & uint32(0xffff)
			gray := ((0.2125 * float64(newRed)) +
					 (0.7154 * float64(newGrn)) +
					 (0.0721 * float64(newBlu))) / 65535.0

			intPixel[0] = int(gray * float64((1 << 8) - 1) + 0.5)

			inData := make([]byte, 1)
			inData[0] = byte(0xff & intPixel[0])

			dstImg.Set(x, y, color.Gray{Y: inData[0]})
		}
	}
	return dstImg
}

// 二值化处理
func binaryImg(srcImg *image.Gray, threshold uint8) * image.Gray{
	for y := 0; y < srcImg.Rect.Dy(); y++ {
		for x := 0; x < srcImg.Rect.Dx(); x++ {
			gray := srcImg.GrayAt(x, y).Y
			if gray > threshold {
				gray = 255
			}else {
				gray = 0
			}
			srcImg.Set(x, y, color.Gray{Y: gray})
		}
	}
	return srcImg
}


var s8Tol16 []int16
func getsRGBToLinearRGB16LUT()  *[]int16{
	if s8Tol16 == nil {
		s8Tol16 = make([]int16, 256)
		var input, output float64
		for i := 0; i <= 255; i++ {
			input = float64(i) / 255.0
			if input <= 0.04045 {
				 output = input / 12.92
			}else {
				output = math.Pow((input + 0.055)/ 1.055, 2.4)
			}
			s8Tol16[i] = int16(math.Round(output * 65535.0))
		}
	}
	return &s8Tol16
}


// 转化图像为png
func convert2Png(srcImgPath , dstImgPath string) {
	imgo.SaveAsPNG(dstImgPath, imgo.MustRead(srcImgPath))
}

// 图片预处理
func ImgProcess(srcImgPath , dstImgPath string)  {
	convert2Png(srcImgPath, srcImgPath)
	imgFile, err := os.Open(srcImgPath)
	if err != nil {
		Log("打开图片失败", err.Error())
	}
	img, err := png.Decode(imgFile)

	err = imgFile.Close()
	if err != nil {
		Log("关闭图片失败", err.Error())
	}

	img = binaryImg(grayImg(clipImg(img)), 58)
	dstFile, _ := os.Create(dstImgPath)
	png.Encode(dstFile, img)
	dstFile.Close()
}

