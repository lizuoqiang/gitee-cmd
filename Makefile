build:
	rm -rf gitee-cmd && go build ./

install:
	go install github.com/lizuoqiang/gitee-cmd@$(shell git rev-parse HEAD)