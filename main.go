package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the generator
type Config struct {
	ProtoDir    string // Input proto directory
	GoOutDir    string // Output directory for Go generated files
	CsOutDir    string // Output directory for C# generated files
	GoPkg       string // Go package name for generated files
	ProtocPath  string // Path to protoc executable
	ProtocGenGo string // Path to protoc-gen-go executable
	Language    string // Output language: go, cs, all
}

func main() {
	// Parse command line flags
	lang := flag.String("lang", "all", "Language: go, cs, all")
	goOut := flag.String("go_out", "./pb", "Go output directory")
	csOut := flag.String("cs_out", "./pb/cs", "C# output directory")
	flag.Parse()

	// Get absolute paths
	exeDir, _ := filepath.Abs("../proto")
	protoDir, _ := filepath.Abs("../proto")

	// Default paths
	cfg := &Config{
		ProtoDir:    protoDir,
		GoOutDir:    *goOut,
		CsOutDir:    *csOut,
		GoPkg:       "server/pb",
		ProtocPath:  filepath.Join(exeDir, "protoc.exe"),
		ProtocGenGo: filepath.Join(exeDir, "protoc-gen-go.exe"),
		Language:    *lang,
	}

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run executes the code generation process
func Run(cfg *Config) error {
	// Generate Go code
	if cfg.Language == "go" || cfg.Language == "all" {
		// Ensure Go output directory exists
		goOutDir, _ := filepath.Abs(cfg.GoOutDir)
		if err := os.MkdirAll(goOutDir, 0755); err != nil {
			return fmt.Errorf("create Go output directory: %w", err)
		}
		cfg.GoOutDir = goOutDir

		if err := GenerateGo(cfg); err != nil {
			return fmt.Errorf("generate Go: %w", err)
		}
	}

	// Generate C# code
	if cfg.Language == "cs" || cfg.Language == "all" {
		// Ensure C# output directory exists
		csOutDir, _ := filepath.Abs(cfg.CsOutDir)
		if err := os.MkdirAll(csOutDir, 0755); err != nil {
			return fmt.Errorf("create C# output directory: %w", err)
		}
		cfg.CsOutDir = csOutDir

		if err := GenerateCSharp(cfg); err != nil {
			return fmt.Errorf("generate C#: %w", err)
		}
	}

	fmt.Println("Code generation completed successfully")
	return nil
}
