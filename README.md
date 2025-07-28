# Viper + Validator Example

A Go program demonstrating configuration management with [Viper](https://github.com/spf13/viper) and input validation with [go-playground/validator](https://github.com/go-playground/validator).

## Features

- Command-line flag parsing with [pflag](https://github.com/spf13/pflag)
- Configuration file support (YAML)
- Environment variable binding
- Input validation with detailed error messages

## Usage

```bash
# Build the program
go build

# Show help
./viper-validator-example --help

# Run with required color flag
./viper-validator-example --color red --size 25 --count 5

# Run with minimal required flags (size and count have defaults)
./viper-validator-example --color blue
```

## Configuration Options

| Flag                    | Type   | Required | Range/Values             | Default |
| ----------------------- | ------ | -------- | ------------------------ | ------- |
| `--color`               | string | Yes      | red, green, blue, yellow | -       |
| `--size`                | int    | Yes      | 1-100                    | 10      |
| `--count`               | int    | No       | 1-1000                   | 1       |
| `--include-cup-holders` | bool   | No       | true, false              | false   |

## Configuration Sources (in priority order)

1. Command-line flags
2. Environment variables (`EXAMPLE_COLOR`, `EXAMPLE_SIZE`, `EXAMPLE_COUNT`, `EXAMPLE_INCLUDE_CUPHOLDERS`)
3. Configuration file (`config.yaml`)
4. Default values

## Example Output

```bash
# Valid input
$ ./viper-validator-example --color green --size 50
No config file found, using defaults and command line flags
Configuration loaded successfully:
Color: green
Size: 50
Count: 1
Include cup holders: false

# Valid input
$ EXAMPLE_INCLUDE_CUPHOLDERS=true ./viper-validator-example --color green --size 50
No config file found, using defaults and command line flags
Configuration loaded successfully:
Color: green
Size: 50
Count: 1
Include cup holders: true

# Invalid input
$ ./viper-validator-example --color purple
Validation failed:
- Color must be one of: red, green, blue, yellow
```
