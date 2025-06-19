# Build the application from source
FROM golang:1.24 AS build-stage
LABEL authors="kuraiji"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .env /app
COPY cmd /app/cmd
COPY internal /app/internal
# COPY public /app/public
RUN CGO_ENABLED=0 GOOS=linux go build -o /jidou ./cmd/jidou
# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /jidou /jidou
COPY public /public
EXPOSE 1323
USER nonroot:nonroot
CMD ["/jidou"]