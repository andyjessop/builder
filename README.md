# Go Project README

## Overview

This Go project consists of a single `main.go` file that performs the following tasks:

1. Loads environment variables from a `.env` file.
2. Parses command-line flags to get the file path to overwrite.
3. Prompts the user to enter a prompt and reads the input.
4. Converts the content of the specified file to a string.
5. Sends the file content and user prompt to an API to get a code response.
6. Generates a README.md file based on the code response.
7. Prompts the user to enter a new branch name.
8. Initializes a Git repository in the current directory.
9. Creates a new branch with the specified name.
10. Writes the code response to the specified file.
11. Writes the generated README.md content to a file.
12. Adds and commits the changes to the Git repository on the new branch.

## Dependencies

The project uses the following external dependencies:

- `github.com/joho/godotenv`: Used for loading environment variables from a `.env` file.

## Usage

To run the project, follow these steps:

1. Set the `API_KEY` environment variable in a `.env` file.
2. Run the `main.go` file with the following command-line flag:
   - `--file`: Specifies the path of the file to overwrite (default: `./main.go`).
3. Enter a prompt when prompted.
4. Enter a new branch name when prompted.

The project will generate a new code file based on the provided prompt, create a README.md file, and commit the changes to a new branch in the Git repository.

## Functions

The `main.go` file contains the following functions:

- `main()`: The entry point of the program. It orchestrates the entire process of loading environment variables, parsing flags, getting user input, generating code and README, initializing a Git repository, creating a new branch, writing files, and committing changes.
- `getCodeResponse(targetFileContent, prompt, apiKey string) (string, error)`: Sends the target file content and user prompt to an API to get a code response. It retries the API request multiple times if needed.
- `getReadmeResponse(codeResponse string, apiKey string) (string, error)`: Generates a README.md file based on the code response. It retries the API request multiple times if needed.
- `ask(message string, apiKey string) (string, error)`: Sends a request to the Anthropic API with the provided message and API key. It returns the API response as a string.
- `convertFileToString(filePath string) (string, error)`: Reads the content of a file and returns it as a string.
- `writeStringToFile(content string, filename string) error`: Writes the provided content to a file with the specified filename.
- `initGitRepo() error`: Initializes a Git repository in the current directory if it doesn't exist.
- `createBranch(branchName string) error`: Creates a new branch with the specified name in the Git repository.
- `addAndCommitChanges(branchName string) error`: Adds all changes to the Git repository and commits them with a message including the branch name.

## Error Handling

The project includes error handling for various operations, such as loading environment variables, parsing flags, reading files, making API requests, and executing Git commands. Errors are logged to the console, and the program terminates gracefully if an error occurs.

## Acknowledgements

The project utilizes the Anthropic API for generating code responses and README content. It also relies on the `godotenv` package for loading environment variables from a `.env` file.