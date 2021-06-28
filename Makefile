all: one
.PHONY: all

one:
	git add .
	git commit -m "sdsd"
	git push

build:
	go build && .\windowServerTest