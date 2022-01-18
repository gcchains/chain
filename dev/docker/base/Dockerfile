FROM debian:testing

ARG DEBIAN_FRONTEND=noninteractive
ARG DEBCONF_NOWARNINGS=yes

RUN echo "deb http://deb.debian.org/debian/ testing main non-free contrib" > /etc/apt/sources.list
# apt-get for scripting and apt for interactive use.
RUN apt-get update; apt-get -y upgrade; apt-get -y install locales tzdata

# set the locale
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && locale-gen
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL en_US.UTF-8     

# set location
RUN echo $TZ > /etc/timezone 

# user
RUN adduser --disabled-login --gecos '' gcchain
RUN apt-get -y install sudo
RUN echo 'gcchain ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# dev
RUN apt-get -y install python3 python3-pip

RUN apt-get -y install golang-go

RUN apt-get clean

ENV HOME /home/gcchain

USER gcchain
WORKDIR $HOME