current_dir=$(shell pwd)
version=$(VERSION)
tail=$(tail)
project_dir=socialserver

clean:
	cd $(current_dir)/../logs && rm -f ./socialserver.log
	cd $(current_dir)/src/cmd/socialserver && rm -rf ./socialserver
	sudo chmod +777 $(current_dir)/build/kill.sh
	$(current_dir)/build/kill.sh

build: clean
	chmod +x ./build/*.sh
#	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_dir) --name build_$(project_dir) ccr.ccs.tencentyun.com/webankpartners/golang-ext:v1.15.6 /bin/bash /go/src/github.com/WeBankPartners/$(project_dir)/build/build-server.sh
	docker run --rm -v $(current_dir):/tmp/$(project_dir) --name build_$(project_dir) golang:1.21 /bin/bash /tmp/$(project_dir)/build/build-server.sh



start: build
	cd $(current_dir)/scripts && source ./environment.sh && $(current_dir)/src/cmd/socialserver/socialserver -c $(current_dir)/src/configs/socialserver.yaml &
	#echo ".env===="${TASKSERVER_DB_HOST}


restart:
	$(current_dir)/build/restart-container.sh $(version)

log:
	$(current_dir)/build/log-container.sh $(version) $(tail)

