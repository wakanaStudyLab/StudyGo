package main

import (
	"fmt"
	"os"
	"os/exec"
)

// ============================================================================
// Modern Go Crash Course - Main Runner
// For Rust / C# / Python / Java Developers
// ============================================================================

type module struct {
	path  string
	title string
}

func main() {
	printBanner("MODERN GO CRASH COURSE (For Rust / C# / Python / Java Developers)")

	modules := []module{
		{path: "./01_basics", title: "01: Basic Syntax, Types, Functions, and Loops"},
		{path: "./02_errors", title: "02: Error Handling, Custom Errors, and defer"},
		{path: "./03_interfaces", title: "03: Interfaces, Duck Typing, and Generics"},
		{path: "./04_lambdas_and_closures", title: "04: Anonymous Functions, Closures, and Go 1.22+ Per-Iteration"},
	}

	for _, mod := range modules {
		printSection(mod.title)
		cmd := exec.Command("go", "run", mod.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Module %s failed: %v\n", mod.path, err)
			os.Exit(1)
		}
	}

	printBanner("ALL GO TUTORIAL MODULES COMPLETED SUCCESSFULLY!")
}

func printBanner(title string) {
	fmt.Println("\n================================================================")
	fmt.Printf("  %s\n", title)
	fmt.Println("================================================================\n")
}

func printSection(title string) {
	fmt.Println("\n################################################################")
	fmt.Printf("# %s\n", title)
	fmt.Println("################################################################\n")
}
