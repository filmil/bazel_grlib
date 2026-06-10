package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	prefix := ""
	inputH := ""
	outputVhd := ""

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--prefix" && i+1 < len(os.Args) {
			prefix = strings.ToUpper(os.Args[i+1])
			i++
		} else if inputH == "" {
			inputH = os.Args[i]
		} else if outputVhd == "" {
			outputVhd = os.Args[i]
		}
	}

	if inputH == "" || outputVhd == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [--prefix PREFIX] <input_h> <output_vhd>\n", os.Args[0])
		os.Exit(1)
	}

	sf, err := os.Open(inputH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer sf.Close()

	df, err := os.Create(outputVhd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer df.Close()

	fmt.Fprintln(df, "library ieee;")
	fmt.Fprintln(df, "use ieee.std_logic_1164.all;")
	fmt.Fprintln(df, "library grlib;")
	fmt.Fprintln(df, "use grlib.config_types.all;")
	fmt.Fprintln(df, "")
	fmt.Fprintln(df, "package config is")

	re := regexp.MustCompile(`^#define\s+(CONFIG_[a-zA-Z0-9_]+)\s+(.*)$`)
	reHex := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	
	defs := make(map[string]string)
	
	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 2 {
			name := match[1]
			val := strings.TrimSpace(match[2])
			// If it's a string literal, we keep the quotes for later type detection
			defs[name] = val
		}
	}

	// 1. Output all namespaced symbols
	fmt.Fprintln(df, "  -- Namespaced symbols")
	for name, val := range defs {
		cfgName := strings.Replace(name, "CONFIG_", "CFG_", 1)
		outputVHDLConstant(df, cfgName, val, reHex)
	}

	// 2. Promote selected design symbols to generic names
	if prefix != "" {
		fmt.Fprintf(df, "\n  -- Promoted symbols for prefix: %s\n", prefix)
		fullPrefix := "CONFIG_" + prefix + "_"
		for name, val := range defs {
			if strings.HasPrefix(name, fullPrefix) {
				genericName := "CFG_" + strings.TrimPrefix(name, fullPrefix)
				// Strip numeric uniqueness suffix
				genericName = regexp.MustCompile(`_[0-9]+$`).ReplaceAllString(genericName, "")
				outputVHDLConstant(df, genericName, val, reHex)
			}
		}
	}

	// Add the mandatory GRLIB_CONFIG_ARRAY
	fmt.Fprintln(df, "")
	fmt.Fprintln(df, "  constant GRLIB_CONFIG_ARRAY : grlib_config_array_type := (")
	fmt.Fprintln(df, "    others => 0);")
	fmt.Fprintln(df, "")
	fmt.Fprintln(df, "end;")
}

func outputVHDLConstant(df io.Writer, name, val string, reHex *regexp.Regexp) {
	// If the name is CFG_ACTIVE_DESIGN_PREFIX, don't output as integer
	if name == "CFG_ACTIVE_DESIGN_PREFIX" {
		fmt.Fprintf(df, "  constant %s : string := %s;\n", name, val)
		return
	}

	trimmedVal := strings.Trim(val, "\"")

	if trimmedVal == "1" || trimmedVal == "0" {
		fmt.Fprintf(df, "  constant %s : integer := %s;\n", name, trimmedVal)
	} else if strings.HasPrefix(trimmedVal, "0x") {
		fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", name, strings.TrimPrefix(trimmedVal, "0x"))
	} else if reHex.MatchString(trimmedVal) && len(trimmedVal) > 1 {
		fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", name, trimmedVal)
	} else if trimmedVal == "y" {
		fmt.Fprintf(df, "  constant %s : integer := 1;\n", name)
	} else if trimmedVal == "n" {
		fmt.Fprintf(df, "  constant %s : integer := 0;\n", name)
	} else if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		// String literal
		fmt.Fprintf(df, "  constant %s : string := %s;\n", name, val)
	} else {
		// Try parsing as integer
		fmt.Fprintf(df, "  constant %s : integer := %s;\n", name, trimmedVal)
	}
}
