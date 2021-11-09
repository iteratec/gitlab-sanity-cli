package gitlabsanitycli

import (
	"log"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

// UserHandler implements AbstractHandler
// implementation for interacting with Gitlab User Resources
type UserHandler struct {
	Git         *gitlab.Client
	Printer     Printer
	Config      *Config
	listOptions gitlab.ListUsersOptions
}

func (u UserHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()
	for id := range rc {

		user, _, err := u.Git.Users.GetUser(id, gitlab.GetUsersOptions{})
		if err != nil {
			log.Fatal(err)
		}

		wc <- content{
			Req:  User,
			ID:   user.ID,
			Name: user.Name,
		}

	}
}

func (u UserHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	u.listOptions.Page = page

	values, resp, err := u.Git.Users.ListUsers(&u.listOptions)

	// Copy gitlab.user slice into int slice
	var xIds []int
	for _, value := range values {
		if len(u.Config.Query) > 0 && !strings.Contains(strings.ToLower(value.Name), strings.ToLower(u.Config.Query)) {
			continue
		}
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}

func (u UserHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	u.listOptions = gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{PerPage: PageSize},
	}

	if len(u.Config.Query) > 0 {
		u.listOptions.Search = &u.Config.Query
	}
	_, resp, err := u.Git.Users.ListUsers(&u.listOptions)

	if err != nil {
		log.Fatal(err)
	}

	genericHandler(u, channelHandlerFunc, resp.TotalPages)
}

func (u *UserHandler) List() {
	Run(u, u.Printer.Print)
}
