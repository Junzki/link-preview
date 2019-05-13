package link_preview

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"image"
	"net/http"
)

type LinkPreview struct {
	Title		string
	Description	string
	Image       *image.Image
	ImageURL	string

	Link		string

	parsed      *goquery.Document
}

type HTMLMetaAttr struct {
	Key		string
	Value   string
}

func PreviewLink(link string) (*LinkPreview, error) {
	preview := LinkPreview{}

	preview.Link = link
	err := preview.preview()

	if err != nil {
		return nil, err
	}
	return &preview, nil
}


func (p *LinkPreview) preview() error {
	err := p.request()
	if nil != err {
		return err
	}

	err = p.readTags()
	return err
}


func (p *LinkPreview) request() error {
	res, err := http.Get(p.Link)
	if nil != err {
		return err
	}
	defer res.Body.Close()


	doc, err := goquery.NewDocumentFromReader(res.Body)
	if nil != err {
		return err
	}

	p.parsed = doc
	return nil
}


func (p *LinkPreview) readTags() error {
	titleNode := p.parsed.Find("html > head > title")
	if 0 == titleNode.Length() {
		return errors.New("title not found")
	}
	p.Title = titleNode.Text()

	metaNodes := p.parsed.Find("html > head > meta")
	if 0 == metaNodes.Length() {
		return errors.New("meta not found")
	}


	return nil
}
