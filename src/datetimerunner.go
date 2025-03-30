package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println()
	fmt.Println("datetimerunner v1.1.0 - https://github.com/gmagana/datetimerunner")
	fmt.Println("\u00A9 2025 - Gabriel Magana-Gonzalez - g@g11a.com")
	fmt.Println()

	// Define the test flag
	testFlag := flag.Bool("t", false, "Test mode - don't execute, just show what would be executed")

	// Parse flags but keep the non-flag arguments
	flag.Parse()
	args := flag.Args()

	// Check if we have enough arguments
	if len(args) < 2 {
		fmt.Println("Error: Missing required arguments.")
		fmt.Println("Usage: datetimerunner <executable_name> <argument_template> [-t]")
		os.Exit(1)
	}

	// Get the executable name
	executableName := args[0]

	// Join all remaining arguments as the template, preserving quotes
	// This ensures that quoted arguments with spaces are treated as a single argument
	argumentTemplate := strings.Join(args[1:], " ")

	// Get the current time
	now := time.Now()

	// Process the argument template by replacing placeholders
	processedArgument := replacePlaceholders(argumentTemplate, now)

	// Quote the executable path if it contains spaces
	quotedExecutableName := quotePathIfNeeded(executableName)

	// Display information in test mode or execute the command
	if *testFlag {
		fmt.Println("Test Mode: Command will not be executed")
		fmt.Println()
		fmt.Printf("Executable: %s\n", quotedExecutableName)
		fmt.Println()
		fmt.Printf("Original Argument Template: %s\n", argumentTemplate)
		fmt.Println()
		fmt.Printf("Processed Argument: %s\n", processedArgument)
		fmt.Println()
		fmt.Printf("Command that would be executed: %s %s\n", quotedExecutableName, processedArgument)
		fmt.Println()
	} else {
		fmt.Printf("Executing command: %s %s\n", quotedExecutableName, processedArgument)

		// Execute the command
		// For Windows, we need a different approach to handle nested quotes
		if isWindows() {
			// Basic approach: directly execute the command using Go's exec.Command
			cmd := exec.Command(executableName)

			// Split the arguments while preserving quotes
			cmdArgs := parseCommandLine(processedArgument)

			// Set the arguments
			cmd.Args = append([]string{executableName}, cmdArgs...)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error executing command:", err)
				os.Exit(1)
			}
		} else {
			// On other platforms, use the command directly
			cmdArgs := parseCommandLine(processedArgument)
			cmd := exec.Command(executableName, cmdArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error executing command:", err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("Done.")
}

// quotePathIfNeeded adds quotes around a path if it contains spaces and isn't already quoted
func quotePathIfNeeded(path string) string {
	// If the path already starts and ends with quotes, leave it as is
	if strings.HasPrefix(path, "\"") && strings.HasSuffix(path, "\"") {
		return path
	}

	// If the path contains spaces, enclose it in quotes
	if strings.Contains(path, " ") {
		return "\"" + path + "\""
	}

	// Otherwise, return the path as is
	return path
}

// isWindows returns true if the current OS is Windows
func isWindows() bool {
	return runtime.GOOS == "windows"
}

// parseCommandLine parses a command line string into individual arguments
// preserving quoted sections
func parseCommandLine(cmdLine string) []string {
	fmt.Printf("Parsing command line: |%s|\n", cmdLine)

	var args []string
	var currentArg strings.Builder
	inQuotes := false
	escapeNext := false

	for i := 0; i < len(cmdLine); i++ {
		char := cmdLine[i]

		// Handle escape character (backslash)
		if escapeNext {
			currentArg.WriteByte(char)
			escapeNext = false
			continue
		}

		if char == '\\' && i+1 < len(cmdLine) {
			if cmdLine[i+1] == '"' {
				// Escape the next character (likely a quote)
				escapeNext = true
				continue
			}
		}

		// This method of counting quotes may not work in all platfrms, Go on Windows does pre-processing
		// of the command line, so we cannot control accurate processing of escaped quotes. Advise the
		// user to used escpaed quotes (\") whenever possible to ensure correct processing.
		switch char {
		case '"':
			// Check for double quotes (escaped quotes)
			if i+1 < len(cmdLine) && cmdLine[i+1] == '"' {
				// Handle double quotes (""this is quoted"")
				currentArg.WriteByte('"')
				i++ // Skip the next quote
			} else {
				// Toggle quote state
				inQuotes = !inQuotes
			}
		case ' ':
			if inQuotes {
				// Space inside quotes, add to current argument
				currentArg.WriteByte(' ')
			} else {
				// Space outside quotes, finish current argument if not empty
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			}
		default:
			// Add any other character to the current argument
			currentArg.WriteByte(char)
		}
	}

	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	// Process the arguments to properly handle nested quotes
	for i, arg := range args {
		// If the argument has double quotes, convert them to proper escaping
		if strings.Count(arg, "\"") >= 2 {
			// This is a proper quoted argument with actual quotes
			// Ensure it's properly quoted for the command line
			args[i] = ensureProperQuoting(arg)
		}
	}

	return args
}

// ensureProperQuoting ensures that an argument is properly quoted
// for command line execution
func ensureProperQuoting(arg string) string {
	// If it starts and ends with quotes, it's already fully quoted
	if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
		return arg
	}

	// For arguments that contain quotes but aren't fully quoted,
	// ensure they're properly quoted
	if strings.Contains(arg, "\"") {
		// Replace pairs of quotes with a single quote
		// This handles the case of ""quoted text"" which needs to become "quoted text"
		arg = strings.ReplaceAll(arg, "\"\"", "\"")
	}

	// If the argument contains spaces, add quotes around it
	if strings.Contains(arg, " ") && !strings.HasPrefix(arg, "\"") {
		return "\"" + arg + "\""
	}

	return arg
}

// replacePlaceholders replaces all date/time placeholders in the input string with current values
func replacePlaceholders(input string, t time.Time) string {
	// Define all the replacements based on the current time
	replacements := map[string]string{
		"<y>":       strconv.Itoa(t.Year() % 100),
		"<yy>":      fmt.Sprintf("%02d", t.Year()%100),
		"<yyyy>":    strconv.Itoa(t.Year()),
		"<m>":       strconv.Itoa(int(t.Month())),
		"<mm>":      fmt.Sprintf("%02d", t.Month()),
		"<mmm>":     t.Month().String()[:3],
		"<mmmm>":    t.Month().String(),
		"<d>":       strconv.Itoa(t.Day()),
		"<dd>":      fmt.Sprintf("%02d", t.Day()),
		"<ddd>":     fmt.Sprintf("%03d", t.YearDay()),
		"<h24>":     strconv.Itoa(t.Hour()),
		"<hh24>":    fmt.Sprintf("%02d", t.Hour()),
		"<h12>":     strconv.Itoa(t.Hour()%12 + 1),
		"<hh12>":    fmt.Sprintf("%02d", t.Hour()%12+1),
		"<ampm>":    getAMPM(t),
		"<i>":       strconv.Itoa(t.Minute()),
		"<mi>":      fmt.Sprintf("%02d", t.Minute()),
		"<s>":       strconv.Itoa(t.Second()),
		"<ss>":      fmt.Sprintf("%02d", t.Second()),
		"<dow>":     t.Weekday().String()[:3],
		"<weekday>": t.Weekday().String(),
	}

	// Replace each placeholder with its corresponding value
	result := input
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// getAMPM returns "AM" or "PM" based on the hour
func getAMPM(t time.Time) string {
	if t.Hour() < 12 {
		return "AM"
	}
	return "PM"
}
