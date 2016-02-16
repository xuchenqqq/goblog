FROM ubuntu:trusty
MAINTAINER chenqijing2 <chenqijing2@163.com>

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD views /goblog/views
ADD conf /goblog/conf
ADD static /goblog/static
ADD beego_goblog /goblog/goblog
RUN ["cp","/usr/share/zoneinfo/Asia/Shanghai","/etc/localtime"]

EXPOSE 80
EXPOSE 8080

VOLUME ["/goblog/log"]

WORKDIR /goblog
# CMD ["/goblog/goblog"]
ENTRYPOINT ["./goblog"]
