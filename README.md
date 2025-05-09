# Realtime Chat Application

## Setup Instructions

1. Clone the repository
2. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

3. Update the `.env` file with your database credentials

4. Start the services with Docker:
```bash
docker-compose up -d
```

5. Run the application:
```bash
go run main.go
```

## Environment Variables

The following environment variables need to be set in your `.env` file:

- `POSTGRES_USER`: PostgreSQL username
- `POSTGRES_PASSWORD`: PostgreSQL password
- `POSTGRES_DB`: Database name
- `POSTGRES_HOST`: Database host (default: localhost)
- `POSTGRES_PORT`: Database port (default: 5432)
- `REDIS_ADDR`: Redis address (default: localhost:6379)

## Security Note

Never commit the `.env` file to version control. The `.env` file contains sensitive information and is already added to `.gitignore`.