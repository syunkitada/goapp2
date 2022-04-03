.PHONY: vendor
vendor:
	go mod tidy; go mod vendor;

test:
	go test -mod=vendor -race --coverpkg ./pkg/... -covermode atomic -coverprofile=.coverage.out ./pkg/...

bench:
	go test -mod=vendor -bench . ./pkg/... -benchmem

cover:
	go tool cover -func=.coverage.out
	go tool cover -html=.coverage.out -o .coverage.html

srv:
	python3 -m http.server 8000

#
# tasks for ci
#
cienv:
	ci/actions.sh setup_go

citest:
	ci/actions.sh test_go
