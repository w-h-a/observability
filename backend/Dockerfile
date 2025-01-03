FROM golang AS build
WORKDIR /trace-blame-backend
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/trace-blame-backend ./

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/trace-blame-backend /bin/trace-blame-backend
ENTRYPOINT [ "/bin/trace-blame-backend" ]