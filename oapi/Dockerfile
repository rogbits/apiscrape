FROM golang
COPY go.mod /app/
COPY oapi/ /app/oapi
COPY lib/ /app/lib
WORKDIR /app/oapi/
RUN go build main.go
CMD ./main