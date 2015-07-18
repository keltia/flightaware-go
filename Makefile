# Main Makefile for fa-export

VPATH=	fa-export:fa-tail:flightaware:config:utils
DEST=	bin
GOBIN=	${GOPATH}/bin

SRCS=	config.go client.go fa-export.go utils.go cli.go

all:	${DEST}/fa-export ${DEST}/fa-tail

install:
	go install fa-export/fa-export.go fa-export/cli.go
	go install fa-tail/fa-tail.go

clean:
	go clean -v
	rm -f ${DEST}/fa-export

${DEST}/fa-export:    ${SRCS}
	go build -v -o $@ fa-export/fa-export.go fa-export/cli.go

${DEST}/fa-tail:   fa-tail/fa-tail.go
	go build -v -o $@ fa-tail/fa-tail.go

push:
	git push --all
	git push --all origin
	git push --all backup
