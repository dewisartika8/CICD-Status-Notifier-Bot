# Dashboard Frontend README

# Dashboard Frontend

This project is a frontend dashboard for monitoring the CI/CD Status Notifier Bot backend. It provides a user-friendly interface to visualize project metrics, build statuses, and analytics through a set of interactive components.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Components](#components)
- [API Integration](#api-integration)
- [WebSocket Integration](#websocket-integration)
- [Contributing](#contributing)
- [License](#license)

## Installation

To get started with the dashboard frontend, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd dashboard-frontend
   ```

2. Install the dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

4. Open your browser and navigate to `http://localhost:3000` to view the dashboard.

## Usage

The dashboard allows users to:

- View an overview of key metrics and statuses.
- Access detailed project information and build history.
- Analyze build success rates and durations through visual charts.

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── dashboard/
│   │   ├── charts/
│   │   ├── common/
│   │   └── layout/
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Projects.tsx
│   │   ├── Analytics.tsx
│   │   └── Settings.tsx
│   ├── services/
│   │   ├── api.ts
│   │   └── websocket.ts
│   ├── store/
│   │   └── index.ts
│   └── App.tsx
└── package.json
```

## Components

The dashboard is built using the following key components:

- **Layout**: Provides a common structure including header and sidebar.
- **Header**: Displays the title and navigation options.
- **Sidebar**: Contains links to different sections of the dashboard.
- **Overview**: Summarizes key metrics and statuses.
- **ProjectList**: Lists all projects with their current statuses.
- **MetricsCard**: Displays individual metrics in a card format.
- **SuccessRateChart**: Visualizes the success rate of builds over time.
- **BuildDurationChart**: Visualizes the average build duration over time.

## API Integration

The frontend communicates with the backend through the following API endpoints:

- `GET /api/v1/dashboard/overview`: Fetches the overview metrics.
- `GET /api/v1/dashboard/projects`: Retrieves the list of projects.
- `GET /api/v1/projects/:id/details`: Gets detailed information about a specific project.
- `GET /api/v1/projects/:id/builds`: Fetches the build history for a specific project.
- `GET /api/v1/dashboard/projects/:id/metrics`: Fetches specific metrics for a project.
- `GET /api/v1/dashboard/notifications/stats`: Retrieves notification statistics.
- `GET /api/v1/events/stream`: Establishes a server-sent events stream for real-time updates.
- `WS /ws/events`: Establishes a WebSocket connection for real-time events.
- `GET /api/v1/analytics/trends`: Fetches analytics trends data.
- `GET /api/v1/analytics/reports`: Retrieves analytics reports.
- `POST /api/v1/analytics/export`: Exports analytics data.
- `GET /api/v1/subscriptions`: Lists all subscriptions.
- `POST /api/v1/subscriptions`: Creates a new subscription.
- `PUT /api/v1/subscriptions/:id`: Updates an existing subscription.
- `DELETE /api/v1/subscriptions/:id`: Deletes a subscription.

Ensure that the backend is running and accessible for the frontend to function correctly.

## WebSocket Integration

The dashboard supports real-time updates through WebSocket connections. This allows users to receive live notifications about project status changes and build events.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your branch and create a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.