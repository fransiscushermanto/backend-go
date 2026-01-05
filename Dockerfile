# Stage 1: Build the Go application
# Use a specific Go version from official Golang image with Alpine Linux
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's build cache.
# This ensures that if only source code changes, dependencies aren't re-downloaded.
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/air-verse/air@latest

COPY . .

# Build the Go application.
# CGO_ENABLED=0: Disable CGO for static binaries, making the final image smaller and more portable.
# GOOS=linux: Target Linux OS for the build (essential for Alpine base image).
# -ldflags="-s -w": Strips debugging information and symbol tables, reducing binary size.
# -o /go/bin/api: Specifies the output executable name and its location.
# ./cmd/api: The path to your main package for the API server application.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/api ./cmd/api

# Stage 2: Create the final, smaller runtime image
# Use a minimal base image like Alpine Linux for a tiny footprint
FROM alpine:latest AS base

# Define a non-root user and group
# Using ARGs for flexibility, though you could hardcode if preferred
ARG APP_USER=appuser
ARG APP_GROUP=appgroup
#Common UID for non-root users
ARG APP_UID=1000 
#Common GID for non-root users
ARG APP_GID=1000 

# Create the group and user.
# -g $APP_GID: specify group ID.
# -G $APP_GROUP: add user to this primary group.
# -u $APP_UID: specify user ID.
# -s /bin/sh: set default shell for the user.
# -D: don't assign a password.
RUN addgroup -g $APP_GID $APP_GROUP \
    && adduser -u $APP_UID -G $APP_GROUP -s /bin/sh -D $APP_USER

# Install ca-certificates required for HTTPS communication (e.g., connecting to databases with SSL)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the compiled executable from the 'builder' stage to the final image
COPY --from=builder /go/bin/api .
COPY --from=builder /go/bin/migrate /usr/local/bin/

# Copy your application's configuration file.
# Your Go app expects this at `configs/config.yaml` relative to its working directory.
COPY configs/ ./configs/
COPY scripts/ ./scripts/
COPY migrations/ ./migrations/

COPY public_key.pem ./
COPY private_key.pem ./

RUN chown -R $APP_USER:$APP_GROUP /app
RUN chmod +x /app/scripts/migrate.sh
RUN chown $APP_USER:$APP_GROUP /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate

# Stage 3: Development environment
FROM golang:1.25-alpine AS development

ARG APP_USER=appuser
ARG APP_GROUP=appgroup
ARG APP_UID=1000 
ARG APP_GID=1000 

RUN addgroup -g $APP_GID $APP_GROUP \
    && adduser -u $APP_UID -G $APP_GROUP -s /bin/sh -D $APP_USER

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy everything including development tools
COPY --from=builder /go/bin/api .
COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY --from=builder /go/bin/air /usr/local/bin/

COPY configs/ ./configs/
COPY scripts/ ./scripts/
COPY migrations/ ./migrations/
COPY public_key.pem ./
COPY private_key.pem ./

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOROOT="/usr/local/go"

RUN chown -R $APP_USER:$APP_GROUP /app
RUN chmod +x /app/scripts/migrate.sh
RUN chown $APP_USER:$APP_GROUP /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate
RUN chown $APP_USER:$APP_GROUP /usr/local/bin/air
RUN chmod +x /usr/local/bin/air

USER $APP_USER

EXPOSE 8080

# Stage 4: Production environment
FROM base AS production

RUN chown -R $APP_USER:$APP_GROUP /app

USER $APP_USER

# Production startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'export PUBLIC_KEY="$(cat /app/public_key.pem)"' >> /app/start.sh && \
    echo 'export PRIVATE_KEY="$(cat /app/private_key.pem)"' >> /app/start.sh && \
    echo 'exec ./api' >> /app/start.sh && \
    chmod +x /app/start.sh

EXPOSE 8080

FROM production