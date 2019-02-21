# Main Makefile for fa-export
#
# XXX Need to be cleaned up at some point

GO=		go
OPTS=	-ldflags="-s -w" -v

SRCS=	client.go client_subr.go config.go filters.go types.go \
	auth.go datalog.go decode.go utils.go \

ESRC=	cmd/fa-export/cli.go cmd/fa-export/fa-export.go \
	cmd/fa-export/utils.go

TSRC=	cmd/fa-tail/fa-tail.go cmd/fa-tail/cli.go

EBIN=	fa-export
TBIN=	fa-tail

all: ${EBIN} ${TBIN}
	${GO} build ${OPTS} ./cmd/...

${EBIN}:	${ESRC}

${TBIN}:	${TSRC}

install:	${EBIN} ${TBIN}
	${GO} install ./cmd/...

lint:
	gometalinter ./...

clean:
	${GO} clean -v ./cmd/...

push:
	git push --all origin
	git push --tags origin
