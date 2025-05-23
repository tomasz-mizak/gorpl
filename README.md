# Medicinal Products Data Fetcher

A simple application for fetching and storing medicinal products data from the Polish e-Health Registry.

## Overview

This application downloads medicinal products data from the Polish e-Health Registry API and stores it locally. It provides a simple interface to view and search through the stored data.

## Features

- Downloads medicinal products data from the Polish e-Health Registry API
- Stores data locally for offline access
- Simple interface for viewing and searching data
- Special API interface for Unitbox integration
- Error handling with fallback to previously downloaded data

## Data Source

The application fetches data from the following endpoint:
```
https://rejestry.ezdrowie.gov.pl/api/rpl/medicinal-products/public-pl-report/6.0.0/overall.xml
```

## Usage Notes

- The application downloads data only once and stores it locally
- No automatic updates are performed
- If a download fails, the application will use the previously downloaded data
- An error will be thrown if no data is available locally

## Technical Details

- The application is designed to be simple and lightweight
- Data is stored locally to minimize API calls
- Error handling ensures graceful degradation when network issues occur

## API Integration

A special API interface is available for Unitbox integration. Please refer to the API documentation for details.

## Docker Deployment

The application can be easily deployed using Docker. The following files are provided:

- `Dockerfile` - Multi-stage build configuration
- `docker-compose.yml` - Docker Compose configuration
- `.dockerignore` - Excludes unnecessary files from the build

### Quick Start

1. Build and start the container:
   ```bash
   docker-compose up --build
   ```

2. Access the application at `http://localhost:1532`

### Docker Compose Commands

- Start in detached mode:
  ```bash
  docker-compose up -d --build
  ```

- Stop the container:
  ```bash
  docker-compose down
  ```

### Configuration

- The application runs on port 1532
- Data is persisted in the `./data` directory
- Timezone is set to Europe/Warsaw
- Automatic restart is enabled

### Publishing Docker Image

To publish the Docker image to a registry (e.g., Docker Hub, GitHub Container Registry), follow these steps:

1. Build the image with a tag:
   ```bash
   docker build -t yourusername/gorpl:latest .
   ```

2. Log in to your registry:
   ```bash
   docker login
   ```

3. Push the image:
   ```bash
   docker push yourusername/gorpl:latest
   ```

#### Using GitHub Container Registry (GHCR)

1. Log in to GHCR:
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
   ```

2. Build and tag the image:
   ```bash
   docker build -t ghcr.io/yourusername/gorpl:latest .
   ```

3. Push the image:
   ```bash
   docker push ghcr.io/yourusername/gorpl:latest
   ```

Note: Replace `yourusername` with your actual username and ensure you have the necessary permissions to push to the registry.