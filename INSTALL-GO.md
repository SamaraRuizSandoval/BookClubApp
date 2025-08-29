# 📥 Go Installation on Windows

This guide will help you install Go on Windows to run the BookClubApp HTTP server.

## 🎯 Prerequisites

- Windows 10/11
- Administrator privileges
- Internet connection

## 🚀 Method 1: Installation via Official Installer (Recommended)

### Step 1: Download Go
1. Go to [https://go.dev/dl/](https://go.dev/dl/)
2. Click the download button for Windows
3. Choose version **Go 1.23.3** (or latest)

### Step 2: Install Go
1. Run the downloaded `.msi` file
2. Follow the installation wizard
3. Keep the default directory: `C:\Go`
4. Complete the installation

### Step 3: Verify Installation
1. Close all open PowerShell/Command Prompt windows
2. Open a **new** PowerShell window
3. You should see something like:
   ```
   go version go1.23.3 windows/amd64
   ```

## 🔧 Method 2: Installation via Chocolatey

If you already have Chocolatey installed:

```powershell
# Install Go
choco install golang

# Verify installation
go version
```

## 🔧 Method 3: Installation via Scoop

If you already have Scoop installed:

```powershell
# Install Go
scoop install go

# Verify installation
go version
```

## 🌍 Environment Variables Configuration

The Go installer should automatically configure the environment variables:

- `GOROOT`: Points to Go installation directory
- `GOPATH`: Points to your Go workspace
- `Path`: Includes Go binary directory

### If not configured, configure manually:

1. Press `Win + R`, type `sysdm.cpl`, press Enter
2. Go to "Advanced" tab
3. Click "Environment Variables"
4. In "System Variables", click "New"
5. Add Go root:
   - **Variable name**: `GOROOT`
   - **Variable value**: `C:\Go`
6. Add Go path:
   - **Variable name**: `GOPATH`
   - **Variable value**: `%USERPROFILE%\go`
7. In "System Variables", find `Path`
8. Click "Edit" and add:
   - `%GOROOT%\bin`
   - `%GOPATH%\bin`

## ✅ Installation Verification

After installation, open a **new** PowerShell and run:

```powershell
# Check Go version
go version

# Check environment variables
go env

# Check if Go is working
go help
```

## 🚀 Running BookClubApp

Now you can run the server:

```powershell
# Navigate to project
cd BookClubApp

# Install dependencies
go mod tidy

# Run server
go run main.go
```

## 🔍 Troubleshooting

### Error: "go is not recognized as a command"
- **Solution**: Restart PowerShell or restart computer
- **Check**: Environment variables are configured correctly

### Error: "go: command not found"
- **Solution**: Check if Go was installed in the correct directory
- **Check**: Run `go env` to see configurations

### Permission error
- **Solution**: Run PowerShell as administrator
- **Alternative**: Install Go in a user directory

## 📚 Additional Resources

- [Official Go Documentation](https://golang.org/doc/)
- [Go Installation Guide](https://golang.org/doc/install)
- [Go Environment Variables](https://golang.org/doc/install#environment)

## 🆘 Still having problems?

1. Restart your computer
2. Check if Go is in the correct directory
3. Check environment variables
4. Consult the [official documentation](https://golang.org/doc/install)

---

**🎉 Congratulations! Now you can run the BookClubApp HTTP server!**
