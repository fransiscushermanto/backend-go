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

echo "Preparing database (Dropping and Recreating)..."

MAX_DB_OPS_ATTEMPTS=5 # Try dropping/creating 5 times (5 * 1 = 5 seconds)
DB_OPS_ATTEMPT=0
SUCCESS=false

DB_NAME=$(echo "$DATABASE_URL" | sed -r 's/.*\/([a-zA-Z0-9_]+)(\?.*)?/\1/')
PG_HOST=$(echo "$DATABASE_URL" | sed -r 's/.*@([^:]+)(:.*)?\/.*/\1/')
PG_PORT=$(echo "$DATABASE_URL" | sed -r 's/.*:([0-9]+)\/.*/\1/')
PG_USER=$(echo "$DATABASE_URL" | sed -r 's/.*:\/\/(.*):.*@.*/\1/')
PG_PASSWORD=$(echo "$DATABASE_URL" | sed -r 's/.*:(.*)@.*/\1/')

until $SUCCESS || [ $DB_OPS_ATTEMPT -ge $MAX_DB_OPS_ATTEMPTS ]; do
  echo "Dropping database '$DB_NAME' (Attempt $((DB_OPS_ATTEMPT + 1))/$MAX_DB_OPS_ATTEMPTS)..."
  if docker compose -f "$PROJECT_ROOT/docker-compose.yaml" exec db psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d postgres -c "DROP DATABASE IF EXISTS \"$DB_NAME\""; then
    echo "Creating database '$DB_NAME' (Attempt $((DB_OPS_ATTEMPT + 1))/$MAX_DB_OPS_ATTEMPTS)..."
    if docker compose -f "$PROJECT_ROOT/docker-compose.yaml" exec db psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d postgres -c "CREATE DATABASE \"$DB_NAME\""; then
      SUCCESS=true
    else
      echo "Create database failed. Retrying..."
      sleep 1
    fi
  else
    echo "Drop database failed. Retrying..."
    sleep 1
  fi
  DB_OPS_ATTEMPT=$((DB_OPS_ATTEMPT+1))
done

if ! $SUCCESS; then
  echo "Error: Failed to drop/create database '$DB_NAME' after $MAX_DB_OPS_ATTEMPTS attempts."
  echo "Please check Docker logs for 'backend-db-1' using 'docker compose logs db' and review permissions/connection."
  exit 1
fi
echo "Database '$DB_NAME' recreated."