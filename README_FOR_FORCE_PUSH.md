# Yggdrasil API Go Project

This is a Yggdrasil API server implementation in Go, compatible with Minecraft authentication and skin serving.

## Features

- ✅ Yggdrasil protocol compliance
- 🔐 Secure authentication system
- 👥 User management with profile features
- 🎨 Skin and texture serving
- 🛡️ Rate limiting and security measures
- 📊 Performance monitoring
- 🎮 Player name registration system
- 🔔 Announcement system
- 🗂️ Multiple storage backends support

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/httye/yggdrasil-skins-go.git
   ```

2. Navigate to the project directory:
   ```bash
   cd yggdrasil-skins-go
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   go build -o yggdrasil-api-go
   ```

## Configuration

1. Copy the example configuration file:
   ```bash
   cp conf/example.yml conf/config.yml
   ```

2. Edit `conf/config.yml` with your settings

3. Generate RSA key pairs:
   ```bash
   make keys
   ```

## Running the Server

Use the provided Makefile:

```bash
# Build the project
make build

# Run the server
make run

# Or run in development mode
make dev
```

## License

This project is licensed under the terms specified in the LICENSE file.
