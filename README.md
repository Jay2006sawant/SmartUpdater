# SmartUpdater

SmartUpdater is an intelligent dependency update scheduler that optimizes the timing and execution of dependency updates in GitHub repositories. By analyzing historical commit patterns, CI pipeline success rates, and repository activity, it determines the optimal windows for performing updates, reducing the risk of failed updates and minimizing disruption to development workflows.

## 🎯 Problem Solved

Dependency updates are crucial for security and feature enhancements, but they often fail due to:
- Poor timing (e.g., during active development)
- Incompatible with current CI pipeline state
- Lack of historical success pattern analysis
- Manual scheduling leading to human error

SmartUpdater addresses these challenges by:
- Analyzing repository activity patterns to find quiet periods
- Tracking CI pipeline success rates to identify stable windows
- Automating the update process with intelligent scheduling
- Providing data-driven insights for update timing
- Creating pull requests with detailed update information

## 💡 Benefits

- **Reduced Update Failures**: By analyzing historical patterns and CI success rates, SmartUpdater significantly reduces the likelihood of failed updates.
- **Time Optimization**: Automatically identifies the best times for updates based on repository activity and CI pipeline health.
- **Developer Productivity**: Minimizes disruption to development workflows by scheduling updates during optimal periods.
- **Security Enhancement**: Ensures timely dependency updates while maintaining system stability.
- **Data-Driven Decisions**: Provides insights into update patterns and success rates for better decision-making.

## 🎯 Use Cases

- **Enterprise Development**: Manage dependency updates across multiple repositories with minimal disruption.
- **Open Source Projects**: Automate dependency maintenance while respecting contributor activity patterns.
- **CI/CD Pipelines**: Integrate with existing CI/CD workflows to optimize update timing.
- **Security Teams**: Ensure timely security updates while maintaining system stability.
- **DevOps Teams**: Reduce manual intervention in dependency management processes.

## 🌟 Features

- GitHub API integration for repository monitoring
- Automated dependency update scheduling
- Pull request creation for dependency updates
- REST API for repository management
- PostgreSQL for data persistence
- Docker support for easy deployment

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

### Running Locally

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

### Running with Docker

1. Build and start the containers:
   ```bash
   docker-compose up --build
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
│   └── services/   # Core business logic
├── Dockerfile      # Container definition
└── docker-compose.yml  # Container orchestration
```

## 🔌 API Endpoints

- `GET /repositories` - List all repositories
- `POST /repositories` - Add a new repository
- `GET /repositories/{id}` - Get repository details
- `PUT /repositories/{id}` - Update repository settings
- `DELETE /repositories/{id}` - Remove a repository
- `GET /repositories/{id}/updates` - List update history

## 📝 License

This project is licensed under the MIT License. 

