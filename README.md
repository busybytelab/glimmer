> ⚠️ **Status**: This project is in early development and not ready yet. Features, documentation, and stability are actively being improved.

# Glimmer

A simple, self-hosted app helping kids (~6-12 yrs) practice specific skills like grammar or basic math. Parents define the topics, an LLM generates the content. Focused on easy setup (Docker) and providing safe, supplementary practice.

## Features

- **Simple Setup**: Easy to run with Docker
- **Parent-Configured**: Create practice topics tailored to your child's needs
- **LLM-Powered**: Generates varied practice content automatically
- **Privacy-Focused**: Self-hosted for full control of your data
- **Child-Friendly UI**: Simple interface designed for tablets and computers

## Usage

Glimmer allows parents to:
1. Set up practice topics (like punctuation exercises or math problems)
2. Define parameters (age level, difficulty, etc.)
3. Let the LLM generate appropriate practice content

Children can then:
1. Access the web interface from their device
2. Select from available practice activities
3. Complete exercises and receive immediate feedback

## Quick Start

```bash
# Build the Docker image
docker build -t glimmer .

# or pull the image from Docker Hub
docker pull ghcr.io/busybytelab/glimmer:latest

# Run with Docker
docker run -p 8787:8787 glimmer

# Or using docker-compose
docker-compose up
```

For detailed installation instructions, see the [installation guide](docs/installation.md).

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! This is a personal project in its early stages.

## Status

Glimmer is currently in early development. Features and documentation are actively being added.
