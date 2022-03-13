#!/bin/bash -ex

COMMAND="${@:-help}"
NAME=${PWD##*/}
WD=${PWD}

function setup_go() {
    GOPATH=/opt/ci/go-1.17.4
    MODULE=`grep '^module ' go.mod | awk '{print $2}'`

	mkdir -p /opt/ci
	test -e /opt/ci/go || \
        (wget -q https://go.dev/dl/go1.17.4.linux-amd64.tar.gz && tar xf go1.17.4.linux-amd64.tar.gz -C /opt/ci/ && rm -f go1.17.4.linux-amd64.tar.gz)

    cat << EOS > /opt/ci/envgorc
export PATH=/opt/ci/go/bin:${GOPATH}/bin:$PATH
export GOROOT=/opt/ci/go
export GOPATH=${GOPATH}
export GO111MODULE=on
EOS

	. /opt/ci/envgorc
    go version
    cd $GOROOT

    # for codecov
    go install github.com/axw/gocov/gocov@latest
    go install github.com/AlekSi/gocov-xml@latest

    # for coveralls
    # go install github.com/jandelgado/gcov2lcov@latest
}

function test_go() {
	. /opt/ci/envgorc

	go test -mod=vendor -race --coverpkg ./pkg/... -covermode atomic -coverprofile=.coverage.out ./pkg/...

    # for codecov
    gocov convert .coverage.out | gocov-xml > /tmp/coverage.xml

    # for coveralls
    # gcov2lcov -infile .coverage.out -outfile /tmp/coverage.lcov
}

$COMMAND
