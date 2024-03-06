package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Use the loaded environment variables
	apiKey := os.Getenv("API_KEY")

	cwf, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cwd := filepath.Dir(cwf)

	// Define the command line flag for the file path
	filePath := flag.String("file", "./main.go", "Path of the file to overwrite")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
		fmt.Println("Error: File path not provided. Please use the --file flag to specify the file path.")
		flag.Usage()
		return
	}

	// Get the directory of the file path
	fileDir := filepath.Dir(*filePath)

	// Change to the file directory
	err = os.Chdir(fileDir)
	if err != nil {
		fmt.Printf("Error changing to directory %s: %v\n", fileDir, err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your prompt: ")
	scanner.Scan()
	prompt := scanner.Text()

	targetFileContent, err := convertFileToString(*filePath)
	if err != nil {
		fmt.Printf("Error converting file %s to string: %v\n", *filePath, err)
		return
	}

	p := `You are an expert Golang programmer with many years of experience - nothing is beyond you. But you know only code, not spoken languages. When you return your response, you must only return code. Golang code. Nothing else. It is absolutely crucial that you adhere to this rule. What you return will be written directly to disk in a mission-critical file. The code you write will be a Go file, so it should include the "package" declaration, any imports, and the necessary functions. Your code should be written in a human-readable manner, and all error messages should be very explicit and lead the reader by the hand to the right fix.

Below, you will be given the current code for the target file. It is your job to return new code adhering to the prompt. You should return the full file code. Do not import anything outside of the standard library. Only use the standard library.

Please ensure your code is accurate and error-free. Double-check everything before returning your response. Aim for zero mistakes.

You must output the full code, from start to finish. Do not leave anything out.

Your response is limited to 4028 tokens, so you must absolutely be sure that your response is under that. I suggest giving yourself a hard limit of 3000 tokens, just to be sure.

Remember that you should only return code without explanation. Your output must start with 'package' and end with '}'.

The code will be given after '=CODE=' and the prompt will be given after '=PROMPT='.`

	p += "\n\n=CODE=\n\n" + targetFileContent
	p += "\n\n=PROMPT=\n\n" + prompt

	var apiResponse struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		response, err := ask(p, apiKey)
		if err != nil {
			fmt.Printf("Error (attempt %d/%d): %v\n", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		err = json.Unmarshal([]byte(response), &apiResponse)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON response (attempt %d/%d): %v\n", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		if len(apiResponse.Content) == 0 {
			fmt.Printf("\033[33mEmpty LLM response content. Retrying...\033[0m\n")
			time.Sleep(retryDelay)
			continue
		}

		break
	}

	if len(apiResponse.Content) == 0 {
		fmt.Println("Failed to get a valid response after", maxRetries, "attempts")
		return
	}

	responseText := apiResponse.Content[0].Text
	unescapedText := html.UnescapeString(responseText)

	// Trim any leading or trailing whitespace, newlines, or backticks
	trimmedText := strings.TrimSpace(unescapedText)
	trimmedText = strings.Trim(trimmedText, "\n")
	trimmedText = strings.Trim(trimmedText, "`")

	// Prompt for the branch name
	fmt.Printf("\033[32mSuccessfully generated new code.\033[0m\n")
	fmt.Print("Enter the new branch name: ")
	scanner.Scan()
	branchName := scanner.Text()

	// Initialize a new Git repository if it doesn't exist
	err = initGitRepo()
	if err != nil {
		fmt.Printf("Error initializing Git repository: %v\n", err)
		return
	}

	// Create a new branch
	err = createBranch(branchName)
	if err != nil {
		fmt.Printf("Error creating branch: %v\n", err)
		return
	}

	// Write the response to the specified file
	err = writeStringToFile(trimmedText, filepath.Base(*filePath))
	if err != nil {
		fmt.Printf("Error writing response to file %s: %v\n", *filePath, err)
		return
	}

	// Add and commit the changes
	err = addAndCommitChanges(branchName)
	if err != nil {
		fmt.Printf("Error adding and committing changes: %v\n", err)
		return
	}

	// Change back to the original working directory
	err = os.Chdir(cwd)
	if err != nil {
		fmt.Printf("Error changing back to the original directory: %v\n", err)
		return
	}

	// Print success message in green color
	fmt.Printf("\033[32mSuccessfully created branch %s, wrote response to %s, and committed the changes\033[0m\n", branchName, *filePath)
}

func ask(message string, apiKey string) (string, error) {
	// URL and data for the POST request
	url := "https://api.anthropic.com/v1/messages"
	data := map[string]interface{}{
		"model":      "claude-3-opus-20240229",
		"max_tokens": 4096,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": message,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return "", err
	}

	return string(body), nil
}

func convertFileToString(filePath string) (string, error) {
	// Read the contents of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert the content to a string
	return string(content), nil
}

func writeStringToFile(content string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func initGitRepo() error {
	// Check if the current directory is already a Git repository
	_, err := os.Stat(".git")
	if err == nil {
		// Git repository already exists, no need to initialize
		return nil
	}

	// Initialize a new Git repository
	cmd := exec.Command("git", "init")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error initializing Git repository: %w", err)
	}

	return nil
}

func createBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating branch: %w", err)
	}
	return nil
}

func addAndCommitChanges(branchName string) error {
	// Add changes
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding changes: %w", err)
	}

	// Commit changes
	commitMsg := fmt.Sprintf("Update file on branch %s", branchName)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %w", err)
	}

	return nil
}
