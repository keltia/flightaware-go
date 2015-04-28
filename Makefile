SRCS=config/config.go flightaware/client.go fa-export/main.go

all: main

main: ${SRCS}
	go build -v fa-export/main.go
