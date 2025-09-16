# Order Management System (OMS)

An order management system built with Go, PostgreSQL, containerized with Docker.

## Prerequisites

- Docker and Docker Compose installed
- Git (for cloning the repository)

## Quick Start

### 1. Clone and Setup
```bash
git clone <repository-url>
cd order-management-system
```

### 2. Build and Run Everything
```bash
# Build and start all services
docker-compose up --build

# Run in background (detached mode)
docker-compose up --build -d
```

### 3. Verify Installation
```bash
# Check if all services are running
docker-compose ps

# Test the application
curl http://localhost:8089/health

# Check application logs
docker-compose logs -f app
```

## Service Management

### Individual Service Operations
```bash
# Start only specific services
docker-compose up postgres redis

# Rebuild only the app service
docker-compose build app
docker-compose up app

# View service logs
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis

# Restart a specific service
docker-compose restart app
```

### Service Status
```bash
# Check running containers
docker-compose ps

# View resource usage
docker stats $(docker-compose ps -q)
```

## Database Management

### Basic Operations
```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U raisul -d order_management_system

# Run SQL commands directly
docker-compose exec postgres psql -U raisul -d order_management_system -c "SELECT * FROM stores LIMIT 5;"

# List all tables
docker-compose exec postgres psql -U raisul -d order_management_system -c "\dt"
```

### Data Management
```bash
# Backup database
docker-compose exec postgres pg_dump -U raisul order_management_system > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore from backup
docker-compose exec -T postgres psql -U raisul -d order_management_system < backup.sql

# View database size
docker-compose exec postgres psql -U raisul -d order_management_system -c "SELECT pg_size_pretty(pg_database_size('order_management_system'));"
```

### Reset Database
```bash
# Stop services
docker-compose down

# Remove database volume (WARNING: This deletes all data)
docker volume rm order-management-system_postgres_data

# Restart services (will recreate database)
docker-compose up --build
```

## Development Workflow

### Code Changes
```bash
# After making code changes
docker-compose build app && docker-compose up app

# Or rebuild everything
docker-compose up --build
```

### Clean Development Environment
```bash
# Stop all services
docker-compose down

# Remove containers, networks, and anonymous volumes
docker-compose down --volumes --remove-orphans

# Remove all unused Docker resources
docker system prune -f
```

## API Testing

### Sample API Calls
```bash
# Health check
curl http://localhost:8089/health

# List all orders (requires JWT token)
curl --location 'http://localhost:8089/api/v1/orders/all' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'

# Create a new order
curl --location 'http://localhost:8089/api/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--data '{
    "store_id": 1,
    "recipient_name": "John Doe",
    "recipient_phone": "01712345678",
    "recipient_address": "123 Test Street, Dhaka",
    "recipient_city": 1,
    "recipient_zone": 1,
    "delivery_type": 1,
    "item_type": 1,
    "item_quantity": 1,
    "item_weight": 0.5,
    "amount_to_collect": 1000
}'
```

### Authentication
```bash
# Login to get JWT token
curl --location 'http://localhost:8089/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "email": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

## Configuration

### Environment Variables
The system uses the following configuration through environment variables:

- **Database**: PostgreSQL connection settings
- **Redis**: Cache server configuration
- **JWT**: Token expiration settings
- **Application**: Port and other app settings

### Port Mappings
- **Application**: `localhost:8089` → `container:8089`
- **PostgreSQL**: `localhost:5432` → `container:5432`
- **Redis**: `localhost:6379` → `container:6379`

### Data Persistence
- **PostgreSQL data**: `postgres_data` volume
- **Redis data**: `redis_data` volume

Data persists between container restarts unless volumes are explicitly removed.

## Monitoring and Health Checks

### Service Health
```bash
# Check PostgreSQL health
docker-compose exec postgres pg_isready -U raisul

# Check Redis connectivity
docker-compose exec redis redis-cli ping

# View application metrics
curl http://localhost:8089/metrics
```

### Logs and Debugging
```bash
# Follow all logs
docker-compose logs -f

# View logs for specific time range
docker-compose logs --since "2024-01-01T00:00:00Z" --until "2024-01-02T00:00:00Z"

# Export logs to file
docker-compose logs > oms_logs_$(date +%Y%m%d_%H%M%S).log
```

## Troubleshooting

### Application Issues
```bash
# Check if app container is running
docker-compose ps app

# View application logs
docker-compose logs app

# Execute commands inside app container
docker-compose exec app sh
```

### Database Connection Issues
```bash
# Test PostgreSQL connectivity
docker-compose exec postgres pg_isready -U raisul

# Check database exists
docker-compose exec postgres psql -U raisul -l

# Verify network connectivity from app to database
docker-compose exec app ping postgres
```

### Redis Connection Issues
```bash
# Test Redis connectivity
docker-compose exec redis redis-cli ping

# Check Redis configuration
docker-compose exec redis redis-cli config get "*"

# Monitor Redis operations
docker-compose exec redis redis-cli monitor
```

### Performance Issues
```bash
# Check container resource usage
docker stats

# View system resource usage
docker system df

# Check database performance
docker-compose exec postgres psql -U raisul -d order_management_system -c "SELECT * FROM pg_stat_activity;"
```

### Common Error Solutions

**Error: Port already in use**
```bash
# Find process using the port
lsof -i :8089
# Kill the process or use different port
```

**Error: Database connection refused**
```bash
# Ensure PostgreSQL is healthy
docker-compose logs postgres
# Restart database service
docker-compose restart postgres
```

**Error: Out of disk space**
```bash
# Clean up Docker resources
docker system prune -a
# Remove unused volumes
docker volume prune
```

#### Future Scope of Work
- Implement caching for holding sessions and other frequently used data