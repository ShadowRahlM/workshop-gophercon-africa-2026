package components

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/domain/sighting"
)

const sparkW, sparkH = 100.0, 24.0

func years(sum sighting.Summary) []int {
	ys := make([]int, 0, sum.Latest.Year()-sum.Earliest.Year()+1)
	for y := sum.Earliest.Year(); y <= sum.Latest.Year(); y++ {
		ys = append(ys, y)
	}
	return ys
}

func peak(byYear map[int]int) (year, count int) {
	for y, n := range byYear {
		if n > count || (n == count && y < year) {
			year, count = y, n
		}
	}
	return year, count
}

func sparkPoints(byYear map[int]int, sum sighting.Summary) string {
	ys := years(sum)
	_, high := peak(byYear)
	var b strings.Builder
	for i, y := range ys {
		x := 0.0
		if len(ys) > 1 {
			x = float64(i) / float64(len(ys)-1) * sparkW
		}
		top := sparkH
		if high > 0 {
			top = sparkH - float64(byYear[y])/float64(high)*(sparkH-1)
		}
		fmt.Fprintf(&b, "%.1f,%.1f ", x, top)
	}
	return strings.TrimSpace(b.String())
}

func sparkArea(byYear map[int]int, sum sighting.Summary) string {
	return fmt.Sprintf("0,%v %s %v,%v", sparkH, sparkPoints(byYear, sum), sparkW, sparkH)
}

func sparkLabel(byYear map[int]int, sum sighting.Summary) string {
	y, n := peak(byYear)
	return fmt.Sprintf("Sightings per year, %d to %d, most in %d with %d",
		sum.Earliest.Year(), sum.Latest.Year(), y, n)
}

func bar(count, longest int) map[string]string {
	if longest <= 0 {
		return map[string]string{"--bar": "0%"}
	}
	return map[string]string{"--bar": strconv.Itoa(count*100/longest) + "%"}
}
