package stats

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type ProgressBar struct {
	total   int
	current int
	mutex   sync.Mutex
}

func (pb *ProgressBar) Increment() {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.current++
	pb.render()
}

func (pb *ProgressBar) render() {
	percentage := float64(pb.current) / float64(pb.total) * 100
	barWidth := 50
	filled := int(percentage / 100 * float64(barWidth))
	empty := barWidth - filled

	fmt.Fprintf(os.Stderr, "\r[%s%s] %d%% (%d/%d)\n",
		strings.Repeat("#", filled),
		strings.Repeat("-", empty),
		int(percentage), pb.current, pb.total)
}
