module github.com/iteratec/gitlab-sanity-cli

go 1.16

replace github.com/iteratec/gitlab-sanity-cli/pkg => ./pkg

require (
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2
	github.com/xanzy/go-gitlab v0.50.1
)
