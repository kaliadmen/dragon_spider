package dragonSpider

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"time"
)

func (ds *DragonSpider) TakeScreenshot(pageUrl, testName string, width, height float64) {
	page := rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageUrl).MustWaitLoad()

	screenshot, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: nil,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
			Scale:  1,
		},
		FromSurface:           true,
		CaptureBeyondViewport: false,
	})
	if err != nil {
		ds.ErrorLog.Println(err)
	}

	fileName := time.Now().Format("2006-01-02-15-04-05.000000")

	err = utils.OutputFile(fmt.Sprintf("%s/screenshots/%s-%s.png", ds.RootPath, testName, fileName), screenshot)
	if err != nil {
		ds.ErrorLog.Println(err)
	}

}

func (ds *DragonSpider) FetchPage(pageUrl string) *rod.Page {
	return rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(pageUrl).MustWaitLoad()
}

func (ds *DragonSpider) SelectElementById(page *rod.Page, id string) *rod.Element {
	return page.MustElement(fmt.Sprintf("#%s", id))
}
