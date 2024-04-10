# SmartUpdater

SmartUpdater is an intelligent dependency update scheduler that analyzes GitHub commit patterns and CI success rates to determine optimal update windows.

## 🌟 Features

- Data-driven dependency update scheduling
- GitHub API integration for real-time analytics
- Time-series analysis for trend detection
- REST API for monitoring and control
- Prometheus/Grafana metrics integration
- PostgreSQL for data persistence

## 🚀 Setup

### Prerequisites

1. Go 1.21 or higher
2. PostgreSQL database
3. GitHub API token
4. Docker (optional, for containerization)

### Environment Variables

Create a `.env` file in the project root:

```env
GITHUB_TOKEN=your_github_token
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smartupdater
PORT=8080
```

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/smartupdater.git
   cd smartupdater
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   ```bash
   psql -U postgres
   CREATE DATABASE smartupdater;
   ```

4. Run migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

### Running the Application

1. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

2. The server will start on `http://localhost:8080`

## 📚 API Documentation

### Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/metrics` - Prometheus metrics
- `GET /api/v1/stats` - Update statistics
- `POST /api/v1/repositories` - Add repository for monitoring
- `GET /api/v1/repositories` - List monitored repositories
- `GET /api/v1/schedule` - Get update schedule
- `POST /api/v1/schedule` - Modify update schedule

## 📊 Monitoring

Access metrics dashboard:
1. Prometheus: `http://localhost:9090`
2. Grafana: `http://localhost:3000`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 