package interpolation

import (
  "math"
)

type MonotonicCubic func(x float64) float64

func CreateMonotonicCubic(xs []float64 , ys []float64) MonotonicCubic {
  var i = len(xs)
  var length = i

  if (length != len(ys)) {
    return nil
  }//, errors.New("Interpolation: Need an equal count of xs and ys.") }
  if (length == 0) {
    return func(x float64) float64 {
      return 0
    }
  }
  if (length == 1) {
    // Impl: Precomputing the result prevents problems if ys is mutated later and allows garbage collection of ys
    // Impl: Unary plus properly converts values to numbers
    var result = ys[0]
    return func(x float64) float64 {
      return result
    }
  }

  var indexes = make([]int, length)
	for i := 0; i < length; i++ {
    indexes[i] = i
  }

  //TMP: Assume the xs are in ASC order
  //indexes.sort(function(a, b) { return xs[a] < xs[b] ? -1 : 1; })

  // Get consecutive differences and slopes
  var dys = []float64{}
  var dxs = []float64{}
  var ms = []float64{}

  var dx float64 = 0.0
  var dy float64 = 0.0

  for i := 0; i < length - 1; i++ {
    dx = xs[i + 1] - xs[i]
    dy = ys[i + 1] - ys[i]

    dxs = append(dxs, dx)
    dys = append(dys, dy)
    ms = append(ms, dy/dx)
  }

  // Get degree-1 coefficients
	var degree1Coefficients = []float64 {ms[0]}
	for i = 0; i < len(dxs) - 1; i++ {
		var m = ms[i]
    var mNext = ms[i + 1];
		if (m*mNext <= 0) {
			degree1Coefficients = append(degree1Coefficients, 0)
		} else {
			var dx_ = dxs[i]
      var dxNext = dxs[i + 1]
      var common = dx_ + dxNext

      degree1Coefficients = append(degree1Coefficients, 3*common/((common + dxNext)/m + (common + dx_)/mNext));
		}
	}
	degree1Coefficients = append(degree1Coefficients, ms[len(ms) - 1])

  // Get degree-2 and degree-3 coefficients
  var degree2Coefficients = []float64{}
  var degree3Coefficients = []float64{}

  for i = 0; i < len(degree1Coefficients) - 1; i++ {
    var c1 = degree1Coefficients[i]
    var m_ = ms[i]
    var invDx = 1/dxs[i]
    var common_ = c1 + degree1Coefficients[i + 1] - m_ - m_

    degree2Coefficients = append(degree2Coefficients, (m_ - c1 - common_)*invDx)
    degree3Coefficients = append(degree3Coefficients, common_*invDx*invDx)
  }

  // Return interpolant function
  return func(x float64) float64 {
    // The rightmost point in the dataset should give an exact result
    var i = len(xs) - 1
    if (x == xs[i]) {
      return ys[i]
    }

    // Search for the interval x is in, returning the corresponding y if x is one of the original xs
    var low = 0
    var mid int = len(degree3Coefficients) - 1
    var high = len(degree3Coefficients) - 1

    for low <= high {
      mid = int(math.Floor(0.5*(float64(low) + float64(high))))
      var xHere = xs[mid]
      if (xHere < x) {
        low = mid + 1
      } else if (xHere > x) {
        high = mid - 1
      } else {
        return ys[mid]
      }
    }
    i = int(math.Max(0.0, float64(high)));

    // Interpolate
    var diff = x - xs[i]
    var diffSq = diff*diff
    return ys[i] + degree1Coefficients[i]*diff + degree2Coefficients[i]*diffSq + degree3Coefficients[i]*diff*diffSq
  }

  return nil
}
