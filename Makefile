.PHONY: clean releases

TAG=latest

clean:
	rm cover-generator-*.tar.gz

releases:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o cover-generator main.go    && tar czvf cover-generator-$(TAG)-linux_amd64.tar.gz cover-generator
	GOARCH=386   GOOS=linux CGO_ENABLED=0 go build -o cover-generator main.go    && tar czvf cover-generator-$(TAG)-linux_386.tar.gz cover-generator
	GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 go build -o cover-generator main.go   && tar czvf cover-generator-$(TAG)-darwin_amd64.tar.gz cover-generator
	GOARCH=386   GOOS=darwin CGO_ENABLED=0 go build -o cover-generator main.go   && tar czvf cover-generator-$(TAG)-darwin_386.tar.gz cover-generator
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o cover-generator main.go  && tar czvf cover-generator-$(TAG)-windows_amd64.tar.gz cover-generator
	GOARCH=386   GOOS=windows CGO_ENABLED=0 go build -o cover-generator main.go  && tar czvf cover-generator-$(TAG)-windows_386.tar.gz cover-generator
	rm ./cover-generator
