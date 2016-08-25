install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_clean_repo/aptly_clean_repo.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_copy_package/aptly_copy_package.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_create_repo/aptly_create_repo.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_delete_package/aptly_delete_package.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_delete_repo/aptly_delete_repo.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_package_lister/aptly_package_lister.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_package_version/aptly_package_version.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_repo_lister/aptly_repo_lister.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/aptly_upload/aptly_upload.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
vet:
	go tool vet .
	go tool vet .-shadow .
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
