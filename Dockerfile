FROM openjdk

RUN yum update -y
RUN yum install -y vim make

WORKDIR /root
