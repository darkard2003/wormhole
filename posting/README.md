# 📡 API Testing with Posting

This directory contains API testing examples for Wormhole using [Posting](https://github.com/darrenburns/posting) YAML format.

## 🚀 Quick Setup

1. **Install Posting** (if you haven't already):
   ```bash
   pip install posting
   ```

2. **Configure test credentials**:
   - Copy the `posting.env` file and update with your test credentials
   - The file contains placeholder values that you need to replace:
     ```env
     BASE_URL=http://localhost:8080/api/v1
     BASE_USERNAME=your_test_username
     BASE_PASSWORD=your_test_password
     ```

3. **Start your Wormhole server**:
   ```bash
   cd .. && make run
   ```

4. **Run the tests** using Posting UI or command line

## 📁 Test Structure

```
posting/
├── posting.env              # Environment variables (configure this!)
├── auth/                    # Authentication endpoints
│   ├── signup.posting.yaml
│   └── signin.posting.yaml
├── channels/                # Channel management
│   ├── create-channel.posting.yaml
│   └── get-channels.posting.yaml
├── items/                   # Item push/pop operations
│   ├── push-item.posting.yaml
│   └── pop-item.posting.yaml
├── user/                    # User management
│   └── me.posting.yaml
└── scripts/                 # Helper scripts for token management
    ├── add_token_to_request.py
    └── store_auth_token.py
```

## 🔄 Testing Workflow

1. **Sign up** a test user (`auth/signup.posting.yaml`)
2. **Sign in** to get authentication token (`auth/signin.posting.yaml`)
3. **Create channels** for organizing your items (`channels/create-channel.posting.yaml`)
4. **Push items** (text or files) to channels (`items/push-item.posting.yaml`)
5. **Pop items** from channels (`items/pop-item.posting.yaml`)
6. **Check user details** (`user/me.posting.yaml`)

## 🔐 Security Notes

- **Never commit real credentials** to version control
- The `posting.env` file is gitignored to prevent accidental exposure
- Use test data only - this is for development and testing
- Scripts automatically handle JWT token management

## ⚙️ Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `BASE_URL` | Wormhole API base URL | `http://localhost:8080/api/v1` |
| `BASE_USERNAME` | Test username for authentication | `testuser` |
| `BASE_PASSWORD` | Test password for authentication | `testpassword123` |

## 🛠️ Custom Scripts

- **`add_token_to_request.py`**: Automatically adds JWT bearer token to authenticated requests
- **`store_auth_token.py`**: Extracts and stores JWT token from signin response

These scripts work together to handle the authentication flow automatically.
