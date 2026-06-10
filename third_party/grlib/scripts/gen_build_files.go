package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	// Assuming the script is run while the files are still present to generate the BUILD file
	libs, _ := readLines("lib/libs.txt")
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
	entries, _ := os.ReadDir("lib")
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "tech" {
				techEntries, _ := os.ReadDir("lib/tech")
				for _, te := range techEntries {
					if te.IsDir() {
						if _, err := os.Stat(filepath.Join("lib/tech", te.Name(), "dirs.txt")); err == nil {
							discoveredLibs = append(discoveredLibs, "tech/"+te.Name())
						}
					}
				}
			} else {
				if _, err := os.Stat(filepath.Join("lib", entry.Name(), "dirs.txt")); err == nil {
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
	for _, l := range activeLibs {
		addToAll(l)
	}
	for _, l := range discoveredLibs {
		addToAll(l)
	}

	libPathMap := make(map[string]string)
	for _, l := range allLibs {
		libPathMap[filepath.Base(l)] = l
	}

	libDeps := make(map[string][]string)
	libFiles := make(map[string][]vhdlFile)

	for _, lib := range allLibs {
		libPath := filepath.Join("lib", lib)
		dirs, err := readLines(filepath.Join(libPath, "dirs.txt"))
		if err != nil {
			continue
		}

		if lib == "gaisler" {
			// specifically add components that are often missing from dirs.txt or needed for noelv
			dirs = append(dirs, "l5nv/shared")
			dirs = append(dirs, "noelv/pkg")
			dirs = append(dirs, "noelv/core")
			dirs = append(dirs, "noelv/subsys")
			dirs = append(dirs, "noelv")
		}

		var files []vhdlFile
		depsSet := make(map[string]bool)
		for _, dir := range dirs {
			dir = strings.TrimSpace(dir)
			if dir == "" || strings.HasPrefix(dir, "#") {
				continue
			}

			for _, dp := range strings.Fields(dir) {
				if strings.HasPrefix(dp, "#") {
					break
				}
				subdirPath := filepath.Join(libPath, dp)
				linesSyn, _ := readLines(filepath.Join(subdirPath, "vhdlsyn.txt"))
				linesSim, _ := readLines(filepath.Join(subdirPath, "vhdlsim.txt"))

				if len(linesSyn) == 0 && len(linesSim) == 0 {
					// No source list files, try to glob all .vhd files in this subdir
					entries, _ := os.ReadDir(subdirPath)
					for _, entry := range entries {
						if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vhd") {
							linesSyn = append(linesSyn, entry.Name())
						}
					}
				}

				for _, fLine := range append(linesSyn, linesSim...) {
					fLine = strings.TrimSpace(fLine)
					if fLine == "" || strings.HasPrefix(fLine, "#") {
						continue
					}
					fParts := strings.Fields(fLine)
					fileName := filepath.Join(dp, fParts[0])
					fullPath := filepath.Join(libPath, fileName)
					if _, err := os.Stat(fullPath); err == nil {
						found := false
						for _, existing := range files {
							if existing.path == fileName {
								found = true
								break
							}
						}
						if !found {
							files = append(files, vhdlFile{path: fileName})
							fileDeps := parseLibDeps(fullPath)
							for _, d := range fileDeps {
								if d != filepath.Base(lib) {
									depsSet[d] = true
								}
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
				if p, ok := libPathMap[d]; ok {
					deps = append(deps, p)
				}
			}
			libDeps[lib] = deps
		}
	}

	// Generate grlib.BUILD
	gb, _ := os.Create("third_party/grlib/grlib.BUILD")
	fmt.Fprintln(gb, "load(\"@rules_nvc//nvc:rules.bzl\", \"vhdl_library\")")
	fmt.Fprintln(gb, "")

	for _, lib := range allLibs {
		files := libFiles[lib]
		if len(files) == 0 {
			continue
		}

		libBase := filepath.Base(lib)

		// Filegroup
		fmt.Fprintln(gb, "# do not sort")
		fmt.Fprintf(gb, "filegroup(\n")
		fmt.Fprintf(gb, "    name = \"%s_files\",\n", libBase)
		fmt.Fprintf(gb, "    srcs = [\n")
		for _, f := range files {
			filePath := fmt.Sprintf("lib/%s/%s", lib, f.path)
			if filePath == "lib/grlib/stdlib/config.vhd" {
				// Use the generated one instead
				fmt.Fprintf(gb, "        \"@@//third_party/grlib:config.vhd\",\n")
			} else {
				fmt.Fprintf(gb, "        \"%s\",\n", filePath)
			}
		}
		fmt.Fprintf(gb, "    ],\n")
		fmt.Fprintf(gb, "    visibility = [\"//visibility:public\"],\n")
		fmt.Fprintf(gb, ")\n\n")

		// vhdl_library
		fmt.Fprintf(gb, "vhdl_library(\n")
		fmt.Fprintf(gb, "    name = \"%s\",\n", libBase)
		fmt.Fprintf(gb, "    srcs = [\":%s_files\"],\n", libBase)

		std := libStds[lib]
		if std == "" {
			std = "1993"
		} else if std == "93" {
			std = "1993"
		} else if std == "08" {
			std = "2008"
		}
		fmt.Fprintf(gb, "    standard = \"%s\",\n", std)

		fmt.Fprintf(gb, "    deps = [\n")
		for _, d := range libDeps[lib] {
			fmt.Fprintf(gb, "        \":%s\",\n", filepath.Base(d))
		}
		fmt.Fprintf(gb, "    ],\n")
		fmt.Fprintf(gb, "    visibility = [\"//visibility:public\"],\n")
		fmt.Fprintf(gb, ")\n\n")
	}
	gb.Close()
	fmt.Println("Generated grlib.BUILD")
}

var libRegex = regexp.MustCompile(`(?i)^\s*library\s+([a-zA-Z0-9_]+)\s*;`)

func parseLibDeps(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var deps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		match := libRegex.FindStringSubmatch(line)
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
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
