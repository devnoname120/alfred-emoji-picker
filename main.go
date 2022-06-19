package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/deanishe/awgo"
	"github.com/hackebrot/turtle"
	"github.com/samber/lo"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	query := os.Args[1]
	results := search(query)

	for _, result := range results {
		wf.NewItem(result.Name).
			Subtitle(fmt.Sprintf("Input \"%s\" (%s) into foremost application", result.Char, result.Name)).
			Arg(result.Char).
			Icon(&aw.Icon{Value: fmt.Sprintf("emojis/%s.png", result.Name)}).
			Valid(true)
	}

	wf.WarnEmpty("No matching emojis", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

func search(query string) []*turtle.Emoji {
	if query == "" {
		return make([]*turtle.Emoji, 0)
	}

	nameExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name == query
	})

	nameMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name != query && strings.Contains(e.Name, query)
	})

	keywordExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword == query {
				return true
			}
		}
		return false
	})

	keywordMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword != query && strings.Contains(keyword, query) {
				return true
			}
		}
		return false
	})

	categoryExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category == query
	})

	categoryMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category != query && strings.Contains(e.Category, query)
	})

	results := [][]*turtle.Emoji{
		nameExactMatches,
		nameMatches,
		keywordExactMatches,
		keywordMatches,
		categoryExactMatches,
		categoryMatches,
	}

	consolidated := lo.Flatten(results)
	return consolidated
}

//func textToImage() {
//	//bgImage, _ := gg.LoadImage(request.BgImgPath)
//	imgWidth := 64
//	imgHeight := 64
//
//	dc := gg.NewContext(imgWidth, imgHeight)
//	//dc.DrawImage(bgImage, 0, 0)
//
//	//_ := dc.LoadFontFace(request.FontPath, request.FontSize)
//
//	//x := float64(imgWidth / 2)
//	//y := float64((imgHeight / 2) - 80)
//	//maxWidth := float64(imgWidth) - 60.0
//	dc.SetColor(color.Black)
//	dc.DrawString("lol", 0, 64)
//
//	img := dc.Image()
//
//	f, _ := os.OpenFile("myimage.png", os.O_RDWR, 0777)
//	_ = png.Encode(f, img)
//}

//type PageController struct {
//	fontFile string
//}
//
//func (z *PageController) loadFont() (*truetype.Font, error) {
//	z.fontFile = "static/fonts/UbuntuMono-R.ttf"
//	fontBytes, err := ioutil.ReadFile(z.fontFile)
//	if err != nil {
//		return nil, err
//	}
//	f, err := freetype.ParseFont(fontBytes)
//	if err != nil {
//		return nil, err
//	}
//	return f, nil
//}
//
//func (z *PageController) generateImage(textContent string, fgColorHex string, bgColorHex string, fontSize float64) ([]byte, error) {
//
//	fgColor := color.RGBA{0xff, 0xff, 0xff, 0xff}
//	if len(fgColorHex) == 7 {
//		_, err := fmt.Sscanf(fgColorHex, "#%02x%02x%02x", &fgColor.R, &fgColor.G, &fgColor.B)
//		if err != nil {
//			log.Println(err)
//			fgColor = color.RGBA{0x2e, 0x34, 0x36, 0xff}
//		}
//	}
//
//	bgColor := color.RGBA{0x30, 0x0a, 0x24, 0xff}
//	if len(bgColorHex) == 7 {
//		_, err := fmt.Sscanf(bgColorHex, "#%02x%02x%02x", &bgColor.R, &bgColor.G, &bgColor.B)
//		if err != nil {
//			log.Println(err)
//			bgColor = color.RGBA{0x30, 0x0a, 0x24, 0xff}
//		}
//	}
//
//	loadedFont, err := z.loadFont()
//	if err != nil {
//		return nil, err
//	}
//
//	code := strings.Replace(textContent, "\t", "    ", -1) // convert tabs into spaces
//	text := strings.Split(code, "\n")                      // split newlines into arrays
//
//	fg := image.NewUniform(fgColor)
//	bg := image.NewUniform(bgColor)
//	rgba := image.NewRGBA(image.Rect(0, 0, 1200, 630))
//	draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src)
//	c := freetype.NewContext()
//	c.SetDPI(72)
//	c.SetFont(loadedFont)
//	c.SetFontSize(fontSize)
//	c.SetClip(rgba.Bounds())
//	c.SetDst(rgba)
//	c.SetSrc(fg)
//	c.SetHinting(font.HintingNone)
//
//	textXOffset := 50
//	textYOffset := 10 + int(c.PointToFixed(fontSize)>>6) // Note shift/truncate 6 bits first
//
//	pt := freetype.Pt(textXOffset, textYOffset)
//	for _, s := range text {
//		_, err = c.DrawString(strings.Replace(s, "\r", "", -1), pt)
//		if err != nil {
//			return nil, err
//		}
//		pt.Y += c.PointToFixed(fontSize * 1.5)
//	}
//
//	b := new(bytes.Buffer)
//	if err := png.Encode(b, rgba); err != nil {
//		log.Println("unable to encode image.")
//		return nil, err
//	}
//	return b.Bytes(), nil
//}
//
//func textToImage() {
//	//fgColor := color.RGBA{0xff, 0xff, 0xff, 0xff}
//	//if len(fgColorHex) == 7 {
//	//	_, err := fmt.Sscanf(fgColorHex, "#%02x%02x%02x", &fgColor.R, &fgColor.G, &fgColor.B)
//	//	if err != nil {
//	//		log.Println(err)
//	//		fgColor = color.RGBA{0x2e, 0x34, 0x36, 0xff}
//	//	}
//	//}
//
//	//bgColor := color.RGBA{0x30, 0x0a, 0x24, 0xff}
//	//if len(bgColorHex) == 7 {
//	//	_, err := fmt.Sscanf(bgColorHex, "#%02x%02x%02x", &bgColor.R, &bgColor.G, &bgColor.B)
//	//	if err != nil {
//	//		log.Println(err)
//	//		bgColor = color.RGBA{0x30, 0x0a, 0x24, 0xff}
//	//	}
//	//}
//
//	//fg := image.NewUniform(fgColor)
//	//bg := image.NewUniform(bgColor)
//
//	loadedFont, _ := z.loadFont()
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	rgba := image.NewRGBA(image.Rect(0, 0, 1200, 630))
//	//draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src)
//	c := freetype.NewContext()
//	c.SetDPI(72)
//	c.SetFont(loadedFont)
//	c.SetFontSize(32)
//	c.SetClip(rgba.Bounds())
//	c.SetDst(rgba)
//	//c.SetSrc(fg)
//	c.SetHinting(font.HintingNone)
//
//	textXOffset := 50
//	textYOffset := 10 + int(c.PointToFixed(32)>>6) // Note shift/truncate 6 bits first
//
//	pt := freetype.Pt(textXOffset, textYOffset)
//	stuff, _ := c.DrawString("yeah", pt)
//
//	buffer := new(bytes.Buffer)
//	_ = png.Encode(buffer, stuff)
//}
