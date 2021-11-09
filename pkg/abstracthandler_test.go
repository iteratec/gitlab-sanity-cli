package gitlabsanitycli

import (
	"sort"
	"sync"
	"testing"

	"github.com/xanzy/go-gitlab"
)

// InMemoryPrinter implementation of Printer
type InMemoryPrinter struct {
	sync.RWMutex
	rows []content
}

// Print get results from egress channels and store them in a list
func (p *InMemoryPrinter) Print(rc <-chan content) {
	go func() {
		for r := range rc {
			p.Lock()
			p.rows = append(p.rows, r)
			p.Unlock()
		}
	}()
}

// GetRowsSortedById returns all received rows sorted by id
func (p *InMemoryPrinter) GetRowsSortedById() []content {
	p.Lock()
	sort.Slice(p.rows, func(i, j int) bool {
		return p.rows[i].ID < p.rows[j].ID
	})
	p.Unlock()
	return p.rows
}

// TestHandler implementation of AbstractHandler
type TestHandler struct {
	Printer Printer
}

func (u TestHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()
	for row := range rc {
		wc <- content{
			ID: row,
		}
	}
}

func (u TestHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	resp := &gitlab.Response{}
	xIds := []int{page}
	return xIds, resp, nil
}

func (u TestHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	genericHandler(u, channelHandlerFunc, 10)
}

func (u *TestHandler) Execute() {
	Run(u, u.Printer.Print)
}

// TestHandlerExecute verifies handler execution
func TestHandlerExecute(t *testing.T) {
	printer := &InMemoryPrinter{}
	testHandler := TestHandler{Printer: printer}
	testHandler.Execute()
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(want) != len(printer.GetRowsSortedById()) {
		t.Fatalf(`row count does not match. (actual=%v, want=%v)`, len(printer.GetRowsSortedById()), len(want))
	}
	for idx, actual := range printer.GetRowsSortedById() {
		if want[idx] != actual.ID {
			t.Fatalf(`mismatch at idx=%v: actual=%v, want=%v`, idx, actual.ID, want[idx])
		}
	}
}
