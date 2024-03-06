# Main.go

This Go program implements a simple command-line tool for generating random passwords.

## Description

The `main` function is the entry point of the program. It performs the following steps:

1. Parses command-line flags to determine the desired length of the password and whether to include symbols in the generated password.
2. Defines a character set containing lowercase letters, uppercase letters, digits, and optionally symbols based on the user's preference.
3. Initializes a `strings.Builder` to efficiently build the password string.
4. Seeds the random number generator with the current time to ensure different random sequences across program runs.
5. Generates a random password by repeatedly selecting random characters from the character set and appending them to the `strings.Builder` until the desired password length is reached.
6. Prints the generated password to the console.

## Usage

To use the program, run the `main.go` file with the following optional command-line flags:

- `-length`: Specifies the desired length of the generated password. Default is 16.
- `-symbols`: Specifies whether to include symbols in the generated password. Default is false.

Example usage:
```
go run main.go -length=20 -symbols
```

This will generate a random password with a length of 20 characters, including symbols.

## Dependencies

The program uses the following packages from the Go standard library:
- `flag`: For parsing command-line flags.
- `fmt`: For printing the generated password.
- `math/rand`: For generating random numbers.
- `strings`: For efficiently building the password string.
- `time`: For seeding the random number generator.

## License

This program is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).