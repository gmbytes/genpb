package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the configuration for the generator
type Config struct {
	ProtoDir    string // Input proto directory
	GoOutDir    string // Output directory for Go generated files
	CsOutDir    string // Output directory for C# generated files
	RustOutDir  string // Output directory for Rust generated files
	GoPkg       string // Go package name for generated files
	ProtocPath  string // Path to protoc executable
	ProtocGenGo string // Path to protoc-gen-go executable
	Language    string // Output language: go, Pb, rust, all
	Flag        string // Export flag: server (all files), client (exclude data_srv.proto, data_fwd.proto)
}

func main() {
	lang := flag.String("lang", "all", "Language: go, Pb, rust, all")
	goOut := flag.String("go_out", "./pb", "Go output directory")
	csOut := flag.String("cs_out", "./pb/Pb", "C# output directory")
	rustOut := flag.String("rust_out", "", "Rust output directory (for cmd_ext.rs, protocol.desc, protocol_meta.json)")
	protoIn := flag.String("proto_in", "./proto", "Proto input directory")
	toolsDir := flag.String("tools_dir", "", "Directory containing protoc and protoc-gen-go (default: ../proto)")
	flagType := flag.String("flag", "server", "Export flag: server (all files), client (exclude data_srv.proto, data_fwd.proto)")
	flag.Parse()

	toolsPath := *toolsDir
	if toolsPath == "" {
		toolsPath = "../proto"
	}
	toolsPath, _ = filepath.Abs(toolsPath)

	protocExe := "protoc"
	protocGenGoExe := "protoc-gen-go"
	if runtime.GOOS == "windows" {
		protocExe = "protoc.exe"
		protocGenGoExe = "protoc-gen-go.exe"
	}

	protoDir, _ := filepath.Abs(*protoIn)

	cfg := &Config{
		ProtoDir:    protoDir,
		GoOutDir:    *goOut,
		CsOutDir:    *csOut,
		RustOutDir:  *rustOut,
		GoPkg:       "server/pb",
		ProtocPath:  filepath.Join(toolsPath, protocExe),
		ProtocGenGo: filepath.Join(toolsPath, protocGenGoExe),
		Language:    *lang,
		Flag:        *flagType,
	}

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run executes the code generation process
func Run(cfg *Config) error {
	if cfg.Language == "go" || cfg.Language == "all" {
		goOutDir, _ := filepath.Abs(cfg.GoOutDir)
		if err := os.MkdirAll(goOutDir, 0755); err != nil {
			return fmt.Errorf("create Go output directory: %w", err)
		}
		cfg.GoOutDir = goOutDir

		if err := GenerateGo(cfg); err != nil {
			return fmt.Errorf("generate Go: %w", err)
		}
	}

	if cfg.Language == "Pb" || cfg.Language == "all" {
		csOutDir, _ := filepath.Abs(cfg.CsOutDir)
		if err := os.MkdirAll(csOutDir, 0755); err != nil {
			return fmt.Errorf("create C# output directory: %w", err)
		}
		cfg.CsOutDir = csOutDir

		if err := GenerateCSharp(cfg); err != nil {
			return fmt.Errorf("generate C#: %w", err)
		}
	}

	if cfg.Language == "rust" || cfg.Language == "all" {
		if cfg.RustOutDir == "" {
			if cfg.Language == "rust" {
				return fmt.Errorf("--rust_out is required when --lang=rust")
			}
		} else {
			rustOutDir, _ := filepath.Abs(cfg.RustOutDir)
			if err := os.MkdirAll(rustOutDir, 0755); err != nil {
				return fmt.Errorf("create Rust output directory: %w", err)
			}
			cfg.RustOutDir = rustOutDir

			if err := GenerateRust(cfg); err != nil {
				return fmt.Errorf("generate Rust: %w", err)
			}
		}
	}

	fmt.Println("Code generation completed successfully")
	return nil
}
