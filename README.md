# Radar

Radar is an application that allows to map networks.
Created by :

- Alexandre Duchesne
- Noé Steiner

## Installation

**To install Radar, you will have to install Go on your computer.**

Clone the repository and run the following command:

```bash
cd /path/to/radar
go get main
```

## Usage

Radar can be used directly on your computer, or you can also use it in remote mode.

### Use it on your computer

You only have to run the following command:

```bash
cd /path/to/radar
go run main.go
```

Then you will see multiple commands, you can directly use Scan-ports and Scan-gateway commands.
The port scan is done on the provided IP.

You can also use the commands directly if you prefer.

### Remote mode

To use remote mode, you will have to put Radar on a server and run the following command:

```bash
cd /path/to/radar
go run main.go
```

Then you will have to select "server-mode"
It will launch Radar daemon on 6666 TCP port.
Every scan will then be done directly from the server.

Then, you can use Radar in client mode on your computer:

```bash
cd /path/to/radar
go run main.go
```

Then you will have to select "client-mode" and enter the server IP address.
You can then use Scan-ports and Scan-gateway commands.

### Note

To run the application, you might need to give it the right to send ICMP packets.
For development, we set `net.ipv4.ping_group_range="0 2147483647"`

## Understand the code

The code is divided in multiple packages:

- cmd
- internal

### cmd

cmd contains every callable command in the application. It is the entry point of the application. Most functions can be
called without the TUI even though the TUI might be more optimized.

You will also find in this package the root view containing the main menu.

### internal

internal contains every function that are used in the application. It is the core of the application.

It mostly contains helpers to scan ports and gateways, and the functions to achieve these scans.