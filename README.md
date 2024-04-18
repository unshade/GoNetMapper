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

### Remote mode

To use remote mode, you will have to put Radar on a server and run the following command:

```bash
cd /path/to/radar
go run main.go
```

Then you will have to select "server-mode"
It will launch Radar daemon on 6666 TCP port.

Then, you can use Radar in client mode on your computer:
    
```bash
cd /path/to/radar
go run main.go
```

Then you will have to select "client-mode" and enter the server IP address.
You can then use Scan-ports and Scan-gateway commands.

## Understand the code

The code is divided in multiple packages:
- cmd
- internal

### cmd

cmd contains every callable command in the application. It is the entry point of the application.

### internal

internal contains every function that are used in the application. It is the core of the application.
