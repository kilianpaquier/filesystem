# Code generated by craft; DO NOT EDIT.

CI_PROJECT_PATH := github.com/kilianpaquier/filesystem
GCI_CONFIG_PATH := build/ci/.golangci.yml

.PHONY: lint
lint:
	@gci write --skip-generated -s standard -s default -s "Prefix(${CI_PROJECT_PATH})" .
	@golangci-lint run -c ${GCI_CONFIG_PATH} --timeout 120s --fast --sort-results \
		--issues-exit-code 0 \
		--out-format colored-line-number $(ARGS)

.PHONY: test
test: lint
	@go test ./... -count 1

.PHONY: test-race
test-race: lint
	@go test ./... -race

.PHONY: test-cover
test-cover: lint
	@go test ./... -coverpkg=./... -covermode=count -coverprofile=coverage.out