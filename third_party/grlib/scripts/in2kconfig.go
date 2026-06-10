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
		if err != nil { return err }
		if !info.IsDir() && strings.HasSuffix(path, ".in") {
			relPath, _ := filepath.Rel(grlibPath, path)
			targetPath := filepath.Join(outDir, relPath+".Kconfig")
			os.MkdirAll(filepath.Dir(targetPath), 0755)
			convertInToKconfig(path, targetPath)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking files: %v\n", err)
	}
}

var (
	reMainmenu  = regexp.MustCompile(`(?i)^mainmenu_name\s+"(.*)"`)
	reComment   = regexp.MustCompile(`(?i)^\s*comment\s+'(.*)'`)
	reBool      = regexp.MustCompile(`(?i)^\s*bool\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)`)
	reInt       = regexp.MustCompile(`(?i)^\s*int\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)\s*(.*)`)
	reHex       = regexp.MustCompile(`(?i)^\s*hex\s+'(.*)'\s+(CONFIG_[a-zA-Z0-9_]+)\s*(.*)`)
	reSource    = regexp.MustCompile(`(?i)^\s*source\s+(.*)`)
	reIf        = regexp.MustCompile(`(?i)^\s*if\s+\[\s*(.*)\s*\]\s*;\s*then`)
	reFi        = regexp.MustCompile(`(?i)^\s*fi\b`)
	reDefine    = regexp.MustCompile(`(?i)^\s*define_(bool|int|hex|string)\s+(CONFIG_[a-zA-Z0-9_]+)\s+(.*)`)
	reChoice    = regexp.MustCompile(`(?i)^\s*choice\s+'(.*)'\s*\\`)
)

func convertInToKconfig(src, dest string) {
	sf, err := os.Open(src)
	if err != nil { return }
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil { return }
	defer df.Close()

	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			fmt.Fprintln(df, line)
			continue
		}

		if m := reBool.FindStringSubmatch(trimmed); m != nil {
			name := strings.TrimPrefix(m[2], "CONFIG_")
			fmt.Fprintf(df, "config %s\n    bool \"%s\"\n", name, m[1])
		} else if m := reInt.FindStringSubmatch(trimmed); m != nil {
			name := strings.TrimPrefix(m[2], "CONFIG_")
			def := "0"
			if m[3] != "" { def = m[3] }
			fmt.Fprintf(df, "config %s\n    int \"%s\"\n    default %s\n", name, m[1], def)
		} else if m := reHex.FindStringSubmatch(trimmed); m != nil {
			name := strings.TrimPrefix(m[2], "CONFIG_")
			def := "0"
			if m[3] != "" { def = m[3] }
			fmt.Fprintf(df, "config %s\n    hex \"%s\"\n    default %s\n", name, m[1], def)
		} else if m := reDefine.FindStringSubmatch(trimmed); m != nil {
			name := strings.TrimPrefix(m[2], "CONFIG_")
			typ := m[1]
			def := m[3]
			if typ == "bool" {
				if def == "y" { def = "y" } else { def = "n" }
			}
			fmt.Fprintf(df, "config %s\n    %s\n    default %s\n", name, typ, def)
		} else if m := reSource.FindStringSubmatch(trimmed); m != nil {
			sfile := m[1]
			fmt.Fprintf(df, "rsource \"%s.Kconfig\"\n", filepath.Join("../../../..", sfile))
		} else {
			// Comment out everything else to ensure valid Kconfig
			fmt.Fprintf(df, "# %s\n", line)
		}
	}
}
