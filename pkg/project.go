package gitlabsanitycli

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/xanzy/go-gitlab"
)

// ProjectHandler implements AbstractHandler
// implementation for interacting with Gitlab Project Resources
type ProjectHandler struct {
	Git         *gitlab.Client
	Printer     Printer
	Config      *Config
	ListOptions gitlab.ListProjectsOptions
}

func (p ProjectHandler) Pipe(wg *sync.WaitGroup, wc chan<- content, rc <-chan int) {
	defer wg.Done()
	now := time.Now()
	for id := range rc {
		project, _, err := p.Git.Projects.GetProject(id, nil)
		if err != nil {
			log.Fatal(err)
		}
		// Age filter
		if now.Sub(*project.LastActivityAt).Hours() > float64(p.Config.Age) {
			wc <- content{
				Req:          Project,
				ID:           project.ID,
				Name:         project.Name,
				LastActivity: math.Round(now.Sub(*project.LastActivityAt).Hours()),
			}
		}
	}
}

func (p ProjectHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	p.ListOptions.Page = page

	values, resp, err := p.Git.Projects.ListProjects(&p.ListOptions)

	var xIds []int
	for _, value := range values {
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}

func (p ProjectHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	p.ListOptions = gitlab.ListProjectsOptions{
		SearchNamespaces: gitlab.Bool(false),
		Membership:       gitlab.Bool(false),
		Owned:            gitlab.Bool(false),
		Simple:           gitlab.Bool(false),
		Sort:             gitlab.String("asc"),
		Archived:         gitlab.Bool(false),
		ListOptions:      gitlab.ListOptions{PerPage: PageSize},
	}
	switch p.Config.ProjectType {
	case "private":
		p.ListOptions.Visibility = gitlab.Visibility(gitlab.PrivateVisibility)
	case "public":
		p.ListOptions.Visibility = gitlab.Visibility(gitlab.PublicVisibility)
	}

	if len(p.Config.Query) > 0 {
		p.ListOptions.Search = gitlab.String(p.Config.Query)
	}

	_, resp, err := p.Git.Projects.ListProjects(&p.ListOptions)

	if err != nil {
		log.Fatal(err)
	}
	genericHandler(p, channelHandlerFunc, resp.TotalPages)

}

func (p ProjectHandler) archiveProject(projectId int, dryRun bool) {
	project, _, err0 := p.Git.Projects.GetProject(projectId, nil)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.Printf("Archiving Project %v\n", project.Name)
	if !dryRun {
		_, _, err := p.Git.Projects.ArchiveProject(projectId)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Project %v archived\n", project.Name)
	} else {
		log.Printf("[DryRun] Project %v was not archived\n", project.Name)
	}
}

func (p ProjectHandler) archiveProjects(projects <-chan content) {
	go func() {
		for project := range projects {
			p.archiveProject(project.ID, p.Config.DryRun)
		}
	}()
}

func (p ProjectHandler) deleteProject(projectId int, dryRun bool) {
	project, _, err0 := p.Git.Projects.GetProject(projectId, nil)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.Printf("Deleting Project %v\n", project.Name)
	if !dryRun {
		_, err1 := p.Git.Projects.DeleteProject(projectId)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("Project %v deleted\n", project.Name)
	} else {
		log.Printf("[DryRun] Project %v was not deleted\n", project.Name)
	}
}

func (p ProjectHandler) deleteProjects(projects <-chan content) {
	go func() {
		for project := range projects {
			p.deleteProject(project.ID, p.Config.DryRun)
		}
	}()
}

func (p ProjectHandler) List() {
	Run(p, p.Printer.Print)
}

func (p ProjectHandler) Archive() {
	p.archiveProject(p.Config.ResourceId, p.Config.DryRun)
}

func (p ProjectHandler) ArchiveAll() {
	Run(p, p.archiveProjects)
}
func (p ProjectHandler) Delete() {
	p.deleteProject(p.Config.ResourceId, p.Config.DryRun)
}

func (p ProjectHandler) DeleteAll() {
	Run(p, p.deleteProjects)
}
