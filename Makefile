install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_clean_repo/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_copy_package/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_create_repo/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_delete_package/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_delete_repo/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_package_lister/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_package_versions/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_package_latest_version/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_repo_lister/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/aptly_upload/*.go
test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)
vet:
	go tool vet .
	go tool vet --shadow .
lint:
	golint -min_confidence 1 ./...
errcheck:
	errcheck -ignore '(Close|Write)' ./...
check: lint vet errcheck
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor
