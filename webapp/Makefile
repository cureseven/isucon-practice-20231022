.PHONY: *

gogo: stop-services build truncate-logs start-services bench

build:
	cd golang && make app

stop-services:
	sudo systemctl stop nginx
	sudo systemctl stop isu-go
	ssh isucon-s2 "sudo systemctl stop mysql"

start-services:
	ssh isucon-s2 "sudo systemctl start mysql"
	sleep 5
	sudo systemctl start isu-go
	sudo systemctl start nginx

truncate-logs:
	sudo truncate --size 0 /var/log/nginx/access.log
	sudo truncate --size 0 /var/log/nginx/error.log
	ssh isucon-s2 "sudo truncate --size 0 /var/log/mysql/mysql-slow.log"
	ssh isucon-s2 "sudo chmod 777 /var/log/mysql/mysql-slow.log"
	sudo journalctl --vacuum-size=1K

pprof: TIME=60
pprof: PROF_FILE=~/pprof.samples.$(shell TZ=Asia/Tokyo date +"%H%M").$(shell git rev-parse HEAD | cut -c 1-8).pb.gz
pprof:
	curl -sSf "http://localhost:6060/debug/fgprof?seconds=$(TIME)" > $(PROF_FILE)
	go tool pprof $(PROF_FILE)

bench:
	/home/isucon/private_isu/benchmarker/bin/benchmarker -u /home/isucon/private_isu/benchmarker/userdata -t http://localhost/

kataribe:
	cd ../.. && sudo cat /var/log/nginx/access.log | ./kataribe
