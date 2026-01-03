# Havoc

[![Go Report Card](https://goreportcard.com/badge/github.com/Jonathanthedeveloper/havoc)](https://goreportcard.com/report/github.com/Jonathanthedeveloper/havoc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)](https://github.com/Jonathanthedeveloper/havoc/releases)

Havoc is a powerful CLI tool designed for **chaos engineering** and **traffic inspection**. It helps developers test their applications' resilience against network instability and proxy configurations, offering controllable latency, jitter, and packet drops.

## Features

- **Traffic Proxy**: Inspect and intercept HTTP/HTTPS traffic.
- **Chaos Injection**: Simulate meaningful network conditions:
  - **Latency**: Add fixed or random delays.
  - **Jitter**: vary the latency to simulate unstable connections.
  - **Packet Drops**: Randomly drop requests to test failure handling.
- **Easy Configuration**: Simple CLI flags for all settings.

## Installation

### From Source

```bash
go install github.com/Jonathanthedeveloper/havoc@latest
```

### Manual Build

```bash
git clone https://github.com/Jonathanthedeveloper/havoc.git
cd havoc
go build -o havoc
```

## Usage

Check the version:
```bash
havoc --version
```

Start the proxy server (example):
```bash
havoc start --port 8080 --target http://localhost:3000
```

Get help:
```bash
havoc --help
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a Pull Request.

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Jonathan** - [Jonathanthedeveloper](https://github.com/Jonathanthedeveloper)
