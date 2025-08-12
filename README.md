# 🌀 Wormhole

A secure, self-hosted stack-like file management server with client-side encryption support. Wormhole enables you to push and pop encrypted files and text in a LIFO (Last In, First Out) manner across different channels.

## ✨ Features

- **🔐 Client-Side Encryption**: Zero-knowledge architecture - files are encrypted on the client before upload
- **📚 Stack-Like Management**: LIFO-based file storage (push/pop operations)
- **🔑 User Authentication**: Secure JWT-based authentication with refresh tokens
- **📂 Channel System**: Organize your files in different channels with optional password protection
- **📄 Multi-Type Support**: Handle both text content and file uploads
- **🏠 Self-Hosted**: Complete control over your data
- **🚀 RESTful API**: Clean HTTP API for integration
- **💾 MySQL Backend**: Reliable database storage with automatic migrations

## 🏗️ Architecture

Wormhole follows a zero-knowledge principle where:
- All file encryption/decryption happens on the client side
- The server only stores encrypted blobs and metadata
- Encryption keys never leave the client
- The server cannot decrypt your data

## 🚀 Quick Start

### Prerequisites

- Go 1.24.4 or later
- MySQL database
- Basic understanding of REST APIs

### Installation

1. Clone the repository:
```bash
git clone https://github.com/darkard2003/wormhole.git
cd wormhole
```

2. Build the application:
```bash
make build
```

3. Set up your environment variables by creating a `.env` file:
```env
# Database Configuration
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=wormhole

# JWT Configuration
ACCESS_SECRET_KEY=your_access_secret_key
REFRESH_SECRET_KEY=your_refresh_secret_key

# Storage Configuration
STORE_PATH=/path/to/your/storage/directory
```

4. Run the server:
```bash
make run
```

The server will start on the default port (usually 8080) and automatically handle database migrations.

## � API Overview

Wormhole provides a RESTful API with the following main endpoints:

- **Authentication**: `/api/v1/signup`, `/api/v1/signin`, `/api/v1/refresh`
- **User Management**: `/api/v1/user/me`, `/api/v1/user/status`
- **Channel Operations**: `/api/v1/user/channels` (GET, POST, DELETE)
- **Item Operations**: `/api/v1/user/items` (GET for pop, POST for push)
- **Health Monitoring**: `/api/v1/health`, `/api/v1/ping`

For detailed API documentation and examples, check the `posting/` directory which contains sample requests in YAML format.

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Password Hashing**: Passwords are hashed using bcrypt
- **Channel Protection**: Optional password protection for channels
- **Client-Side Encryption**: Zero-knowledge architecture ensures server never sees unencrypted data
- **Secure Headers**: Proper HTTP security headers implementation

## 📁 Project Structure

```
wormhole/
├── cmd/server/              # Main application entry point
├── internals/
│   ├── models/             # Data models (User, Channel, Item, etc.)
│   ├── server/
│   │   ├── handlers/       # HTTP request handlers
│   │   └── middleware/     # Authentication middleware
│   ├── services/           # Business logic services
│   │   ├── db/            # Database interface and MySQL implementation
│   │   ├── jwtservice/    # JWT token management
│   │   ├── envservice/    # Environment variable management
│   │   └── storage_service/ # File storage abstraction
│   └── utils/             # Utility functions
├── migrations/            # Database migration files
└── posting/              # API testing examples
```

## 🛠️ Development

### Building
```bash
make build
```

### Running
```bash
make run
```

### Testing API Endpoints
The `posting/` directory contains example API requests in YAML format that can be used with tools like [Posting](https://github.com/darrenburns/posting) or similar HTTP clients.

To set up API testing:
1. Copy `posting/posting.env.template` to `posting/posting.env`
2. Update the credentials in `posting.env` with your test values
3. Use Posting or similar tools to run the test requests

See `posting/README.md` for detailed testing instructions.

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DB_USER` | MySQL database username | Yes |
| `DB_PASSWORD` | MySQL database password | Yes |
| `DB_HOST` | MySQL database host | Yes |
| `DB_PORT` | MySQL database port | Yes |
| `DB_NAME` | MySQL database name | Yes |
| `ACCESS_SECRET_KEY` | JWT access token secret | Yes |
| `REFRESH_SECRET_KEY` | JWT refresh token secret | Yes |
| `STORE_PATH` | File storage directory path | Yes |

### Database Setup

Wormhole automatically handles database migrations on startup. Ensure your MySQL instance is running and accessible with the provided credentials.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## 🔗 API Client Integration

Wormhole is designed to work with client applications that handle encryption/decryption. The recommended workflow:

1. **Client encrypts** files/text using your preferred encryption method
2. **Client sends** encrypted data with metadata (salt, IV, etc.) to Wormhole
3. **Wormhole stores** encrypted blobs without knowledge of content
4. **Client retrieves** encrypted data and decrypts locally

This ensures complete privacy and security of your data.

---

Built with ❤️ using Go, Gin, and MySQL.
