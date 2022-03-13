.PHONY: vendor
vendor:
	go mod tidy; go mod vendor;

test:
	go test -mod=vendor -race --coverpkg ./pkg/... -covermode atomic -coverprofile=.coverage.out ./pkg/...

bench:
	go test -mod=vendor -bench . ./pkg/... -benchmem

cover:
	go tool cover -func=.coverage.out

#
# tasks for ci
#
cienv:
	ci/actions.sh setup_go

citest:
	ci/actions.sh test_go
