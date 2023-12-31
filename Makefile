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
	cd $(current_dir)/scripts && source ./environment.sh && $(current_dir)/src/cmd/socialserver/socialserver -c $(current_dir)/src/configs/socialserver.yaml
	#后台运行
# 	cd $(current_dir)/scripts && source ./environment.sh && $(current_dir)/src/cmd/socialserver/socialserver -c $(current_dir)/src/configs/socialserver.yaml &
	#echo ".env===="${TASKSERVER_DB_HOST}
image:
	 docker build -t rrs/socialserver:$(version) .

start-container:image
	sudo chmod +777 $(current_dir)/build/remove-container.sh
	$(current_dir)/build/remove-container.sh $(version)
	#docker run  -p 8808:8808 -p 9808:9808  -e "PGSQL_PORT=5432" --name syncserver-$(version) -d fancyflink/syncserver:$(version)
	docker run -p 8808:8808 --name socialserver-$(version) -d rrs/socialserver:$(version)
	docker ps -all
	docker logs -f --tail=100 socialserver-$(version)
#	docker logs  socialserver-$(version)
#   docker exec -it 6dfb762e67dc /bin/sh

start-com:
	docker-compose down
	docker-compose up -d --build
	docker-compose logs -f --tail=50
#	docker logs -f --tail=100 go-socialapp-socialserver-1

restart:
	$(current_dir)/build/restart-container.sh $(version)

log:
	$(current_dir)/build/log-container.sh $(version) $(tail)

