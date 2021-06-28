all: one
.PHONY: all

one:
	git add .
	git commit -m "sdsd"
	git push
	go build
	./windoServerTest.exe -u tleugazy_erasil@gmail.com