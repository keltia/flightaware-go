# Main Makefile for fa-export
#
# XXX Need to be cleaned up at some point

VPATH=	flightaware:config

SRCS=	config.go client.go filters.go types.go \
	auth.go client.go datalog.go decode.go filters.go types.go

all:
	go build -v ./...

clean:
	go clean -v

push:
	git push --all origin
	git push --all backup
	git push --all gitlab
	git push --tags origin
	git push --tags backup
	git push --tags gitlab
