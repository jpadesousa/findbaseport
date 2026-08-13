# findbaseport

`findbaseport` is a CLI tool that searches a specified port range and returns the starting port of the first contiguous block of available ports with the requested size.

## Motivation

This tool was created to simplify finding a suitable base port for a [CryoSPARC](https://cryosparc.com/) installation.

CryoSPARC requires a base port and expects a range of **10 contiguous ports** to be available starting from that port. Finding an appropriate base port manually can be inconvenient, particularly on systems where some ports are already in use.

`findbaseport` automates this process by searching a specified port range and returning the first port from which the requested number of contiguous ports are available.

For example, when looking for a base port with the default requirement of 10 contiguous ports, the tool searches for a sequence such as:

```text
8000  8001  8002  8003  8004  8005  8006  8007  8008  8009
```

If all 10 ports are available, `8000` can be used as the base port. If any of those ports are unavailable, the tool continues searching until it finds the next suitable contiguous block.

## Description

The tool searches an **inclusive** port range for available ports.

Given a range and a requested block size, `findbaseport` finds the first sequence of contiguous available ports that is large enough and returns the starting port of that block.

For example, if the requested block size is `3` in a range of `8000` to `8020`:
```bash
findbaseport --start 8000 --end 8020 --count 3
```

and ports `8000`, `8001`, and `8002` are available, the command returns:

```text
8000
```

If those ports are unavailable but `8010`, `8011`, and `8012` are available, it returns:

```text
8010
```

## Installation

### Build from source

Clone the repository and use the included Makefile:

```bash
git clone <repository-url>
cd <repository-directory>

make build
```

The resulting binary is placed in:

```text
bin/findbaseport
```

You can then run it with:

```bash
./bin/findbaseport --help
```

## Usage

```bash
./bin/findbaseport --help
```

The help output describes the available command-line options and arguments.

The command searches within the specified port range and returns the start port of the first contiguous block of available ports with the requested size.

## Development

The project provides a Makefile with commands for formatting, static analysis, linting, and building.

### Run everything

```bash
make
```

The default target runs:

1. `fmt`
2. `vet`
3. `lint`
4. `build`

### Format

Format all Go source files:

```bash
make fmt
```

### Vet

Run Go's static analysis:

```bash
make vet
```

### Lint

Run `golangci-lint`:

```bash
make lint
```

Make sure `golangci-lint` is installed and available in your `PATH`.

### Build

Build the CLI:

```bash
make build
```

The binary is written to:

```text
bin/findbaseport
```

### Clean

Remove build artifacts:

```bash
make clean
```

## Requirements

* Go
* `make`
* `golangci-lint` for linting

## Make Targets

| Target       | Description                        |
| ------------ | ---------------------------------- |
| `make`       | Format, vet, lint, and build |
| `make fmt`   | Format Go source code              |
| `make vet`   | Run `go vet`                       |
| `make lint`  | Run `golangci-lint`                |
| `make build` | Build `findbaseport` into `bin/`   |
| `make clean` | Remove the `bin/` directory        |

## License

See the repository's `LICENSE` file for licensing information.

