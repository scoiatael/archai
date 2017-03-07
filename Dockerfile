FROM golang:1.8
ENV HOME /opt/app
ENV GOPATH $HOME
ENV PATH $GOPATH/bin:$PATH
RUN go get github.com/tools/godep
RUN mkdir -p $HOME/src/github.com/scoiatael/
COPY . $HOME/src/github.com/scoiatael/archai

WORKDIR $HOME/src/github.com/scoiatael/archai
RUN godep restore
RUN go install github.com/scoiatael/archai
EXPOSE 8080
ENTRYPOINT ["archai"]
