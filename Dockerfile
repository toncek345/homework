FROM golang:1.19
WORKDIR /mnt/homework
COPY . .
RUN go build -o binary

ENTRYPOINT ./binary
