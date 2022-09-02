FROM golang
COPY go.mod /app/
COPY omgr/ /app/omgr
COPY lib/ /app/lib
WORKDIR /app/omgr/
RUN go build main.go
CMD ./main