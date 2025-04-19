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

## 📚 Project Structure

```
smartupdater/
├── cmd/
│   ├── migrate/    # Database migration tool
│   └── server/     # Main application server
├── internal/
│   ├── api/        # REST API handlers
│   ├── config/     # Configuration management
│   ├── models/     # Database models
│   ├── services/   # Core business logic
│   └── tests/      # Test suites
└── docs/           # Documentation
```

## 📝 License

This project is licensed under the MIT License. 