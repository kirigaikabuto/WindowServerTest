all: one
.PHONY: all

one:
	git add .
	git commit -m "sdsd"
	git push
	go build

build:
	go build && .\windowServerTest