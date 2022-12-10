FROM golang:1.17 as builder
RUN cd ..
RUN mkdir crabi-solution
WORKDIR crabi-solution
COPY . ./
ARG version=dev
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=$version" -o crabi-solution ./cmd/crabi-solution/main.go

FROM scratch
COPY --from=builder /go/crabi-solution/crabi-solution .

CMD ["./crabi-solution"]
