# locaster

Locaster is a simple and efficient solution for casting your PC screen to cheap TVs that only support displaying JPEG images through a browser. This project allows you to stream your screen to any device with a web browser, making it ideal for older TVs or devices with limited capabilities.

## Features

- **Screen Capture**: Capture your PC screen and convert it to JPEG images.
- **Real-time Streaming**: Stream the captured images to a web browser in real-time.
- **Simple Setup**: Easy to set up and use with minimal configuration.
- **Cross-Platform**: Works on Linux, Windows, and macOS.

## How It Works

1. **Screen Capture**: The application captures the screen of your PC at regular intervals.
2. **Image Conversion**: The captured screen is converted to a JPEG image.
3. **Streaming**: The JPEG image is sent to a web server.
4. **Display**: The web server serves the JPEG image to any connected web browser, updating it in real-time.

## Getting Started

### Prerequisites

- Go 1.23.3 or later
- A device with a web browser

### Installation

Download the binary file for your operating system from the Releases section.

### Usage

1. Open the casting URL on your PC:
    ```
    http://<your-pc-ip>:8080/cast
    ```

2. Open the player URL on your TV or any device with a web browser:
    ```
    http://<your-pc-ip>:8080
    ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.