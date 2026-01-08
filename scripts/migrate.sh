#!/bin/bash
set -e

# --- Configuration ---
# Determine script's directory and project root
SCRIPT_DIR=$(dirname "$0")
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
MIGRATIONS_DIR=${MIGRATIONS_DIR:-"./migrations"}

# --- Load Environment Variables from .env file (Fallback for local execution) ---
# This makes the script more portable and self-contained for local execution.
# It will load .env variables if they are not already set in the shell.
if [ -f "$PROJECT_ROOT/.env" ]; then
  echo "Sourcing $PROJECT_ROOT/.env file for migration script..."
  set -a # Automatically export all variables loaded by 'source'
  source "$PROJECT_ROOT/.env"
  set +a # Disable automatic export afterwards
else
  echo "Warning: .env file not found at $PROJECT_ROOT/.env. Relying on shell environment variables."
fi

# Ensure DATABASE_URL is set (from shell env or .env file)
if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL environment variable is not set."
  echo "Example: export DATABASE_URL=\"postgresql://user:password@localhost:5432/mydatabase?sslmode=disable\""
  exit 1
fi

echo "Using database: $(echo $DATABASE_URL | sed 's/\/\/.*@/\/\/user:*****@/')" # Mask password
echo "Migrations directory: $MIGRATIONS_DIR"

case "$1" in
   up)
    echo "Running database migrations UP..."
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
    echo "Migrations UP complete."
    ;;
  down)
    # Be careful with 'down' in production! Usually 'down 1' for last migration, or 'down all' for development resets.
    if [ -n "$2" ]; then
      # Down specific number of migrations
      echo "Running $2 migration(s) DOWN..."
      read -p "WARNING: Running DOWN migrations can lead to data loss. Are you sure? (type 'yes'): " CONFIRM
      if [ "$CONFIRM" != "yes" ]; then
        echo "Migration aborted."
        exit 1
      fi
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down "$2"
    else
      # Down all migrations (existing behavior)
      echo "Running database migrations DOWN..."
      read -p "WARNING: Running DOWN migrations can lead to data loss. Are you sure? (type 'yes'): " CONFIRM
      if [ "$CONFIRM" != "yes" ]; then
        echo "Migration aborted."
        exit 1
      fi
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down
    fi
    echo "Migrations DOWN complete."
    ;;
  goto)
    if [ -z "$2" ]; then
      echo "Usage: $0 goto <version>"
      exit 1
    fi
    echo "Migrating to version $2..."
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" goto "$2"
    echo "Migration to version $2 complete."
    ;;
  status)
    echo "Checking migration status..."
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" status
    ;;
  force)
    echo "Forcing migrate"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$2"
    echo "Migration forcely moved to version $2"
    ;;
  create)
    if [ -z "$2" ]; then
      echo "Usage: $0 create <name_of_migration>"
      exit 1
    fi
    echo "Creating new migration: $2"
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$2" 
    echo "Migration files created in $MIGRATIONS_DIR"
    ;;
  *)
    echo "Usage: $0 {up|down|goto <version>|status|create <name>}"
    exit 1
    ;;
esac