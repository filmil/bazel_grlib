package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type vhdlFile struct {
	path string
	std  string
}

func main() {
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceDir != "" {
		if err := os.Chdir(workspaceDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error changing to workspace directory: %v\n", err)
			os.Exit(1)
		}
	}

	grlibPath := "."
	if len(os.Args) > 1 {
		grlibPath = os.Args[1]
	}

	outPath := "third_party/grlib/grlib.BUILD"
	if len(os.Args) > 2 {
		outPath = os.Args[2]
	}

	libs, _ := readLines(filepath.Join(grlibPath, "lib/libs.txt"))
	var activeLibs []string
	libStds := make(map[string]string)
	for _, l := range libs {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.Fields(l)
		libName := parts[0]
		activeLibs = append(activeLibs, libName)
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "vhdlstd=") {
				libStds[libName] = strings.TrimPrefix(p, "vhdlstd=")
			}
		}
	}

	var discoveredLibs []string
	entries, _ := os.ReadDir(filepath.Join(grlibPath, "lib"))
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "tech" {
				techEntries, _ := os.ReadDir(filepath.Join(grlibPath, "lib/tech"))
				for _, te := range techEntries {
					if te.IsDir() {
						if _, err := os.Stat(filepath.Join(grlibPath, "lib/tech", te.Name(), "dirs.txt")); err == nil {
							discoveredLibs = append(discoveredLibs, "tech/"+te.Name())
						}
					}
				}
			} else {
				if _, err := os.Stat(filepath.Join(grlibPath, "lib", entry.Name(), "dirs.txt")); err == nil {
					discoveredLibs = append(discoveredLibs, entry.Name())
				}
			}
		}
	}

	var allLibs []string
	added := make(map[string]bool)
	addToAll := func(l string) {
		if !added[l] {
			allLibs = append(allLibs, l)
			added[l] = true
		}
	}
	addToAll("grlib")
	addToAll("techmap")
	for _, l := range activeLibs { addToAll(l) }
	for _, l := range discoveredLibs { addToAll(l) }

	libPathMap := make(map[string]string)
	for _, l := range allLibs {
		libPathMap[filepath.Base(l)] = l
	}

	libDeps := make(map[string][]string)
	libFiles := make(map[string][]string)

	for _, lib := range allLibs {
		libSourcePath := filepath.Join(grlibPath, "lib", lib)
		dirs, err := readLines(filepath.Join(libSourcePath, "dirs.txt"))
		if err != nil { continue }

		if lib == "gaisler" {
			dirs = append(dirs, "l5nv/shared", "noelv/pkg", "noelv/core", "noelv/subsys", "noelv")
		}

		var files []string
		depsSet := make(map[string]bool)
		for _, dir := range dirs {
			dir = strings.TrimSpace(dir)
			if dir == "" || strings.HasPrefix(dir, "#") { continue }

			for _, dp := range strings.Fields(dir) {
				if strings.HasPrefix(dp, "#") { break }
				subdirPath := filepath.Join(libSourcePath, dp)
				linesSyn, _ := readLines(filepath.Join(subdirPath, "vhdlsyn.txt"))
				linesSim, _ := readLines(filepath.Join(subdirPath, "vhdlsim.txt"))

				if len(linesSyn) == 0 && len(linesSim) == 0 {
					entries, _ := os.ReadDir(subdirPath)
					for _, entry := range entries {
						if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vhd") {
							linesSyn = append(linesSyn, entry.Name())
						}
					}
				}

				for _, fLine := range append(linesSyn, linesSim...) {
					fLine = strings.TrimSpace(fLine)
					if fLine == "" || strings.HasPrefix(fLine, "#") { continue }
					fParts := strings.Fields(fLine)
					fileName := fParts[0]
					
					relVhdPath := filepath.Join("lib", lib, dp, fileName)
					fullVhdPath := filepath.Join(grlibPath, relVhdPath)

					if _, err := os.Stat(fullVhdPath); err == nil {
						finalRef := relVhdPath
						if relVhdPath == "lib/grlib/stdlib/config.vhd" {
							finalRef = "@@//third_party/grlib:config.vhd"
						}

						// Avoid duplicates
						found := false
						for _, existing := range files {
							if existing == finalRef { found = true; break }
						}
						if !found {
							files = append(files, finalRef)
							fileDeps := parseLibDeps(fullVhdPath)
							for _, d := range fileDeps {
								if d != filepath.Base(lib) { depsSet[d] = true }
							}
						}
					}
				}
			}
		}

		if len(files) > 0 {
			libFiles[lib] = files
			var deps []string
			for d := range depsSet {
				if p, ok := libPathMap[d]; ok { deps = append(deps, p) }
			}
			libDeps[lib] = deps
		}
	}

	os.MkdirAll(filepath.Dir(outPath), 0755)
	gb, _ := os.Create(outPath)
	fmt.Fprintln(gb, "load(\"@rules_nvc//nvc:rules.bzl\", \"vhdl_library\")")
	fmt.Fprintln(gb, "")
	fmt.Fprintln(gb, "filegroup(")
	fmt.Fprintln(gb, "    name = \"grlib_srcs_all\",")
	fmt.Fprintln(gb, "    srcs = glob([\"**\"]),")
	fmt.Fprintln(gb, "    visibility = [\"//visibility:public\"],")
	fmt.Fprintln(gb, ")")
	fmt.Fprintln(gb, "")
	
	for _, lib := range allLibs {
		files := libFiles[lib]
		if len(files) == 0 { continue }

		libBase := filepath.Base(lib)
		fmt.Fprintln(gb, "# do not sort")
		fmt.Fprintf(gb, "filegroup(\n    name = \"%s_files\",\n    # do not sort\n    srcs = [\n", libBase)
		for _, f := range files {
			fmt.Fprintf(gb, "        \"%s\",\n", f)
		}
		fmt.Fprintf(gb, "    ],\n    visibility = [\"//visibility:public\"],\n)\n\n")

		fmt.Fprintf(gb, "vhdl_library(\n    name = \"%s\",\n    # do not sort\n    srcs = [\":%s_files\"],\n", libBase, libBase)
		std := libStds[lib]
		if std == "" { std = "1993" } else if std == "93" { std = "1993" } else if std == "08" { std = "2008" }
		fmt.Fprintf(gb, "    standard = \"%s\",\n    deps = [\n", std)
		sort.Strings(libDeps[lib])
		for _, d := range libDeps[lib] {
			fmt.Fprintf(gb, "        \":%s\",\n", filepath.Base(d))
		}
		fmt.Fprintf(gb, "    ],\n    visibility = [\"//visibility:public\"],\n)\n\n")
	}
	gb.Close()
	fmt.Println("Generated grlib.BUILD")
}

var libRegex = regexp.MustCompile(`(?i)^\s*library\s+([a-zA-Z0-9_]+)\s*;`)

func parseLibDeps(path string) []string {
	f, err := os.Open(path)
	if err != nil { return nil }
	defer f.Close()
	var deps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := libRegex.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			dep := strings.ToLower(match[1])
			if dep != "ieee" && dep != "std" && dep != "work" {
				deps = append(deps, dep)
			}
		}
	}
	return deps
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil { return nil, err }
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { lines = append(lines, scanner.Text()) }
	return lines, scanner.Err()
}
