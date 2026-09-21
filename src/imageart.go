package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

var quadrantGlyphs = [16]rune{
	' ', '▘', '▝', '▀',
	'▖', '▌', '▞', '▛',
	'▗', '▚', '▐', '▜',
	'▄', '▙', '▟', '█',
}

type rgbColor struct{ r, g, b uint32 }
func displayImage(path string, width int) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return false
	}

	printImageQuadrant(src, width)
	return true
}

func printImageQuadrant(src image.Image, width int) {
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return
	}

	height := width * srcH / srcW / 2
	if height < 1 {
		height = 1
	}

	sampleW, sampleH := width*2, height*2

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			tl := sampleAvgColor(src, bounds, col*2, row*2, sampleW, sampleH)
			tr := sampleAvgColor(src, bounds, col*2+1, row*2, sampleW, sampleH)
			bl := sampleAvgColor(src, bounds, col*2, row*2+1, sampleW, sampleH)
			br := sampleAvgColor(src, bounds, col*2+1, row*2+1, sampleW, sampleH)

			printQuadrantCell(tl, tr, bl, br)
		}
		fmt.Println()
	}
}

func printQuadrantCell(tl, tr, bl, br rgbColor) {
	corners := [4]rgbColor{tl, tr, bl, br}

	bestMask := 0
	bestErr := int64(-1)

	for mask := 0; mask < 16; mask++ {
		fg, bg := splitColors(corners, mask)

		var totalErr int64
		for i, c := range corners {
			ref := bg
			if mask&(1<<i) != 0 {
				ref = fg
			}
			totalErr += colorSqDist(c, ref)
		}

		if bestErr == -1 || totalErr < bestErr {
			bestErr = totalErr
			bestMask = mask
		}
	}

	fg, bg := splitColors(corners, bestMask)
	fmt.Printf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c\x1b[0m",
		fg.r, fg.g, fg.b, bg.r, bg.g, bg.b, quadrantGlyphs[bestMask])
}

func splitColors(corners [4]rgbColor, mask int) (fg, bg rgbColor) {
	var fgSum, bgSum rgbColor
	var fgCount, bgCount int64

	for i, c := range corners {
		if mask&(1<<i) != 0 {
			fgSum.r += c.r
			fgSum.g += c.g
			fgSum.b += c.b
			fgCount++
		} else {
			bgSum.r += c.r
			bgSum.g += c.g
			bgSum.b += c.b
			bgCount++
		}
	}
	return avgColor(fgSum, fgCount), avgColor(bgSum, bgCount)
}

func avgColor(sum rgbColor, count int64) rgbColor {
	if count == 0 {
		return rgbColor{0, 0, 0}
	}
	return rgbColor{sum.r / uint32(count), sum.g / uint32(count), sum.b / uint32(count)}
}

func colorSqDist(a, b rgbColor) int64 {
	dr := int64(a.r) - int64(b.r)
	dg := int64(a.g) - int64(b.g)
	db := int64(a.b) - int64(b.b)
	return dr*dr + dg*dg + db*db
}

func sampleAvgColor(src image.Image, bounds image.Rectangle, col, row, outW, outH int) rgbColor {
	srcW, srcH := bounds.Dx(), bounds.Dy()

	x0 := bounds.Min.X + col*srcW/outW
	x1 := bounds.Min.X + (col+1)*srcW/outW
	y0 := bounds.Min.Y + row*srcH/outH
	y1 := bounds.Min.Y + (row+1)*srcH/outH
	if x1 <= x0 {
		x1 = x0 + 1
	}
	if y1 <= y0 {
		y1 = y0 + 1
	}

	var rSum, gSum, bSum, count uint32
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			rSum += r >> 8
			gSum += g >> 8
			bSum += b >> 8
			count++
		}
	}
	if count == 0 {
		return rgbColor{0, 0, 0}
	}
	return rgbColor{rSum / count, gSum / count, bSum / count}
}