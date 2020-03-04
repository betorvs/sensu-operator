# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.0.1] - 2020-03-02

### Added
- Added Controllers: sensubackend, sensuagent, sensunamespace, sensuasset, sensuhandler, sensucheck, sensumutator, sensufilter.
- Added extra dirs: appcontext (create repository tool), config (environment variables), domain (sensu gateway structs and interfaces), gateway (sensu http gateway configurations), usecase (application logical interface between controller and gateway), utils(some utilities functions).
- Added sensu-certs with json files to help to create all certificates using cfssl tools.
- Added scripts for local test: run.sh, remove.sh, test-all.sh.
