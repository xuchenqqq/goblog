FROM ubuntu:trusty
MAINTAINER deepzz <deepzz.qi@gmail.com>

RUN apt-get update
RUN apt-get install -y ca-certificates
ENV MGO 172.17.42.1
ADD views /goblog/views
ADD conf /goblog/conf
ADD static /goblog/static
ADD goblog /goblog/goblog
RUN ["cp","/usr/share/zoneinfo/Asia/Shanghai","/etc/localtime"]

EXPOSE 80
EXPOSE 8080

VOLUME ["/goblog/log"]

WORKDIR /goblog
# CMD ["/goblog/goblog"]
ENTRYPOINT ["./goblog", "-m", "prod"] 