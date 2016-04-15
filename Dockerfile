FROM golang:1.6.1

RUN echo Asia/Tokyo | tee /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata

RUN git clone https://github.com/if1live/haru.git /haru/src/github.com/if1live/haru
WORKDIR /haru/src/github.com/if1live/haru

ENV GOPATH /haru/
ENV PATH $GOPATH/bin/:$PATH

RUN go get github.com/tools/godep
RUN godep restore

EXPOSE 3000

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh
