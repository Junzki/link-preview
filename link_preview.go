package link_preview

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

type LinkPreview struct {
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	ImageURL	string	`json:"image"`
	Link		string	`json:"website"`

	parsed      *goquery.Document
	client      *http.Request
}

type HTMLMetaAttr struct {
	Key		string
	Value   string
}

func PreviewLink(link string, extraClient *http.Request) (*LinkPreview, error) {
	preview := LinkPreview{}

	preview.Link = link

	if nil == extraClient {
		preview.initClient()
	} else {
		preview.client = extraClient
	}
	err := preview.preview()

	if err != nil {
		return nil, err
	}
	return &preview, nil
}


func (p *LinkPreview) initClient() {
	client, _ := http.NewRequest("GET", p.Link, nil)
	p.client = client
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
	res, err := http.DefaultClient.Do(p.client)
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

	// Find `favicon.ico`.
	linkNodes := p.parsed.Find("html > head > link")
	for _, node := range linkNodes.Nodes {
		for _, attr := range node.Attr {
			switch strings.ToLower(attr.Key) {
			case "rel":
				if attr.Val != "icon" {
					break
				}
				p.parseFavicon(node)
			default:
				continue
			}
		}
	}

	// Parse <meta> tags.
	metaNodes := p.parsed.Find("html > head > meta")
	for _, node := range metaNodes.Nodes {
		for _, attr := range node.Attr {
			switch strings.ToLower(attr.Key) {
			case "property":
				p.parseMetaProperties(attr.Val, node)
				break
			case "name":
				if "description" == strings.ToLower(attr.Val) && "" == p.Description {
					content := p.parseMetaContent(node)
					p.Description = content
					break
				}
			default:
				continue
			}
		}
	}

	return nil
}


func (p *LinkPreview) parseMetaContent(node *html.Node) string {
	var content string
	for _, attr := range node.Attr {
		switch strings.ToLower(attr.Key) {
		case "content":
			content = attr.Val
			break
		default:
			continue
		}
	}

	return content
}


func (p *LinkPreview) parseMetaProperties(nodeType string, node *html.Node) {
	nodeType = strings.ToLower(nodeType)
	if ! strings.HasPrefix(nodeType, "og:") {
		return
	}

	slices := strings.Split(nodeType, ":")
	if 2 != len(slices) {
		return
	}

	nodeType = slices[1]
	content := p.parseMetaContent(node)

	switch nodeType {
	case "description":
		p.Description = content
	case "image":
		p.ImageURL = content
	case "title":
		p.Title = content
	}
}

func (p *LinkPreview) parseFavicon(node *html.Node) {
	var link string

	for _, attr := range node.Attr {
		switch strings.ToLower(attr.Key) {
		case "href":
			link = attr.Val
			break
		default:
			continue
		}
	}

	if "" == link {
		return
	}

	if strings.HasPrefix("http://", link) || strings.HasPrefix("https://", link) {
		p.ImageURL = link
		return
	}

	parsedURL, _ := url.Parse(p.Link)
	joinedURL := url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   link,
	}

	link = joinedURL.String()
	if "" == p.ImageURL {
		p.ImageURL = link
	}
}

func (p *LinkPreview) ToJSON() string {
	serialized, err := json.Marshal(p)
	if err != nil {
		return `{
			"title": "",
			"description": "",
			"image": "",
			"link": ""
		}`
	}
	return string(serialized)
}