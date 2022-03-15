package gitlabsanitycli

import (
	"github.com/voxelbrain/goptions"
)

const (
	// PageSize is the limit of results per result page (from API). Maximum is 100
	PageSize = 100

	// Config Parameters
	List        = "list"
	Archive     = "archive"
	ArchiveAll  = "archive-all"
	Delete      = "delete"
	DeleteAll   = "delete-all"
	Project     = "project"
	Runner      = "runner"
	GroupRunner = "groupRunner"
	User        = "user"

	// GitLab API OK Response Code
	gitlabReponseCodeOk = 204
)

var (
	// ConcurrentAPIRequest set the amount of concurrent goroutines
	ConcurrentAPIRequest = 10
)

// Application Configuration
type Config struct {
	URL                   string        `goptions:"-u, --url, description='Gitlab API Url, can also be set via env \"GITLAB_URL\"', obligatory"`
	Insecure              bool          `goptions:"--insecure, description='Skip certificate Verfication for Gitlab API URL, (bool)'"`
	Token                 string        `goptions:"-t, --token, description='Gitlab API access token, can also be set via env \"GITLAB_TOKEN\" or file \".token\"'"`
	Operation             string        `goptions:"-o, --operation, description='Operation to start, (list, archive, archive-all, delete, delete-all)', obligatory"`
	Resource              string        `goptions:"-r, --resource, description='Resource to interact with, (project, runner, groupRunner, user)', obligatory"`
	ResourceId            int           `goptions:"-i, --identifier, description='Resource ID to interact with, (int)'"`
	ProjectType           string        `goptions:"-p, --project-type, description='Type of project (internal, private, public), (default: internal), (string)'"`
	Age                   int           `goptions:"-a, --age, description='Filter by last activity in months (not available for runner), (int)'"`
	Query                 string        `goptions:"-q, --query, description='Search by name, (string)'"`
	Status                string        `goptions:"-s, --state, description='Filter list by state, (string)'"`
	DryRun                bool          `goptions:"-d, --dry-run, description='Dry run, does not change/delete any resources, (bool)'"`
	Help                  goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Version               bool          `goptions:"-v, --version, description='Show Version'"`
	NumConCurrentApiCalls int           `goptions:"-n, --num-concurrent-api-calls, description='Limit the amount of concurrent API Calls (default: 10), (int)'"`
}
