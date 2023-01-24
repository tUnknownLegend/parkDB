set-up-test-tool:
	go get -u -v github.com/mailcourses/technopark-dbms-forum@master && \
	go build github.com/mailcourses/technopark-dbms-forum

docker-run:
	sudo docker run -d --memory 2G --log-opt max-size=5M --log-opt max-file=3 --name park_perf -p 5000:5000 park

docker-build:
	sudo docker build -t park .

test-func:
	./technopark-dbms-forum func -u http://localhost:5000/api -r report.html

fill-db:
	./technopark-dbms-forum fill --url=http://localhost:5000/api --timeout=900

test-perf:
	./technopark-dbms-forum perf --url=http://localhost:5000/api --duration=600 --step=60	

docker-clean:
	sudo docker system prune -a && sudo docker volume prune

docker-fix:
	sudo killall containerd-shim
