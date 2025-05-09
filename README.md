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
# Clone the repository
git clone https://github.com/busybytelab.com/glimmer.git
cd glimmer

# Build and run
make build

# Edit the environment variables
cp .env.example .env

# Run the app, e.g. on Linux arm64
./build/glimmer-linux-arm64

# Or using docker-compose
docker compose up -d
```

For detailed installation instructions and configuration options, see our [Installation Guide](docs/installation.md).

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! This is a personal project in its early stages.

## Status

Glimmer is currently in early development. Features and documentation are actively being added.
