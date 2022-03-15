# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.2]
## Added
New parameter for resource projects: `-p, --project-type`.
    This parameter selects the project visibility level (public, internal, private). The default visibility level is set to internal.

## [1.0.1]
## Changes
Fix version check (https://github.com/iteratec/gitlab-sanity-cli/issues/1)

## [1.0.0]
## Changes
Mark as public release

## [0.9.0]
## Changes
Fixed security findings

## [0.0.8]
## Added
basic handler test added

## [0.0.7]
## Added
Added parameter to define gitlab host and skip certificate validations

## Changed
Fix race problem in fetching result and closing ingress channel

## [0.0.6]
### Added
Switched to AbstractHandler Interface implementation

## [0.0.5]
### Added
Gitlab token can be set also from env or file

## [0.0.4]
### Added
Deletion of GroupRunners by query filters
Parameter for Version Number output

## [0.0.3]
### Added
Deletion of Projects by query and age filters

### Changes
Argument parsing and code structure refactored

## [0.0.2]
### Added
CHANGELOG.md to keep changes  
VERSION to set/define versions

### Changes
Use templates for listing output

## [0.0.1]
### Added
- Initial PoC Release

