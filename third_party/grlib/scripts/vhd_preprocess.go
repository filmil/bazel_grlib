package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	prefix := ""
	inputPath := ""
	configHPath := ""
	outputPath := ""

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--prefix" && i+1 < len(os.Args) {
			prefix = strings.ToUpper(os.Args[i+1])
			i++
		} else if inputPath == "" {
			inputPath = os.Args[i]
		} else if configHPath == "" {
			configHPath = os.Args[i]
		} else if outputPath == "" {
			outputPath = os.Args[i]
		}
	}

	if inputPath == "" || configHPath == "" || outputPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [--prefix PREFIX] <input_in_vhd> <config_h> <output_vhd>\n", os.Args[0])
		os.Exit(1)
	}

	// 1. Read config.h
	defs := make(map[string]string)
	cf, err := os.Open(configHPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config_h: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(cf)
	reDef := regexp.MustCompile(`^#define\s+(CONFIG_[a-zA-Z0-9_]+)\s+(.*)$`)
	for scanner.Scan() {
		match := reDef.FindStringSubmatch(scanner.Text())
		if len(match) > 2 {
			defs[match[1]] = strings.Trim(strings.TrimSpace(match[2]), "\"")
		}
	}
	cf.Close()

	// 2. Process input
	inf, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer inf.Close()

	outf, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer outf.Close()

	// Regex for macros like CONFIG_IU_NWINDOWS or 16#CONFIG_AHB_IOADDR#
	// We look for CONFIG_ followed by alphanumeric/underscore
	reMacro := regexp.MustCompile(`CONFIG_[a-zA-Z0-9_]+`)

	scanner = bufio.NewScanner(inf)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip C-style preprocessor lines like #include (GRLIB uses them, but we handle hierarchy differently)
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		processedLine := reMacro.ReplaceAllStringFunc(line, func(macro string) string {
			// 1. Try namespaced: PREFIX + MACRO (without CONFIG_ prefix of macro)
			if prefix != "" {
				suffix := strings.TrimPrefix(macro, "CONFIG_")
				namespaced := "CONFIG_" + prefix + "_" + suffix
				if val, ok := defs[namespaced]; ok {
					return formatValue(val)
				}
			}
			
			// 2. Try original name
			if val, ok := defs[macro]; ok {
				return formatValue(val)
			}
			
			// 3. Fallback to 0 if not found (GRLIB default behavior for undefined configs)
			return "0"
		})
		
		fmt.Fprintln(outf, processedLine)
	}
}

func formatValue(val string) string {
	if val == "y" { return "1" }
	if val == "n" { return "0" }
	// If it's a hex value without prefix, just return it (VHDL might need 16# prefix which is usually in the template)
	return val
}
