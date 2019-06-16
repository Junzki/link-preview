package handlers

import (
	"errors"
	"golang.org/x/net/html"
	"strings"
)

type StandardLinkPreview struct {
	*LinkPreviewContext
}

func (p *StandardLinkPreview) Preview() error {
	err := p.request()

	if nil != err {
		return err
	}

	err = p.readTags()
	return err
}

func (p *StandardLinkPreview) readTags() error {
	titleNode := p.Parsed.Find("html > head > title")
	if 0 == titleNode.Length() {
		return errors.New("title not found")
	}
	p.Title = titleNode.Text()

	// Find `favicon.ico`.
	linkNodes := p.Parsed.Find("html > head > link")
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
	metaNodes := p.Parsed.Find("html > head > meta")
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

func (p *StandardLinkPreview) parseMetaContent(node *html.Node) string {
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

func (p *StandardLinkPreview) parseMetaProperties(nodeType string, node *html.Node) {
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
