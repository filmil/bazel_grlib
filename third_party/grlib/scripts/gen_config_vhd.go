package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
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

	reDef := regexp.MustCompile(`^#define\s+(CONFIG_[a-zA-Z0-9_]+)\s+(.*)$`)
	
	defs := make(map[string]string)
	
	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		match := reDef.FindStringSubmatch(line)
		if len(match) > 2 {
			name := match[1]
			val := strings.TrimSpace(match[2])
			defs[name] = val
		}
	}

	// 1. Output root symbols
	fmt.Fprintln(df, "  -- Root symbols")
	var rootNames []string
	for name := range defs {
		if !strings.HasPrefix(name, "CONFIG_LIB_") && !strings.HasPrefix(name, "CONFIG_DESIGNS_") && !strings.HasPrefix(name, "CONFIG_BIN_") {
			rootNames = append(rootNames, name)
		}
	}
	sort.Strings(rootNames)
	for _, name := range rootNames {
		cfgName := strings.Replace(name, "CONFIG_", "CFG_", 1)
		outputVHDLConstant(df, cfgName, defs[name])
	}

	// 2. Promote selected design symbols
	if prefix != "" {
		fmt.Fprintf(df, "\n  -- Promoted symbols for prefix: %s\n", prefix)
		fullPrefix := "CONFIG_" + prefix + "_"
		promoted := make(map[string]bool)
		
		var allNames []string
		for name := range defs { allNames = append(allNames, name) }
		sort.Strings(allNames)

		for _, name := range allNames {
			if strings.HasPrefix(name, fullPrefix) {
				genericName := "CFG_" + strings.TrimPrefix(name, fullPrefix)
				genericName = regexp.MustCompile(`_[0-9]+$`).ReplaceAllString(genericName, "")
				
				if !promoted[genericName] {
					fmt.Fprintf(df, "  -- Promoted from %s\n", name)
					outputVHDLConstant(df, genericName, defs[name])
					promoted[genericName] = true
				}
			}
		}
	}

	// Add the mandatory GRLIB_CONFIG_ARRAY
	fmt.Fprintln(df, "")
	fmt.Fprintln(df, "  constant GRLIB_CONFIG_ARRAY : grlib_config_array_type := (")
	
	getVal := func(name string, def int) int {
		if val, ok := defs["CONFIG_"+name]; ok {
			if val == "y" || val == "1" { return 1 }
			return 0
		}
		return def
	}

	fmt.Fprintf(df, "    4 => %d, -- grlib_sync_reset_enable_all\n", getVal("SYNC_RESET_ENABLE_ALL", 0))
	fmt.Fprintf(df, "    5 => %d, -- grlib_async_reset_enable\n", getVal("ASYNC_RESET_ENABLE", 0))
	fmt.Fprintf(df, "    9 => %d, -- grlib_amba_inc_nirq\n", getVal("AMBA_INC_NIRQ", 0))
	fmt.Fprintf(df, "    10 => %d, -- grlib_little_endian\n", getVal("LITTLE_ENDIAN", 0))
	fmt.Fprintln(df, "    others => 0);")
	fmt.Fprintln(df, "")
	fmt.Fprintln(df, "end;")
}

var reHexOnly = regexp.MustCompile(`^[0-9a-fA-F]+$`)
var reHasAlpha = regexp.MustCompile(`[a-fA-F]`)

func outputVHDLConstant(df io.Writer, name, val string) {
	if name == "CFG_ACTIVE_DESIGN_PREFIX" {
		fmt.Fprintf(df, "  constant %s : string := %s;\n", name, val)
		return
	}

	trimmedVal := strings.Trim(val, "\"")

	if trimmedVal == "1" || trimmedVal == "0" {
		fmt.Fprintf(df, "  constant %s : integer := %s;\n", name, trimmedVal)
	} else if strings.HasPrefix(trimmedVal, "0x") {
		fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", name, strings.TrimPrefix(trimmedVal, "0x"))
	} else if reHexOnly.MatchString(trimmedVal) && len(trimmedVal) > 1 && reHasAlpha.MatchString(trimmedVal) {
		fmt.Fprintf(df, "  constant %s : integer := 16#%s#;\n", name, trimmedVal)
	} else if trimmedVal == "y" {
		fmt.Fprintf(df, "  constant %s : integer := 1;\n", name)
	} else if trimmedVal == "n" {
		fmt.Fprintf(df, "  constant %s : integer := 0;\n", name)
	} else if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		fmt.Fprintf(df, "  constant %s : string := %s;\n", name, val)
	} else {
		// Assume decimal integer
		fmt.Fprintf(df, "  constant %s : integer := %s;\n", name, trimmedVal)
	}
}
