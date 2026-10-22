FROM golang:1.27.3-alpine AS build

WORKDIR /src
COPY . .
RUN mkdir /out && CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/ipify-api .

FROM scratch

COPY --from=build /out/ipify-api /ipify-api
USER 65532:65532
EXPOSE 3000/tcp
ENTRYPOINT ["/ipify-api"]
