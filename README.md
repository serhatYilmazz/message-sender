# Message Sender API

A robust message management and webhook delivery system built with Go, featuring message persistence, scheduled delivery, and reliable webhook integration using the outbox pattern.

## 🚀 Features

- **Message Management**: Create and retrieve messages via REST API
- **Scheduled Delivery**: Configurable scheduler for automatic message processing
- **Webhook Integration**: Reliable webhook delivery with retry mechanisms
- **Caching**: Redis-based caching for webhook delivery records
- **Database Persistence**: PostgreSQL for reliable message storage
- **API Documentation**: Auto-generated Swagger documentation
- **Outbox Pattern**: Ensures reliable message delivery
- **Docker Support**: Complete containerization with Docker Compose

## 🏗️ Architecture

- **Backend**: Go with Fiber web framework
- **Database**: PostgreSQL for message persistence
- **Cache**: Redis for webhook delivery tracking
- **Messaging**: Outbox pattern for reliable delivery
- **Documentation**: Swagger/OpenAPI integration

## 📋 Prerequisites

Before running the application, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (version 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0+)

## 🐳 Docker Setup - Step by Step

### Step 1: Clone and Navigate to Project
```bash
# Clone the repository (if not already done)
git clone <repository-url>
cd cool-project
```

### Step 2: Build and Start All Services
```bash
# Build and start all services (PostgreSQL, Redis, and Application)
docker-compose up --build
```

This command will:
- Build the Go application Docker image
- Pull PostgreSQL 15 Alpine image
- Pull Redis 7 Alpine image
- Start all services with proper networking
- Automatically run database migrations

### Step 3: Verify Services are Running
```bash
# In a new terminal, check if all services are running
docker-compose ps
```

You should see three services running:
- `cool-project-app-1` (Go application on port 8080)
- `cool-project-postgres-1` (PostgreSQL on port 5432)
- `cool-project-redis-1` (Redis on port 6379)

### Step 4: Access the Application

Once all services are running, you can access:

- **API Base URL**: http://localhost:8080
- **Swagger Documentation**: http://localhost:8080/docs
- **Health Check**: http://localhost:8080/api/messages

### Step 5: Test the API

#### Get All Messages
```bash
curl -X GET http://localhost:8080/api/messages
```

#### Create a New Message
```bash
curl -X POST http://localhost:8080/api/messages \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, World!",
    "phoneNumber": "+1234567890"
  }'
```

#### Check Scheduler Status
```bash
curl -X GET http://localhost:8080/api/messages/scheduler-status
```

#### Enable/Disable Message Sender
```bash
curl -X POST http://localhost:8080/api/messages/process-message-sender \
  -H "Content-Type: application/json" \
  -d '{
    "isMessageSenderEnabled": true
  }'
```

## 🛠️ Available Docker Commands

### Development Commands

```bash
# Start services in detached mode (background)
docker-compose up -d --build

# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis

# Stop all services
docker-compose down

# Stop services and remove volumes (⚠️ data will be lost)
docker-compose down -v
```

### Debugging Commands

```bash
# Check service status
docker-compose ps

# Access PostgreSQL database
docker-compose exec postgres psql -U postgres -d postgres

# Access Redis CLI
docker-compose exec redis redis-cli

# Execute commands in the app container
docker-compose exec app sh

# Rebuild without cache
docker-compose build --no-cache
```

## 📊 Services Configuration

### Application Service
- **Port**: 8080
- **Environment**: Development
- **Configuration**: Uses `config.docker.yaml`
- **Restart Policy**: Unless stopped

### PostgreSQL Database
- **Port**: 5432
- **Database**: postgres
- **Username**: postgres
- **Password**: postgres
- **Persistent Storage**: `postgres_data` volume
- **Auto Migration**: Yes (runs scripts from `migrations/` directory)

### Redis Cache
- **Port**: 6379
- **Persistent Storage**: `redis_data` volume
- **Configuration**: Append-only mode enabled for persistence

## 🔧 Configuration

The application uses different configuration files:
- **Local Development**: `configs/config.yaml`
- **Docker Environment**: `configs/config.docker.yaml`

Key configuration sections:
- Database connection settings
- Redis cache settings
- Webhook endpoint configuration
- Scheduler settings (interval, batch size, timeout)

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| GET | `/api/messages` | Retrieve all messages |
| POST | `/api/messages` | Create a new message |
| POST | `/api/messages/process-message-sender` | Enable/disable scheduler |
| GET | `/api/messages/scheduler-status` | Get scheduler status |
| GET | `/api/webhook-delivery/{messageId}` | Get webhook delivery record |
| GET | `` | Swagger documentation |

## 🚨 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :8080
   # Stop the service or change port in docker-compose.yml
   ```

2. **Database Connection Issues**
   ```bash
   # Check PostgreSQL logs
   docker-compose logs postgres
   # Verify database is accessible
   docker-compose exec postgres psql -U postgres -d postgres -c "SELECT 1;"
   ```

3. **Redis Connection Issues**
   ```bash
   # Check Redis logs
   docker-compose logs redis
   # Test Redis connection
   docker-compose exec redis redis-cli ping
   ```

4. **Application Won't Start**
   ```bash
   # Check application logs
   docker-compose logs app
   # Rebuild the application
   docker-compose up --build --force-recreate
   ```

### Reset Everything
```bash
# Stop all services and remove volumes (⚠️ all data will be lost)
docker-compose down -v

# Remove all containers and images
docker-compose down --rmi all

# Start fresh
docker-compose up --build
```