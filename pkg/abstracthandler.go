package gitlabsanitycli

import (
	"log"
	"sync"
	"time"

	"github.com/xanzy/go-gitlab"
)

/*	AbstractHandler Interface
	Pipe() implements the Pull and Push logic to grab Data from API and Push back to output Channel
	Controller() implements the logic to start and run Pipe() and takes the ChannelHandlerFunc (example logic: see genericHandler() func)
	ApiCallFunc() takes the current page and respones with the ID slice, respose and error
*/
type AbstractHandler interface {
	Pipe(*sync.WaitGroup, chan<- content, <-chan int)
	Controller(ChannelHandlerFunc)
	ApiCallFunc(int) ([]int, *gitlab.Response, error)
}

// ChannelHandlerFunc gets results from Pipe
type ChannelHandlerFunc func(<-chan content)

// Run is the entry func() of AbstractHandler to call the Controller() implementation
func Run(a AbstractHandler, channelHandlerFunc ChannelHandlerFunc) {
	a.Controller(channelHandlerFunc)
}

func genericHandler(a AbstractHandler, channelHandler ChannelHandlerFunc, totalPages int) {
	// Use Waitgroup to run concurrent Pipe() functions
	var wg sync.WaitGroup
	defer wg.Wait()

	// egressChan = Report/Back content
	egressChan := make(chan content)
	// ingressChan = Data Information content (the APICall results)
	ingressChan := make(chan int)

	// Call channelHandler
	channelHandler(egressChan)

	log.Printf("Initialize %d Pipes\n", ConcurrentAPIRequest)
	// Start concurrent Pipe's
	for i := 0; i < ConcurrentAPIRequest; i++ {
		wg.Add(1)
		go a.Pipe(&wg, egressChan, ingressChan)
	}

	for i := 0; i < totalPages; i++ {
		ids, _, err := a.ApiCallFunc(i)
		if err != nil {
			log.Fatal(err)
		}

		for _, value := range ids {
			ingressChan <- value
		}
	}

	// Wait seconds before closing ingress channel and end API call processing
	// Without the second wait, the channel is closing too fast before last results are reached egress channel
	time.Sleep(time.Second * 2)

	close(ingressChan)
}
