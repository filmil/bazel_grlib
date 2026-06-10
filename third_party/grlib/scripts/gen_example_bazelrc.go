package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceDir != "" {
		if err := os.Chdir(workspaceDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error changing to workspace directory: %v\n", err)
			os.Exit(1)
		}
	}

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_list> <output_bazelrc>\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	sf, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer sf.Close()

	df, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer df.Close()

	fmt.Fprintln(df, "# Example bazelrc showing all available Kconfig parameters for GRLIB.")
	fmt.Fprintln(df, "# To use, copy relevant lines to your user.bazelrc or .bazelrc.")
	fmt.Fprintln(df, "")

	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		kind := parts[0]
		label := parts[2]
		
		// label is @grlib_config//:CONFIG_NAME
		name := strings.TrimPrefix(label, "@grlib_config//:")

		val := ""
		switch kind {
		case "bool_flag":
			val = "True"
		case "int_flag":
			val = "0"
		case "string_flag":
			val = "\"\""
		default:
			val = "0"
		}

		fmt.Fprintf(df, "build --@grlib_config//:%s=%s\n", name, val)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning input: %v\n", err)
	}
}
