#! /bin/bash
VERSION=v0.1.$(git rev-parse --short HEAD)
echo $VERSION

[ -d /usr/local/go1.7 ] && {
    export GOROOT=/usr/local/go1.7/
    export PATH=/usr/local/go1.7/bin/:$PATH
}

export GOPATH=$(pwd)/gopath/
go version

ROOT=./
mkdir -p gopath/src/elastic-trib/
rsync -arvz ${ROOT}/Makefile ${ROOT}/elastic-trib.yaml ${ROOT}/*.go gopath/src/elastic-trib
rsync -arvz ${ROOT}/vendor/ gopath/src/elastic-trib
rsync -arvz ${ROOT}/vendor/ gopath/src/

output=${PWD}/output

echo "Building elastic-trib"
cd gopath/src/elastic-trib
go build -i -ldflags "-X main.gitCommit=${COMMIT} -X main.version=${VERSION}" -o ${output}/elastic-trib .

if [ "$?" != "0" ]; then 
    echo "build elastic-trib failed"
    cd - && rm -rf ./gopath
    exit 1
fi

cp -rf elastic-trib.yaml ${output}/

cd - && rm -rf ./gopath

echo "end building"
