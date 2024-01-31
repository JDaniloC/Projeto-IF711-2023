package runner

import (
	"testing"
)

func BenchmarkCrawler(b *testing.B) {
	const link string = "https://hackerspaces.org/"
	const depth int = 3

	for i := 0; i < b.N; i++ {
		TimeoutCrawl(link, depth)
	}
}
