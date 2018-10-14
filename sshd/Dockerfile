FROM ubuntu:16.04 AS ubuntu

RUN apt-get update && apt-get install -y openssh-server
RUN mkdir /var/run/sshd
# RUN echo 'root:qwe123' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

# ENV NOTVISIBLE "in users profile"
# RUN echo "export VISIBLE=now" >> /etc/profile

# EXPOSE 22
# CMD ["/usr/sbin/sshd", "-D"]

FROM docker

RUN apk add openssh py-pip
RUN pip install docker-compose
RUN echo 'root:qwe123' | chpasswd
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

COPY --from=ubuntu /var/run/sshd /var/run/sshd
COPY --from=ubuntu /etc/ssh /etc/ssh
COPY --from=ubuntu /etc/pam.d/sshd /etc/pam.d/sshd

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
