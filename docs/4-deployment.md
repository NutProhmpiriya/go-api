# Deployment Guide

## Overview

This guide covers the deployment process for the Vongga Platform backend service. We support multiple deployment options including Docker containers and traditional deployment.

## Prerequisites

- Docker and Docker Compose for containerized deployment
- Access to a PostgreSQL database
- Access to a Redis instance
- SSL certificates for HTTPS (production)
- Firebase configuration

## Environment Configuration

### Required Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080
ENV=production
API_SECRET=your-secret-key

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=vongga
DB_USER=postgres
DB_PASSWORD=your-password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# Firebase Configuration
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_PRIVATE_KEY=your-private-key
FIREBASE_CLIENT_EMAIL=your-client-email
```

## Deployment Options

### 1. Docker Deployment

1. **Build the Docker image**
   ```bash
   docker build -t vongga-backend .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Verify deployment**
   ```bash
   docker ps
   docker logs vongga-backend
   ```

### 2. Traditional Deployment

1. **Build the application**
   ```bash
   go build -o app cmd/api/main.go
   ```

2. **Run database migrations**
   ```bash
   make migrate-up
   ```

3. **Start the application**
   ```bash
   ./app
   ```

## Health Checks

Monitor the following endpoints:

- `/health`: Basic health check
- `/metrics`: Prometheus metrics
- `/debug/pprof`: Performance profiling (development only)

## Monitoring

### Metrics Collection

1. **Prometheus Integration**
   - Metrics available at `/metrics`
   - Configure Prometheus to scrape this endpoint

2. **Logging**
   - Logs are written to stdout/stderr
   - Use your platform's log aggregation service

### Alert Configuration

Set up alerts for:
- High error rates (> 1%)
- High latency (p99 > 500ms)
- High CPU usage (> 80%)
- High memory usage (> 80%)

## Backup and Recovery

### Database Backup

1. **Automated Backups**
   ```bash
   # Daily backup script
   pg_dump -U postgres vongga > backup_$(date +%Y%m%d).sql
   ```

2. **Backup Verification**
   - Regularly test backup restoration
   - Maintain backup rotation policy

### Recovery Procedures

1. **Database Recovery**
   ```bash
   psql -U postgres vongga < backup_file.sql
   ```

2. **Application Recovery**
   - Keep multiple versions for rollback
   - Document rollback procedures

## SSL/TLS Configuration

For production deployments:

1. **Generate SSL Certificate**
   - Use Let's Encrypt or your SSL provider
   - Place certificates in `/etc/ssl/certs/`

2. **Configure HTTPS**
   - Update nginx/reverse proxy configuration
   - Enable HTTP/2
   - Configure SSL termination

## Security Considerations

1. **Firewall Configuration**
   - Allow only necessary ports
   - Implement rate limiting
   - Configure DDoS protection

2. **Access Control**
   - Use least privilege principle
   - Regularly rotate credentials
   - Monitor access logs

## Troubleshooting

Common issues and solutions:

1. **Database Connection Issues**
   - Check network connectivity
   - Verify credentials
   - Check connection limits

2. **Performance Issues**
   - Monitor resource usage
   - Check slow queries
   - Review application logs

3. **Memory Leaks**
   - Use pprof for profiling
   - Monitor memory usage
   - Check for goroutine leaks
