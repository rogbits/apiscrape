FROM golang
COPY go.mod /app/
COPY owrk/ /app/owrk
COPY lib/ /app/lib
WORKDIR /app/owrk/
RUN go build main.go
CMD ./main