package runner

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	c "github.com/JDaniloC/Projeto-IF711-2023/internal/utils"
	soup "github.com/anaskhan96/soup"
)

const (
	defaultTimeout = 10 * time.Second
)

func crawl(wg *sync.WaitGroup, control *c.Controller, url string, depth int) {
	defer wg.Done()

	if !control.CanVisit(url, depth) {
		return
	}

	control.AddVisitedLink(url)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		control.AddInvalidLink(url)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		control.AddInvalidLink(url)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		control.AddInvalidLink(url)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		control.AddInvalidLink(url)
		return
	}

	html := soup.HTMLParse(string(body))
	links := html.FindAll("a")
	control.AddValidLink(url)

	for _, link := range links {
		href := link.Attrs()["href"]
		if strings.HasPrefix(href, "http") {
			wg.Add(1)
			go crawl(wg, control, href, depth+1)
		}
	}
}

func TimeoutCrawl(control *c.Controller, url string) {
	var wg sync.WaitGroup

	wg.Add(1)
	crawl(&wg, control, url, 0)
	wg.Wait()
}
