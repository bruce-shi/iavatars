package main

import (
	"bytes"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"github.com/hsluv/hsluv-go"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	dpi       = float64(72)
	fontBytes []byte
	font      *truetype.Font
)

func generateImage(text string, size int, hue float64) (*bytes.Buffer, error) {
	var buf = new(bytes.Buffer)
	fontSize := float64(size / 2)
	backgroundWidth := size
	backgroundHeight := size

	red, green, blue := hsluv.HsluvToRGB(hue, 60, 60)
	foregroundColor, backgroundColor :=
		image.NewUniform(color.RGBA{255, 255, 255, 255}),
		image.NewUniform(color.RGBA{R: uint8(red * 255), G: uint8(green * 255), B: uint8(blue * 255), A: 255})

	background := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))
	draw.Draw(background, background.Bounds(), backgroundColor, image.ZP, draw.Src)
	face := truetype.NewFace(font, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: xfont.HintingFull,
	})
	fontDrawer := xfont.Drawer{Dst: background, Src: foregroundColor, Face: face}

	bounds, _ := fontDrawer.BoundString(text)
	width := bounds.Max.X - bounds.Min.X
	height := face.Metrics().Ascent - face.Metrics().Descent
	fontDrawer.Dot = fixed.P((backgroundWidth-width.Ceil())/2, (backgroundHeight-height.Ceil())/2+height.Ceil())

	fontDrawer.DrawString(text)
	err := png.Encode(buf, background)
	return buf, err
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 80
	} else {
		return h.Sum32()
	}
}

func parseParams(name string, size string) (string, int, float64) {
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 150
	}
	sizeInt = int(math.Min(float64(sizeInt), 1024))
	hue := hash(name) % 360
	name = strings.ToUpper(name)
	matched, _ := regexp.MatchString(`^\p{Han}`, name)

	var letters string
	if matched {
		letters = string([]rune(name)[0:1])
	} else {
		letters = name[0:2]
	}

	return letters, sizeInt, float64(hue)
}

func init() {

	f, err := static.LocalFile("./ttf", true).Open("WenQuanYi-Zen-Hei.ttf")
	if err != nil {
		panic(err)
	}
	fontBytes, _ = ioutil.ReadAll(f)
	font, _ = truetype.Parse(fontBytes)

}
func main() {
	router := gin.Default()
	router.GET("/image", func(ctx *gin.Context) {
		name := ctx.DefaultQuery("name", "IA")
		sizeStr := ctx.DefaultQuery("size", "150")
		letters, size, hue := parseParams(name, sizeStr)
		buffer, _ := generateImage(letters, size, hue)
		extraHeaders := map[string]string{
			"Content-Disposition": `inline; filename="avatar.png"`,
		}
		ctx.DataFromReader(200, int64(buffer.Len()), "image/png", buffer, extraHeaders)

	})
	router.Run(":8080")
}
