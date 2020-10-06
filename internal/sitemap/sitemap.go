package sitemap

import (
	"net/url"
	"path/filepath"
	"strings"
)

type Page struct {
	Assets []string `json:"assets"`
	Children map[string][]Page `json:"children"`
}

type Map struct {
	Hosts map[string]Page `json:"hosts"`
}

func NewMap() *Map {
	return &Map{
		Hosts: map[string]Page{},
	}
}

func (m *Map) Add(u *url.URL, links []*url.URL) {
	// Dont add assets to the sitemap individually
	if m.isAsset(u) {
		return
	}

	host := u.Hostname()
	if _, ok := m.Hosts[host]; !ok {
		m.Hosts[host] = Page{
			Assets: []string{},
			Children: map[string][]Page{},
		}
	}

	p := m.Hosts[host]
	page := &p
	segments := strings.Split(u.Path, "/")

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		prettySegment := "/" + segment

		if _, ok := page.Children[prettySegment]; !ok {
			np := Page{
				Assets: []string{},
				Children: map[string][]Page{},
			}

			page.Children[prettySegment] = []Page{np}
			page = &(page.Children[prettySegment][0])
		}
	}

	main:
	for _, l := range links {
		if !m.isAsset(l) {
			continue
		}

		for _, a := range page.Assets {
			if l.String() == a {
				continue main
			}
		}

		page.Assets = append(page.Assets, l.String())
	}

	m.Hosts[host] = p
}

// isAsset checks to see if a URL is an asset by simply seeing if it has some form of extension.
// This can easily be improved by checking mime type of URL data.
func (m *Map) isAsset(u *url.URL) bool {
	return filepath.Ext(u.Path) != ""
}