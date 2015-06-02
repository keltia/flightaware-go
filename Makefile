# Main Makefile for surv-export

VPATH=	fa-export:flightaware:config
DEST=	bin

SRCS=	config.go client.go main.go

all:	${DEST}/fa-export

clean:
	rm -f ${DEST}/fa-export

${DEST}/fa-export:    ${SRCS}
	go build -v -o $@ fa-export/main.go

push:
	git push --all
	git push --all backup
	git push --all bitbucket
