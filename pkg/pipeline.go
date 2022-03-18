package gitlabsanitycli

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/xanzy/go-gitlab"
)

// PipelineHandler implements AbstractHandler
// implementation for interacting with Gitlab Project Pipelines
type PipelineHandler struct {
	Git         *gitlab.Client
	Printer     Printer
	Config      *Config
	ListOptions gitlab.ListProjectPipelinesOptions
}

func (p PipelineHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()
	now := time.Now()
	for id := range rc {
		pipeline, _, err := p.Git.Pipelines.GetPipeline(p.Config.ResourceId, id)
		if err != nil {
			log.Fatal(err)
		}
		if pipeline.FinishedAt == nil {
			log.Printf("Unable to fetch data for Pipeline ID: %v, URL: %v", pipeline.ID, pipeline.WebURL)
		} else {
			// Age filter
			if now.Sub(*pipeline.FinishedAt).Hours() > float64(p.Config.Age) {
				wc <- content{
					Req:          Pipeline,
					ID:           pipeline.ID,
					Status:       pipeline.Status,
					URL:          pipeline.WebURL,
					ProjectID:    pipeline.ProjectID,
					LastActivity: math.Round(now.Sub(*pipeline.FinishedAt).Hours()),
				}
			}
		}
	}
}

func (p PipelineHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	p.ListOptions.Page = page

	values, resp, err := p.Git.Pipelines.ListProjectPipelines(p.Config.ResourceId, &p.ListOptions)

	var xIds []int
	for _, value := range values {
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}

func (p PipelineHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	p.ListOptions = gitlab.ListProjectPipelinesOptions{
		Sort:        gitlab.String("desc"),
		ListOptions: gitlab.ListOptions{PerPage: PageSize},
	}

	_, resp, err := p.Git.Pipelines.ListProjectPipelines(p.Config.ResourceId, &p.ListOptions)

	if err != nil {
		log.Fatal(err)
	}

	genericHandler(p, channelHandlerFunc, resp.TotalPages)

}

func (p PipelineHandler) deletePipeline(pipelineId int, dryRun bool) {
	project, _, err0 := p.Git.Projects.GetProject(p.Config.ResourceId, nil)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.Printf("Deleting Pipeline %v in Project: %v\n", pipelineId, project.Name)
	if !dryRun {
		_, err1 := p.Git.Pipelines.DeletePipeline(project.ID, pipelineId)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("Pipeline %v deleted \n", pipelineId)
	} else {
		log.Printf("[DryRun] Pipeline %v was not deleted\n", pipelineId)
	}
}

func (p PipelineHandler) deleteAllPipeline(Pipeline <-chan content) {
	go func() {
		for pipeline := range Pipeline {
			p.deletePipeline(pipeline.ID, p.Config.DryRun)
		}
	}()
}

func (p PipelineHandler) List() {
	Run(p, p.Printer.Print)
}

func (p PipelineHandler) Delete() {
	p.deletePipeline(p.Config.ResourceId, p.Config.DryRun)
}

func (p PipelineHandler) DeleteAll() {
	Run(p, p.deleteAllPipeline)
}
