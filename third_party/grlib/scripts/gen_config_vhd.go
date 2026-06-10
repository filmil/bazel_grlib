package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_h> <output_vhd>\n", os.Args[0])
		os.Exit(1)
	}

	inputH := os.Args[1]
	outputVhd := os.Args[2]

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
	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 2 {
			name := match[1]
			val := strings.TrimSpace(match[2])
			val = strings.Trim(val, "\"") // Strip quotes first

			// GRLIB convention: constants in config.vhd are often named CFG_*
			cfgName := strings.Replace(name, "CONFIG_", "CFG_", 1)

			if val == "1" || val == "0" {
				fmt.Fprintf(df, "  constant %s : integer := %s;\n", cfgName, val)
			} else if strings.HasPrefix(val, "0x") {
				fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", cfgName, strings.TrimPrefix(val, "0x"))
			} else if reHex.MatchString(val) && len(val) > 1 {
				// Likely a hex value without 0x prefix from Kconfig
				fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", cfgName, val)
			} else {
				fmt.Fprintf(df, "  constant %s : integer := %s;\n", cfgName, val)
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
