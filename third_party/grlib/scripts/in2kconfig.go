package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceDir != "" {
		os.Chdir(workspaceDir)
	}

	grlibPath := "."
	if len(os.Args) > 1 {
		grlibPath = os.Args[1]
	}

	outDir := "third_party/grlib/kconfig"
	if len(os.Args) > 2 {
		outDir = os.Args[2]
	}
	os.MkdirAll(outDir, 0755)

	err := filepath.Walk(grlibPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".in") {
			relPath, _ := filepath.Rel(grlibPath, path)
			targetPath := filepath.Join(outDir, relPath+".Kconfig")
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating dir for %q: %v\n", targetPath, err)
				return err
			}
			fmt.Printf("Converting %s to %s\n", path, targetPath)
			convertInToKconfig(grlibPath, outDir, path, targetPath)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking files: %v\n", err)
	}
}

var (
	reBool      = regexp.MustCompile(`(?i)^\s*bool\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)`)
	reInt       = regexp.MustCompile(`(?i)^\s*int\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)\s*(.*)`)
	reHex       = regexp.MustCompile(`(?i)^\s*hex\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)\s*(.*)`)
	reSource    = regexp.MustCompile(`(?i)^\s*source\s+(.*)`)
	reDefine    = regexp.MustCompile(`(?i)^\s*define_(bool|int|hex|string)\s+(CONFIG_[a-zA-Z0-9_]+)\s+(.*)`)
)

func convertInToKconfig(grlibPath, outDir, src, dest string) {
	sf, err := os.Open(src)
	if err != nil { return }
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil { return }
	defer df.Close()

	// Calculate a prefix based on the source path RELATIVE TO grlibPath
	rel, _ := filepath.Rel(grlibPath, src)
	prefix := strings.ToUpper(rel)
	prefix = regexp.MustCompile(`[^A-Z0-9]`).ReplaceAllString(prefix, "_")
	prefix = strings.TrimSuffix(prefix, "_IN")
	prefix = strings.TrimSuffix(prefix, "_VHD")
	prefix = strings.Trim(prefix, "_")

	scanner := bufio.NewScanner(sf)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			fmt.Fprintln(df, line)
			continue
		}

		handleSymbol := func(rawName, kind, label, def string) {
			name := strings.TrimPrefix(rawName, "CONFIG_")
			
			fullName := fmt.Sprintf("%s_%s", prefix, name)
			
			fmt.Fprintf(df, "config %s\n", fullName)
			if label != "" {
				fmt.Fprintf(df, "    %s \"%s\"\n", kind, label)
			} else {
				fmt.Fprintf(df, "    %s\n", kind)
			}
			if def != "" {
				fmt.Fprintf(df, "    default %s\n", def)
			}
		}

		if m := reBool.FindStringSubmatch(trimmed); m != nil {
			handleSymbol(m[2], "bool", m[1], "")
		} else if m := reInt.FindStringSubmatch(trimmed); m != nil {
			def := "0"
			if m[3] != "" { def = m[3] }
			handleSymbol(m[2], "int", m[1], def)
		} else if m := reHex.FindStringSubmatch(trimmed); m != nil {
			def := "0"
			if m[3] != "" { def = m[3] }
			handleSymbol(m[2], "hex", m[1], def)
		} else if m := reDefine.FindStringSubmatch(trimmed); m != nil {
			typ := m[1]
			def := m[3]
			if typ == "bool" {
				if def == "y" { def = "y" } else { def = "n" }
			}
			handleSymbol(m[2], typ, "", def)
		} else if m := reSource.FindStringSubmatch(trimmed); m != nil {
			// Skip generating rsource statements. We'll handle inclusion flatly in gen_master_kconfig.
			fmt.Fprintf(df, "# Original source: %s\n", m[1])
		} else {
			fmt.Fprintf(df, "# %s\n", line)
		}
	}
}
