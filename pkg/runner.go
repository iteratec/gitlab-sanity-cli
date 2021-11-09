package gitlabsanitycli

import (
	"log"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

// RunnerHandler implements AbstractHandler
// implementation for interacting with Gitlab Runner Resources
type RunnerHandler struct {
	Git         *gitlab.Client
	Printer     Printer
	Config      *Config
	ListOptions gitlab.ListRunnersOptions
}

func (r RunnerHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()

	for id := range rc {
		runnerDetails, _, err := r.Git.Runners.GetRunnerDetails(id, nil)
		if err != nil {
			log.Fatal(err)
		}

		wc <- content{
			Req:         Runner,
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

func (r RunnerHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	r.ListOptions.Page = page

	values, resp, err := r.Git.Runners.ListRunners(&r.ListOptions)

	var xIds []int
	for _, value := range values {
		if (len(r.Config.Query) > 0 && !strings.Contains(strings.ToLower(value.Description), strings.ToLower(r.Config.Query))) ||
			(len(r.Config.Status) > 0 && !strings.Contains(strings.ToLower(value.Status), strings.ToLower(r.Config.Status))) {
			continue
		}
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}

func (r RunnerHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	r.ListOptions = gitlab.ListRunnersOptions{
		ListOptions: gitlab.ListOptions{PerPage: PageSize},
	}

	_, resp, err := r.Git.Runners.ListRunners(&r.ListOptions)
	if err != nil {
		log.Fatal(err)
	}

	genericHandler(r, channelHandlerFunc, resp.TotalPages)
}

func (r RunnerHandler) deleteAllRunner(runners <-chan content) {
	go func() {
		for runner := range runners {
			r.deleteRunner(runner.ID, r.Config.DryRun)
		}
	}()
}

func (r RunnerHandler) deleteRunner(runnerId int, dryRun bool) {
	runner, resp, err := r.Git.Runners.GetRunnerDetails(runnerId)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Runner with ID %v not found or accessable", runnerId)
	}
	if runner.Status == "Online" {
		log.Fatalf("Runner with ID %v is not offline", runnerId)
		return
	}
	if !dryRun {
		resp, err = r.Git.Runners.DeleteRegisteredRunnerByID(runnerId)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != gitlabReponseCodeOk {
			log.Fatalf("Failed to remove Runner %v [ID: %v]: Response: %v, Code: %v\n", runner.Description, runner.ID, resp.Response, resp.StatusCode)
		}
		log.Printf("Succesfully removed runner %v with ID: %v\n", runner.Name, runnerId)

	} else {
		log.Printf("[DryRun] Runner %v with ID: %v was not deleted\n", runner.Name, runnerId)
	}
}

func (r *RunnerHandler) List() {
	Run(r, r.Printer.Print)
}

func (r *RunnerHandler) Delete() {
	r.deleteRunner(r.Config.ResourceId, r.Config.DryRun)
}

func (r RunnerHandler) DeleteAll() {
	Run(r, r.deleteAllRunner)
}
