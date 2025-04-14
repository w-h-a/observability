FROM golang AS build
WORKDIR /observability-backend
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/observability-backend ./

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/observability-backend /bin/observability-backend
ENTRYPOINT [ "/bin/observability-backend" ]