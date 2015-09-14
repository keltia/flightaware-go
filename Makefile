# Main Makefile for fa-export
#
# XXX Need to be cleaned up at some point

VPATH=	fa-export:fa-tail:flightaware:config:utils
DEST=	bin
GOBIN=	${GOPATH}/bin

SRCS=	config.go client.go fa-export.go utils.go cli.go version.go

all:	${DEST}/fa-export

install:
	go install fa-export/fa-export.go fa-export/cli.go fa-export/version.go

clean:
	go clean -v
	rm -f ${DEST}/fa-export

${DEST}/fa-export:    ${SRCS}
	go build -v -o $@ fa-export/fa-export.go fa-export/cli.go fa-export/version.go

push:
	git push --all
	git push --all origin
	git push --all backup
