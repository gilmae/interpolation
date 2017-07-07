package main

import (
  "image"
  "image/color"
  "image/draw"
  "image/jpeg"
  "os"
  "github.com/gilmae/interpolation"
)

func main() {
  var xSequence = []float64{0.0, .16, .42, .6425, .8675, 1}
  var redpoints =  []float64{0.0,32.0,237.0,255.0,0.0,0.0}
  var greenpoints =  []float64{7.0,107.0,255.0, 170.0, 2.0, 7.0}
  var bluepoints =  []float64{100.0, 203.0,255.0, 0.0, 0.0, 100.0}

  var redInterpolant = interpolation.CreateMonotonicCubic(xSequence, redpoints)
  var greenInterpolant = interpolation.CreateMonotonicCubic(xSequence, greenpoints)
  var blueInterpolant = interpolation.CreateMonotonicCubic(xSequence, bluepoints)

   var black = color.NRGBA{0,0,0,255}

   bounds := image.Rect(0,0,1000,255)
   b := image.NewNRGBA(bounds)
   draw.Draw(b, bounds, image.NewUniform(color.White), image.ZP, draw.Src)

   for ii := 0; ii < 1000; ii++ {
     var point = float64(ii)/1000.0
     var redpoint = redInterpolant(point)
     var greenpoint = greenInterpolant(point)
     var bluepoint = blueInterpolant(point)

     c := color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), base}

     for ij := 0; ij < 255; ij++ {
       b.Set(ii, ij, c)
     }
     b.Set(ii, 255-int(redpoint), black)
     b.Set(ii, 255-int(greenpoint), black)
     b.Set(ii, 255-int(bluepoint), black)

   }

   file, _ := os.Create("gradient.jpg")
   jpeg.Encode(file,b, &jpeg.Options{jpeg.DefaultQuality})
   file.Close()
}
