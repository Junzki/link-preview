package handlers

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

const (
	StandardMetaTags = iota
	WeChatMP
)

type HTMLMetaAttr struct {
	Key   string
	Value string
}

type LinkPreviewContext struct {
	TargetType  int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image"`
	Link        string `json:"website"`

	Client *http.Request     `json:"-"`
	Parsed *goquery.Document `json:"-"`
}

func (p *LinkPreviewContext) initClient() {
	client, _ := http.NewRequest("GET", p.Link, nil)
	p.Client = client
}

func (p *LinkPreviewContext) request() error {
	res, err := http.DefaultClient.Do(p.Client)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if nil != err {
		return err
	}

	p.Parsed = doc
	return nil
}

func (p *LinkPreviewContext) parseFavicon(node *html.Node) {
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

type PreviewHandler interface {
	Preview() error
}

func GetPreviewHandler(c *LinkPreviewContext) (PreviewHandler, error) {
	if nil == c {
		return nil, errors.New("bad link preview cxt, nil given")
	}

	if nil == c.Client {
		c.initClient()
	}

	var handler PreviewHandler

	switch c.TargetType {
	case StandardMetaTags:
		handler = &StandardLinkPreview{
			c,
		}
	default:
		return nil, errors.New("unknown target type")
	}

	return handler, nil
}
