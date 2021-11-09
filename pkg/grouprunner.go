
package gitlabsanitycli

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

// GroupRunnerHandler implements AbstractHandler
// implementation for interacting with Gitlab Grouprunner Resources
type GroupRunnerHandler struct {
	Git         *gitlab.Client
	Printer     Printer
	Config      *Config
	ListOptions gitlab.ListGroupsRunnersOptions
	GroupID     int
}

func (g GroupRunnerHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()

	for id := range rc {
		runnerDetails, _, err := g.Git.Runners.GetRunnerDetails(id, nil)
		if err != nil {
			log.Fatal(err)
		}
		wc <- content{
			Req:         GroupRunner,
			ID:          runnerDetails.ID,
			Name:        runnerDetails.Name,
			Description: runnerDetails.Description,
			Status:      runnerDetails.Status,
			IPAddress:   runnerDetails.IPAddress,
			Active:      runnerDetails.Active,
			IsShared:    runnerDetails.IsShared,
			Online:      runnerDetails.Online,
		}

	}
}

func (g GroupRunnerHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	g.ListOptions.Page = page

	values, resp, err := g.Git.Runners.ListGroupsRunners(g.GroupID, &g.ListOptions)

	var xIds []int
	for _, value := range values {
		if (len(g.Config.Query) > 0 && !strings.Contains(strings.ToLower(value.Description), strings.ToLower(g.Config.Query))) ||
			(len(g.Config.Status) > 0 && !strings.Contains(strings.ToLower(value.Status), strings.ToLower(g.Config.Status))) {
			continue
		}
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}

func (g GroupRunnerHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	listGroupsOptions := gitlab.ListGroupsOptions{
		TopLevelOnly: gitlab.Bool(true),
	}
	g.ListOptions = gitlab.ListGroupsRunnersOptions{
		ListOptions: gitlab.ListOptions{PerPage: PageSize},
	}

	gids, _, err := g.Git.Groups.ListGroups(&listGroupsOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range gids {
		_, resp, err := g.Git.Runners.ListGroupsRunners(v.ID, &g.ListOptions)
		if err != nil {
			if resp.StatusCode == 403 {
				continue
			} else {
				fmt.Println(err)
			}
		}
		log.Printf("Got %v Runner in Group %v from API\n", resp.TotalItems, v.ID)
		g.GroupID = v.ID
		genericHandler(g, channelHandlerFunc, resp.TotalPages)
	}
}

func (g *GroupRunnerHandler) deleteGroupRunner(groupRunnerId int, dryRun bool) {
	// dry-run einbauen
	groupRunner, resp, err := g.Git.Runners.GetRunnerDetails(groupRunnerId, nil)
	if resp.StatusCode == 404 {
		log.Fatalf("GroupRunner with ID: %v was not found", groupRunnerId)
	}
	if err != nil {
		log.Fatal(err)
	}
	if !dryRun {
		resp, err := g.Git.Runners.RemoveRunner(groupRunnerId)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != gitlabReponseCodeOk {
			log.Fatalf("Failed to remove GroupRunner %v [ID: %v]: Response: %v, Code: %v\n", groupRunner.Description, groupRunner.ID, resp.Response, resp.StatusCode)
		}
		log.Printf("GroupRunner %v [ID: %v] removed\n", groupRunner.Description, groupRunner.ID)
	} else {
		log.Printf("[DryRun] GroupRunner %v [ID: %v] was not removed\n", groupRunner.Description, groupRunnerId)
	}
}

func (g *GroupRunnerHandler) deleteAllGroupRunner(recvChan <-chan content) {
	go func() {
		for groupRunner := range recvChan {
			g.deleteGroupRunner(groupRunner.ID, g.Config.DryRun)
		}
	}()
}

func (g *GroupRunnerHandler) List() {
	Run(g, g.Printer.Print)
}

func (g *GroupRunnerHandler) Delete() {
	g.deleteGroupRunner(g.Config.ResourceId, g.Config.DryRun)
}

func (g *GroupRunnerHandler) DeleteAll() {
	Run(g, g.deleteAllGroupRunner)
}
