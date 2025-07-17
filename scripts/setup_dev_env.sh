#!/bin/bash
set -e

SCRIPT_DIR=$(dirname "$0")
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)

# Source .env file from project root
if [ -f "$PROJECT_ROOT/.env" ]; then
  echo "Sourcing $PROJECT_ROOT/.env file..."
  source "$PROJECT_ROOT/.env"
else
  echo "Error: .env file not found in project root ($PROJECT_ROOT/.env)."
  echo "Please create a .env file with DATABASE_URL and other configurations."
  exit 1
fi

# Ensure DATABASE_URL is set after sourcing .env
if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL environment variable not found after sourcing .env."
  echo "Please ensure DATABASE_URL is defined in your .env file."
  exit 1
fi

MIGRATE_CLI_PATH="$HOME/go/bin/migrate" # Adjust if your GOBIN is different
MIGRATE_SCRIPT="$SCRIPT_DIR/migrate.sh"

echo "\n--- Setting up Go API Development Environment ---\n"

# --- Step 1: Check for Docker and Go installations ---
echo "1. Checking for Docker and Go installations..."
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH. Please install Docker Desktop."
    exit 1
fi
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH. Please install Go."
    exit 1
fi
echo "Docker and Go are installed."

# --- Step 2: Build/Pull Docker Compose services ---
echo "2. Starting and building Docker Compose services (db)..."
# -f specifies the docker-compose file path relative to current dir
# --wait ensures services are healthy before proceeding
docker compose -f "$PROJECT_ROOT/docker-compose.yaml" up --build -d --wait db
echo "Database service started and is healty."

# --- Step 3: Ensure current Go dependencies are downloaded ---
echo "3. Running go mod tidy to ensure all Go dependencies are downloaded..."
(cd "$PROJECT_ROOT" && go mod tidy) # Run go mod tidy from the project root
if [ $? -ne 0 ]; then
  echo "Error: go mod tidy failed. Please check your go.mod file."
  exit 1
fi
echo "Go dependencies are up-to-date."

# --- Step 4: Run Database Migrations (now from API container) ---
# Now, we bring up the API container specifically to run migrations.
echo "4. Starting API service (api) for migrations..."
# We use `docker compose run` here because it runs a command and exits.
# It waits for `db` to be healthy implicitly.
# We are overriding the default `CMD` of the API container to run the migration script.
make migrate action="up" PROJECT_ROOT="$PROJECT_ROOT"
echo "Database migrations applied."

# --- Step 5: Start the API service for regular operation ---
echo "5. Starting API service (api) for regular operation..."
# This will now start the API service with its normal CMD and keep it running in background.
docker compose -f "$PROJECT_ROOT/docker-compose.yaml" up -d api
echo "API service (api) started and running."

echo "--- Development Environment Setup Complete! ---"
echo "Your API server is now running in Docker Compose and accessible at http://localhost:$APP_PORT"
echo "Code changes on your host will trigger hot-reloading inside the container."
echo "You can view logs with: docker compose -f $PROJECT_ROOT/docker-compose.yaml logs -f api"