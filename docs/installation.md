# Installation Guide

This guide covers different methods to install and run Glimmer.

## Local Build

To build and run Glimmer locally:

1. Ensure you have Go 1.21 or later installed
2. Clone the repository:
   ```bash
   git clone https://github.com/busybytelab.com/glimmer.git
   cd glimmer
   ```
3. Build the application:
   ```bash
   make build
   ```
4. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```
5. Run the application:
   ```bash
   # For Linux arm64
   ./build/glimmer-linux-arm64
   # For other platforms, use the appropriate binary from ./build/
   ```

## Docker Compose

The recommended way to run Glimmer is using Docker Compose. This method provides a consistent environment and easy configuration.

### Prerequisites

- Docker
- Docker Compose

### Basic Setup

1. Create a `.env` file with the following variables:
   ```env
   # Required
   ENCRYPTION_KEY=your32characterkeyhere12345678901  # Generate with: openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32 && echo
   
   # Optional
   LISTEN_ADDRESS=0.0.0.0:8787  # Default: 0.0.0.0:8787
   ADMIN_EMAIL=admin@dev.local  # Default: admin@dev.local
   ADMIN_PASSWORD=your-admin-password  # Default: glimmerglimmer
   ```

2. Run the application:
   ```bash
   docker compose up -d
   ```

### Demo Profile

To run Glimmer with demo data:

1. Set the required environment variables:
   ```env
   SEED_DATA=true
   COMPOSE_PROFILES=demo
   ```

2. Run with demo profile:
   ```bash
   docker compose --profile demo up -d
   ```

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENCRYPTION_KEY` | 32-character encryption key | - | Yes |
| `LISTEN_ADDRESS` | Address to listen on | 0.0.0.0:8787 | No |
| `ADMIN_EMAIL` | Admin user email | admin@dev.local | No |
| `ADMIN_PASSWORD` | Admin user password | glimmerglimmer | No |
| `SEED_DATA` | Enable demo data seeding | false | No |
| `COMPOSE_PROFILES` | Docker Compose profiles to enable | demo | No |

## Portainer Stack

To deploy Glimmer using Portainer:

1. Log in to your Portainer instance
2. Navigate to Stacks
3. Click "Add stack"
4. Name your stack (e.g., "glimmer")
5. Choose "Web editor" as the build method
6. Copy the contents of `docker-compose.yaml` into the editor
7. Create a new environment file (stack.env) with your configuration:
   ```env
   ENCRYPTION_KEY=your32characterkeyhere12345678901
   LISTEN_ADDRESS=0.0.0.0:8787
   ADMIN_EMAIL=your-email@example.com
   ADMIN_PASSWORD=your-secure-password
   ```
8. For demo data, add:
   ```env
   SEED_DATA=true
   COMPOSE_PROFILES=demo
   ```
9. Click "Deploy the stack"

### Portainer Environment Variables

When using Portainer, you can set environment variables in two ways:

1. **Stack Environment File**: Create a `stack.env` file in your stack configuration
2. **Environment Variables**: Add them directly in the Portainer stack configuration

### Persistent Data

By default, Glimmer stores its data in `/opt/glimmer-deploy/pb_data`. Ensure this directory exists and has proper permissions:

```bash
sudo mkdir -p /opt/glimmer-deploy/pb_data
sudo chown -R 1000:1000 /opt/glimmer-deploy/pb_data
```

You can change it in docker-compose.yaml by changing the volume path.

## Next Steps

After installation:
1. Access the superuser interface at `http://your-server:8787/_`
2. Log in with your superuser credentials (ADMIN_EMAIL and ADMIN_PASSWORD)
3. Configure your instance through the web interface

To login as a user, need to create user with correct roles (or use demo data)
1. Go to the normal user interface at `http://your-server:8787`
2. Login with the user credentials (u1@home.local / password123 for demo)
