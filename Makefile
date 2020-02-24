.PHONY: all create build clean start doc install checkout up help 

all: clean create build checkout

create:
	@mkdir bin
	@cp ./src/stepstep/conf/config.ini ./bin/

build:
	@go build -v -tags=jsoniter ./src/stepstep 
	@mv stepstep ./bin

clean:
	@rm -rf bin
	@rm -rf stepstep 
	@go clean -i .

start:
	cd bin && ./stepstep

doc:
	@sh tools/showdoc_api.sh src/stepstep/api

install:
	@go install stepstep

checkout:
	@svn checkout svn://127.0.0.1:8099/dhq_pro 
	@mv server bin/data

up:
	@cd bin/data && svn up 

rsync_release:
	@sh tools/rsync_release.sh

release: build rsync_release

vendor: 
	@cd src/stepstep && govendor add +e

units:
	go install units	

test:
	cd bin && ./units

help:
	@echo "make create: 创建bin目录和复杂配置文件"
	@echo "make build: 指定编译 tags=jsonier"
	@echo "make clean: 删除对应的编译数据"
	@echo "make start: 启动程序"
	@echo "make doc: 重新生成对应的Api接口说明到showdoc上"
	@echo "make install: 安装程序"
	@echo "make checkout: 从svn中检出xls的json数据"
	@echo "make up: 更新json数据"
	@echo "make release: 发布程序到服务器"
