package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	gitlabsanitycli "github.com/iteratec/gitlab-sanity-cli/pkg"
	"github.com/voxelbrain/goptions"
	"github.com/xanzy/go-gitlab"
)

const (
	version       = "1.0.2"
	hoursPerMonth = 30 * 24
	newLine       = "\r\n"
)

func init() {
	if len(os.Args) > 1 {
		if strings.Compare(os.Args[1], "-v") == 0 || strings.Compare(os.Args[1], "-version") == 0 || strings.Compare(os.Args[1], "--version") == 0 {
			fmt.Printf("Gitlab Sanity CLI (@iteratec)\nVERSION: %v\n", version)
			os.Exit(0)
		}
	}
}

func parseParameters() gitlabsanitycli.Config {
	config := gitlabsanitycli.Config{ResourceId: -1, Resource: gitlabsanitycli.Project, Age: 36, DryRun: false}
	goptions.ParseAndFail(&config)
	return config
}

func validateParameters(config *gitlabsanitycli.Config) error {
	config.Age = int(math.Max(0, float64(config.Age))) * hoursPerMonth

	// Retrieve Gitlab token
	if len(config.Token) > 0 {
		log.Print("Using GitLab token from args")
	} else {
		tokenFromEnv, hasToken := os.LookupEnv("GITLAB_TOKEN")
		if hasToken {
			log.Print("Using GitLab token from environment")
			config.Token = tokenFromEnv
		} else {
			tokenFromFile, err := ioutil.ReadFile(".token")
			if err == nil && tokenFromFile != nil && len(tokenFromFile) > 0 {
				log.Print("Using GitLab token from file")
				config.Token = strings.TrimRight(string(tokenFromFile), newLine)
			} else {
				goptions.PrintHelp()
				return errors.New("GitLab token missing")
			}
		}
	}
	if len(config.ProjectType) < 1 {
		config.ProjectType = "internal"
	}
	return nil
}

func newHandlers(git *gitlab.Client, config gitlabsanitycli.Config) map[string]map[string]func() {
	handlers := make(map[string]map[string]func())
	printer := &gitlabsanitycli.TemplatePrinter{}
	projectHandler := gitlabsanitycli.ProjectHandler{Git: git, Printer: printer, Config: &config}
	userHandler := gitlabsanitycli.UserHandler{Git: git, Printer: printer, Config: &config}
	runnerHandler := gitlabsanitycli.RunnerHandler{Git: git, Printer: printer, Config: &config}
	groupRunnerHandler := gitlabsanitycli.GroupRunnerHandler{Git: git, Printer: printer, Config: &config}
	handlers[gitlabsanitycli.List] = make(map[string]func())
	handlers[gitlabsanitycli.List][gitlabsanitycli.Project] = projectHandler.List
	handlers[gitlabsanitycli.List][gitlabsanitycli.User] = userHandler.List
	handlers[gitlabsanitycli.List][gitlabsanitycli.Runner] = runnerHandler.List
	handlers[gitlabsanitycli.List][gitlabsanitycli.GroupRunner] = groupRunnerHandler.List
	handlers[gitlabsanitycli.Archive] = make(map[string]func())
	handlers[gitlabsanitycli.Archive][gitlabsanitycli.Project] = projectHandler.Archive
	handlers[gitlabsanitycli.ArchiveAll] = make(map[string]func())
	handlers[gitlabsanitycli.ArchiveAll][gitlabsanitycli.Project] = projectHandler.ArchiveAll
	handlers[gitlabsanitycli.Delete] = make(map[string]func())
	handlers[gitlabsanitycli.Delete][gitlabsanitycli.Project] = projectHandler.Delete
	handlers[gitlabsanitycli.Delete][gitlabsanitycli.Runner] = runnerHandler.Delete
	handlers[gitlabsanitycli.Delete][gitlabsanitycli.GroupRunner] = groupRunnerHandler.Delete
	handlers[gitlabsanitycli.DeleteAll] = make(map[string]func())
	handlers[gitlabsanitycli.DeleteAll][gitlabsanitycli.Project] = projectHandler.DeleteAll
	handlers[gitlabsanitycli.DeleteAll][gitlabsanitycli.Runner] = runnerHandler.DeleteAll
	handlers[gitlabsanitycli.DeleteAll][gitlabsanitycli.GroupRunner] = groupRunnerHandler.DeleteAll
	return handlers
}

func newGitLabClient(config gitlabsanitycli.Config) (git *gitlab.Client, err error) {
	// Skip certificate validation if enabled
	tr := &http.Transport{}
	tr.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS13}

	if config.Insecure {
		tr.TLSClientConfig.InsecureSkipVerify = true
	}
	client := &http.Client{Transport: tr}
	log.Printf("Connect to %s", config.URL)

	// Initialize git
	git, err = gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.URL), gitlab.WithHTTPClient(client))
	return git, err
}

func setNumConcurrentApiCalls(config gitlabsanitycli.Config) {
	// Retrieve Max Concurrent API Call
	if config.NumConCurrentApiCalls != 0 {
		gitlabsanitycli.ConcurrentAPIRequest = config.NumConCurrentApiCalls
	} else {
		numConCurrentApiCallsFromEnv, hasnumConCurrentApiCallsFromEnv := os.LookupEnv("NUM_CONCURRENT_API_CALLS")
		if hasnumConCurrentApiCallsFromEnv {
			num, err := strconv.Atoi(numConCurrentApiCallsFromEnv)
			if err == nil {
				gitlabsanitycli.ConcurrentAPIRequest = num
			}
		}
	}
}

func main() {
	// Parse arguments
	config := parseParameters()
	if err := validateParameters(&config); err != nil {
		log.Fatal(err)
	}

	setNumConcurrentApiCalls(config)

	// Create git client
	git, err := newGitLabClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Register handlers
	handlers := newHandlers(git, config)

	// Determine handler and run it
	if handlers[config.Operation][config.Resource] != nil {
		handlers[config.Operation][config.Resource]()
	} else {
		goptions.PrintHelp()
		log.Fatalf("Operation %v on resource %v not supported", config.Operation, config.Resource)
	}
}
