package runner

import (
	"testing"

	c "github.com/JDaniloC/Projeto-IF711-2023/internal/utils"
)

func BenchmarkCrawler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const link string = "https://hackerspaces.org/"
		const depth int = 2

		controller := c.NewController(depth)
		TimeoutCrawl(controller, link)
	}
}
