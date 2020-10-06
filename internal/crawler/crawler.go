package crawler

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"webcrawler/internal/queue"
	"webcrawler/internal/sitemap"
)
import "golang.org/x/net/html"

type Crawler struct {
	u *url.URL
	workers int
	client *http.Client
}

type Resource struct {
	Segment string
	Children []Resource
}

// New will return a new instance of Crawler or an error
// if given url could not be parsed.
func New(u string, workers int) (*Crawler, error) {
	// We will need at least one worker...
	if workers <= 0 {
		workers = 1
	}

	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		u: url,
		workers: workers,
		client: http.DefaultClient,
	}, nil
}

// Crawl will generate a sitemap by crawling the initial URL.
// A good improvement would be to return any generated errors
// by the scanner. We will log them for now.
func (c *Crawler) Crawl(ctx context.Context) *sitemap.Map {
	sm := sitemap.NewMap()

	queue := queue.NewQueue()
	queue.Add(c.u)

	for i := 0; i < c.workers; i++ {
		go func() {
			for {
				var (
					ok bool
					pu *url.URL
				)

				select {
				case <-ctx.Done():
					return
				case pu, ok = <- queue.URL():
					if !ok {
						return
					}
				}

				// Wrap this part in a function so we can defer the call to queue.Done
				// instead of duplicating for every exit scenario.
				func(pu *url.URL) {
					defer queue.Done()
					newUrls, err := c.scan(pu)
					if err != nil {
						log.Print(errors.Wrap(err, "could not scan url"))
						return
					}

					sm.Add(pu, newUrls)
					queue.AddBulk(newUrls)
				}(pu)
			}
		}()
	}

	// This part may seem complicated but its basically
	// creating a channel for waiting for the queue
	// so we can select which comes first. Context closing
	// or the site crawl finishing.
	select {
	case <-ctx.Done():
		return sm
	case <-func() chan struct{} {
		ch := make(chan struct{})
		go func() {
			queue.Wait()
			queue.Close()
			close(ch)
		}()
		return ch
	}():
	}

	return sm
}

// scan will try to download a given URL and parse it as a HTML page.
// It will then return a list of found urls within the HTML or return an error.
func (c *Crawler) scan(u *url.URL) ([]*url.URL, error) {
	var urls []*url.URL
	allowed := map[string]struct{}{
		"href": {},
		"src": {},
	}

	resp, err := c.client.Get(u.String())
	if err != nil {
		return urls, errors.Wrap(err, "could not get url")
	}

	tree, err := html.Parse(resp.Body)
	if err != nil {
		return urls, errors.Wrap(err, "could not parse response body")
	}

	log.Printf("scanning: %s\n", u)

	var parse func(*html.Node)
	parse = func(node *html.Node) {

		// Could easily improve by only parsing allowed tag types
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			for _, x := range child.Attr {

				// If the attribute is not a supported type, try the next one
				if _, ok := allowed[x.Key]; !ok {
					continue
				}

				// If the url is invalid, ignore and try the next attribute
				pu, err := url.Parse(strings.TrimSpace(x.Val))
				if err != nil {
					continue
				}

				// If the URL is relative, and doesnt contain a #fragment, then we want it.
				// We will ignore fragments for this basic crawler.
				if !pu.IsAbs() && pu.Fragment == "" {
					ref := c.u.ResolveReference(pu)
					if ref != nil {
						pu = ref
					}
				}

				// Double check that the host is part of the same domain as the primary URL. Subdomains supported.
				if pu.IsAbs() && (pu.Host == c.u.Host || strings.HasSuffix(pu.Host, "." + c.u.Host)) {
					urls = append(urls, pu)
				}
			}

			parse(child)
		}
	}

	parse(tree)

	return urls, nil
}
