FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -o /out/federated-coordinator ./cmd/server
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/federated-coordinator /federated-coordinator
EXPOSE 8080
ENTRYPOINT ["/federated-coordinator"]
