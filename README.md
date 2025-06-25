<p align="center">
  <img src="./static/logo.png" alt="Keygeist Logo" width="220"/>
</p>

<h3 align="center"><em>Your keyboard. Haunted by intelligence.</em></h3>

<div align="center">
  
[![Go Report Card](https://goreportcard.com/badge/github.com/mudler/Keygeist)](https://goreportcard.com/report/github.com/mudler/Keygeist)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/mudler/Keygeist)](https://github.com/mudler/Keygeist/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/mudler/Keygeist)](https://github.com/mudler/Keygeist/issues)

</div>

An invisible AI assistant that responds to your key combos, bringing thoughts to life — one keystroke at a time.

Keygeist is an AI-powered keyboard operator that listens for key combinations and responds with AI-generated text typed directly into your Linux box. Works with any DE and application, it emulates a virtual keyboard.

## Why?

I couldn't find a simple tool that does exactly this for Linux, working with local models - takes a key combination and optionally captures the current context (screen, clipboard) and then takes over the keyboard with LLM output. Many times I just wanted to speed up some interactions, and the easiest way would be for the LLM to type with some directions from my side.

## Features

- **AI-Powered Assistant**: Responds to key combinations with AI-generated text
- **Contextual AI**: Choose what context to send to the LLM (clipboard, screenshot, or both)
- **Real-time Interaction**: Non-blocking interactions that run in the background
- **Multiple Context Modes**: Different key combinations for different types of context
- **Interaction Control**: Cancel ongoing interactions with the same key combination

## Keygeist

<p align="center">
  <img src="./static/draw.png" alt="Keygeist Architecture" width="550"/>
</p>

Keygeist is an AI-powered assistant that:
1. Listens for specific key combinations
2. Opens a zenity dialog for user input
3. Sends the input (with optional clipboard and/or screenshot context) to OpenAI API
4. Types the AI response using the keyboard emulator

### Key Combinations

By default, Keygeist uses the following key combinations:
- **Windows + C**: Sends only clipboard content as context
- **Windows + S**: Sends only a screenshot as context  
- **Windows + E**: Sends both clipboard and screenshot as context
- **Windows + T**: Sends no additional context (text-only mode)

You can customize these keybindings using environment variables.

#### Customizing Keybindings

You can customize the keybindings using environment variables:

```bash
export CLIPBOARD_KEY="ctrl+shift+c"
export SCREENSHOT_KEY="ctrl+shift+s"
export ALL_CONTEXT_KEY="ctrl+shift+e"
export TEXT_ONLY_KEY="ctrl+shift+t"
./build/keygeist
```

**Supported Key Formats**

Key combinations use the format `modifier+key` where:
- **Modifiers**: `ctrl`, `alt`, `shift`, `win` (or `windows`, `meta`)
- **Keys**: `a-z`, `0-9`, `f1-f12`, `space`, `enter`, `tab`, `escape`, `backspace`

Examples:
- `win+c` (Windows + C)
- `ctrl+shift+c` (Ctrl + Shift + C)
- `alt+f4` (Alt + F4)
- `ctrl+alt+delete` (Ctrl + Alt + Delete)

**Note**: Press the same combination again to cancel the current interaction

### Features

- **Contextual AI**: Choose what context to send to the LLM (clipboard, screenshot, or both)
- **Interaction Control**: Press the same combination again to stop the current interaction
- **Non-blocking**: Interactions run in the background, allowing you to continue using your system
- **Error Handling**: Graceful handling of API errors and user cancellations

### Prerequisites

- `zenity` installed on your system
- OpenAI API key
- Root privileges (for uinput access) or udev rules setup (see Installation section)

### Environment Variables

- `OPENAI_API_KEY` (required): Your OpenAI API key
- `OPENAI_MODEL` (required): The OpenAI model to use (e.g., `gpt-3.5-turbo`, `gpt-4`)
- `OPENAI_BASE_URL` (optional): Custom base URL for OpenAI API (useful for proxies or alternative endpoints)
- `OPENAI_SYSTEM_PROMPT` (optional): Custom system prompt to use when querying the LLM. If not set, a default helpful assistant prompt will be used.
- `CLIPBOARD_KEY` (optional): Custom key combination for clipboard context (e.g., `ctrl+shift+c`)
- `SCREENSHOT_KEY` (optional): Custom key combination for screenshot context (e.g., `ctrl+shift+s`)
- `ALL_CONTEXT_KEY` (optional): Custom key combination for all context (e.g., `ctrl+shift+e`)
- `TEXT_ONLY_KEY` (optional): Custom key combination for text-only context (e.g., `ctrl+shift+t`)

### Usage

```bash
# Set environment variables
export OPENAI_API_KEY="your_api_key_here"
export OPENAI_MODEL="gpt-3.5-turbo"
export OPENAI_BASE_URL="https://api.openai.com/v1"  # Optional
export OPENAI_SYSTEM_PROMPT="You are a coding assistant. Provide concise, practical code solutions."  # Optional

# Build and run with default keybindings (with sudo)
make build
sudo -E ./build/keygeist

# Run with custom keybindings via environment variables
export CLIPBOARD_KEY="ctrl+shift+c"
export SCREENSHOT_KEY="ctrl+shift+s"
export ALL_CONTEXT_KEY="ctrl+shift+e"
export TEXT_ONLY_KEY="ctrl+shift+t"
sudo -E ./build/keygeist

# Or run without sudo after setting up udev rules (see Installation section)
./build/keygeist
```

### How it works

1. Run the operator with `make run`
2. Press one of the key combinations to activate the AI assistant:
   - `Windows + C` for clipboard context
   - `Windows + S` for screenshot context
   - `Windows + E` for both clipboard and screenshot context
   - `Windows + T` for text-only context (no additional context)
3. A zenity dialog will appear asking for your question
4. Enter your question and click OK
5. The AI response will be automatically typed into the currently focused application
6. Press the same combination again during an interaction to cancel it

## Debugging Tools

The project includes two debugging utilities for development and testing:

### 1. Keyboard Emulator (`cmd/emulator/main.go`)

Interactive keyboard emulator that allows you to simulate keyboard input for testing purposes.

```bash
# Build and run (requires sudo or udev rules setup)
make build-emulator
sudo ./build/emulator

# Usage examples:
press 30    # Press key 'A'
release 30  # Release key 'A'
tap 30      # Tap key 'A'
type "Hello World"  # Type text
hotkey 29 56 23  # Press Ctrl+Alt+F
```

### 2. Keyboard Listener (`cmd/listener/main.go`)

Listens for specific key combinations and executes callbacks for debugging key detection.

```bash
# Build and run (requires sudo or udev rules setup)
make build-listener
sudo ./build/listener

# Currently configured to detect `Windows + C` combination
```

## Installation

### Dependencies

```bash
# Install zenity (for GUI dialogs)
sudo dnf install zenity  # Fedora/RHEL
sudo apt install zenity  # Ubuntu/Debian

# Install Go dependencies
make deps
```

### Udev Rules Setup

To run Keygeist without sudo, you need to set up udev rules that allow your user to access uinput and input devices. This is more secure than running with sudo.

#### Create udev rules file

Create a new udev rules file:

```bash
sudo nano /etc/udev/rules.d/99-uinput-keyboard.rules
```

Add the following content (replace `YOUR_USERNAME` with your actual username):

```
# Allow user to access uinput device
KERNEL=="uinput", MODE="0660", GROUP="uinput", OPTIONS+="static_node=uinput"

# Allow user to read input devices
KERNEL=="input*", MODE="0644", GROUP="input"
```

#### Create required groups and add user

```bash
# Create uinput group if it doesn't exist
sudo groupadd -f uinput

# Create input group if it doesn't exist
sudo groupadd -f input

# Add your user to both groups
sudo usermod -a -G uinput,input $USER

# Load uinput module
echo uinput | sudo tee -a /etc/modules
sudo modprobe uinput
```

#### Reload udev rules

```bash
# Reload udev rules
sudo udevadm control --reload-rules
sudo udevadm trigger

# Verify the rules are working
ls -la /dev/uinput
ls -la /dev/input/
```

#### Log out and log back in

After adding yourself to the groups, you need to log out and log back in for the group changes to take effect.

#### Test without sudo

After setting up the udev rules, you should be able to run the commands without sudo:

```bash
# Test the main application
./build/keygeist

# Test debugging tools
./build/emulator
./build/listener
```

**Note**: If you still encounter permission issues, you may need to reboot your system for all changes to take effect.

### Systemd User Service

To automatically start the keygeist when you log in, you can create a systemd user service.

#### Create the service file

Create a new systemd user service file:

```bash
mkdir -p ~/.config/systemd/user
nano ~/.config/systemd/user/keygeist.service
```

Add the following content (adjust paths and environment variables as needed):

```ini
[Unit]
Description=Keygeist Service
After=graphical-session.target
Wants=graphical-session.target

[Service]
Type=simple
ExecStart=/path/to/your/keygeist/build/keygeist
Restart=always
RestartSec=5
Environment=OPENAI_API_KEY=your_api_key_here
Environment=OPENAI_MODEL=gpt-3.5-turbo
Environment=OPENAI_BASE_URL=https://api.openai.com/v1
Environment=OPENAI_SYSTEM_PROMPT=You are a coding assistant. Provide concise, practical code solutions.
# Optional: Custom keybindings
# Environment=CLIPBOARD_KEY=ctrl+shift+c
# Environment=SCREENSHOT_KEY=ctrl+shift+s
# Environment=ALL_CONTEXT_KEY=ctrl+shift+e
# Add any other environment variables you need

[Install]
WantedBy=default.target
```

**Important**: Replace `/path/to/your/keygeist/build/keygeist` with the actual path to your keygeist binary, and set your actual OpenAI API key.

#### Enable and start the service

```bash
# Reload systemd user daemon
systemctl --user daemon-reload

# Enable the service to start on login
systemctl --user enable keygeist.service

# Start the service immediately
systemctl --user start keygeist.service

# Check the service status
systemctl --user status keygeist.service

# View service logs
journalctl --user -u keygeist.service -f
```

#### Service management commands

```bash
# Stop the service
systemctl --user stop keygeist.service

# Restart the service
systemctl --user restart keygeist.service

# Disable auto-start
systemctl --user disable keygeist.service

# Check if service is enabled
systemctl --user is-enabled keygeist.service
```

#### Troubleshooting

If the service fails to start, check the logs:

```bash
# View recent logs
journalctl --user -u keygeist.service --since "1 hour ago"

# View all logs
journalctl --user -u keygeist.service
```

Common issues:
- **Permission denied**: Make sure you've set up the udev rules correctly
- **Environment variables**: Ensure all required environment variables are set in the service file
- **Path issues**: Verify the ExecStart path is correct and the binary exists

#### Alternative: Using environment file

Instead of hardcoding environment variables in the service file, you can use an environment file:

1. Create an environment file:
```bash
nano ~/.config/systemd/user/keygeist.env
```

2. Add your environment variables:
```bash
OPENAI_API_KEY=your_api_key_here
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_SYSTEM_PROMPT=You are a coding assistant. Provide concise, practical code solutions.
```

3. Update the service file to use the environment file:
```ini
[Service]
Type=simple
ExecStart=/path/to/your/keygeist/build/keygeist
Restart=always
RestartSec=5
EnvironmentFile=%h/.config/systemd/user/keygeist.env
```

### Build

```bash
# Build the main application
make build

# Build for specific platform
make build-linux
make build-all

# Build debugging tools
make build-tools
make build-emulator
make build-listener
```

### Install system-wide

```bash
make install
```

## Development

```bash
# Format code
make fmt

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Clean build artifacts
make clean
```


## License

See [LICENSE](LICENSE) file for details.
