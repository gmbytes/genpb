package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds code-generator configuration.
type Config struct {
	ProtoDir       string // Input proto directory
	GoOutDir       string // Output directory for Go generated files
	RustOutDir     string // Output directory for Rust generated files (pb.rs, typed_protocol.rs, …)
	GodotOutDir    string // Output directory for Godot bridge gen (godot_bridge_gen.rs); optional
	GoPkg          string // Go package name for generated files
	ProtocPath     string // Path to protoc executable
	ProtocGenGo    string // Path to protoc-gen-go executable
	ProtocGenProst string // Path to protoc-gen-prost executable
	Language       string // Output language: go | rust | all
	Flag           string // Export flag: server (all files) | client (exclude data_srv.proto, data_fwd.proto)
}

func main() {
	lang := flag.String("lang", "all", "Language to generate: go, rust, all")
	goOut := flag.String("go_out", "./pb", "Go output directory")
	rustOut := flag.String("rust_out", "", "Rust output directory (pb.rs, typed_protocol.rs, protocol_manifest.json, protocol.desc)")
	godotOut := flag.String("godot_out", "", "Godot bridge output directory (godot_bridge_gen.rs); optional, only used with --lang rust|all")
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
	protocGenProstExe := "protoc-gen-prost"
	if runtime.GOOS == "windows" {
		protocExe = "protoc.exe"
		protocGenGoExe = "protoc-gen-go.exe"
		protocGenProstExe = "protoc-gen-prost.exe"
	}

	protoDir, _ := filepath.Abs(*protoIn)

	cfg := &Config{
		ProtoDir:       protoDir,
		GoOutDir:       *goOut,
		RustOutDir:     *rustOut,
		GodotOutDir:    *godotOut,
		GoPkg:          "server/pb",
		ProtocPath:     filepath.Join(toolsPath, protocExe),
		ProtocGenGo:    filepath.Join(toolsPath, protocGenGoExe),
		ProtocGenProst: filepath.Join(toolsPath, protocGenProstExe),
		Language:       *lang,
		Flag:           *flagType,
	}

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run executes the code generation process.
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

	if cfg.Language == "rust" || cfg.Language == "all" {
		if cfg.RustOutDir == "" {
			if cfg.Language == "rust" {
				return fmt.Errorf("--rust_out is required when --lang=rust")
			}
			// lang=all without rust_out: skip silently
		} else {
			rustOutDir, _ := filepath.Abs(cfg.RustOutDir)
			if err := os.MkdirAll(rustOutDir, 0755); err != nil {
				return fmt.Errorf("create Rust output directory: %w", err)
			}
			cfg.RustOutDir = rustOutDir

			if cfg.GodotOutDir != "" {
				godotOutDir, _ := filepath.Abs(cfg.GodotOutDir)
				cfg.GodotOutDir = godotOutDir
			}

			if err := GenerateRust(cfg); err != nil {
				return fmt.Errorf("generate Rust: %w", err)
			}
		}
	}

	fmt.Println("Code generation completed successfully")
	return nil
}
