SRCS=config/config.go flightaware/client.go fa-export/fa-export.go

.PATH: bin

all: bin/fa-export

clean:
	rm -f bin/fa-export

bin/fa-export: ${SRCS}
	go build -v -o bin/fa-export fa-export/fa-export.go
