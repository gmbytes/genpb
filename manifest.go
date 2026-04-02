package main

// manifest.go – single source of truth for Rust/Godot codegen.
//
// Flow:
//   proto files ──► protoc --descriptor_set_out ──► protocol.desc
//   protocol.desc ──► buildProtocolManifest ──► ProtocolManifest
//   ProtocolManifest ──► writeTypedProtocol    → typed_protocol.rs
//                    ──► writeGodotBridgeGen   → godot_bridge_gen.rs
//                    ──► writeProtocolManifestJSON → protocol_manifest.json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"google.golang.org/protobuf/proto"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// ───────────────────────── Data model ─────────────────────────

// ProtocolManifest is serialised to protocol_manifest.json and consumed by
// gnet (typed routing) and gdbridge (GodotClass generation).
type ProtocolManifest struct {
	Version   int             `json:"version"`
	ProtoPkg  string          `json:"proto_pkg"`
	Messages  []ManifestEntry `json:"messages"`
	DataTypes []DataTypeEntry `json:"data_types"`
}

// ManifestEntry describes one EKey-bound command message.
type ManifestEntry struct {
	EKeyValue      int           `json:"ekey_value"`
	EKeyName       string        `json:"ekey_name"`
	MessageName    string        `json:"message_name"`
	MessageFull    string        `json:"message_full_name"`
	Direction      string        `json:"direction"`     // client_to_server | server_to_client
	Kind           string        `json:"kind"`          // req | rsp | dsp
	EventName      string        `json:"event_name"`    // snake_case event identifier
	Fields         []FieldSchema `json:"field_schema"`
	Fingerprint    string        `json:"fingerprint"`   // recursive schema string (pre-hash)
	FingerprintU64 uint64        `json:"fingerprint_u64"` // fnv1a64(Fingerprint)
	ClientExport   bool          `json:"client_export"`
	ServerExport   bool          `json:"server_export"`
	BridgeMode     string        `json:"bridge_mode"`       // typed_with_hotfix_fallback | typed_only
	HotfixFallback bool          `json:"hotfix_fallback"`
}

// DataTypeEntry is a non-command data message (from data.proto) used by the bridge.
type DataTypeEntry struct {
	MessageName string        `json:"message_name"`
	MessageFull string        `json:"message_full_name"`
	Fields      []FieldSchema `json:"field_schema"`
}

// FieldSchema holds per-field info needed by codegen.
type FieldSchema struct {
	Number   int32  `json:"number"`
	Name     string `json:"name"`      // proto field name (snake_case)
	TypeName string `json:"type_name"` // "bool","int64","string","e.EnumShort","m.MsgShort"
	Repeated bool   `json:"repeated"`
}

// ───────────────────────── Descriptor loading ─────────────────────────

// loadDescriptorSet reads a binary FileDescriptorSet.
func loadDescriptorSet(path string) (*dpb.FileDescriptorSet, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read descriptor set %s: %w", path, err)
	}
	var fds dpb.FileDescriptorSet
	if err := proto.Unmarshal(b, &fds); err != nil {
		return nil, fmt.Errorf("unmarshal descriptor set: %w", err)
	}
	return &fds, nil
}

// ───────────────────────── Descriptor index ─────────────────────────

// descIndex provides O(1) lookup of message descriptors by full or short name.
type descIndex struct {
	byFull  map[string]*dpb.DescriptorProto // "pb.RspLogin" → descriptor
	byShort map[string]*dpb.DescriptorProto // "RspLogin"    → descriptor (first match)
}

func buildDescIndex(fds *dpb.FileDescriptorSet) *descIndex {
	idx := &descIndex{
		byFull:  make(map[string]*dpb.DescriptorProto),
		byShort: make(map[string]*dpb.DescriptorProto),
	}
	for _, fd := range fds.File {
		pkg := fd.GetPackage()
		for _, m := range fd.MessageType {
			full := pkg + "." + m.GetName()
			idx.byFull[full] = m
			if _, ok := idx.byShort[m.GetName()]; !ok {
				idx.byShort[m.GetName()] = m
			}
			// Index nested messages (e.g. EKey.T lives inside EKey)
			for _, nm := range m.NestedType {
				nfull := full + "." + nm.GetName()
				idx.byFull[nfull] = nm
				if _, ok := idx.byShort[nm.GetName()]; !ok {
					idx.byShort[nm.GetName()] = nm
				}
			}
		}
	}
	return idx
}

// ───────────────────────── EKey parsing ─────────────────────────

// parseEKeyFromDescriptor extracts EKey.T enum values from cmd.proto descriptor.
// Returns map[enumName]numericValue, e.g. "ReqLogin" → 1.
func parseEKeyFromDescriptor(fds *dpb.FileDescriptorSet) (map[string]int32, error) {
	out := make(map[string]int32)
	for _, fd := range fds.File {
		if filepath.Base(fd.GetName()) != "cmd.proto" {
			continue
		}
		for _, m := range fd.MessageType {
			if m.GetName() != "EKey" {
				continue
			}
			for _, en := range m.EnumType {
				for _, ev := range en.Value {
					name := ev.GetName()
					num := ev.GetNumber()
					if name == "Invalid" || name == "Max" {
						continue
					}
					out[name] = num
				}
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no EKey enum values found – is cmd.proto in the descriptor set?")
	}
	return out, nil
}

// ───────────────────────── Manifest builder ─────────────────────────

// buildProtocolManifest derives a ProtocolManifest from a FileDescriptorSet.
func buildProtocolManifest(fds *dpb.FileDescriptorSet, flag string) (*ProtocolManifest, error) {
	ekeys, err := parseEKeyFromDescriptor(fds)
	if err != nil {
		return nil, err
	}
	idx := buildDescIndex(fds)

	type cmdSpec struct {
		file string
		kind string
		dir  string
	}
	specs := []cmdSpec{
		{"cmd_req.proto", "req", "client_to_server"},
		{"cmd_rsp.proto", "rsp", "server_to_client"},
		{"cmd_dsp.proto", "dsp", "server_to_client"},
	}

	var messages []ManifestEntry
	for _, spec := range specs {
		fd := fileByBase(fds, spec.file)
		if fd == nil {
			return nil, fmt.Errorf("descriptor set missing %s", spec.file)
		}
		pkg := fd.GetPackage()
		for _, mt := range fd.MessageType {
			name := mt.GetName()
			num, ok := ekeys[name]
			if !ok || num == 0 || num == 65535 {
				continue
			}
			full := pkg + "." + name
			fields := fieldsFromDescriptor(mt, idx)
			fpStr := recursiveFingerprintStr(full, idx, make(map[string]string))
			fpU64 := fnv1a64(fpStr)
			hotfix := spec.dir == "server_to_client"
			messages = append(messages, ManifestEntry{
				EKeyValue:      int(num),
				EKeyName:       name,
				MessageName:    name,
				MessageFull:    full,
				Direction:      spec.dir,
				Kind:           spec.kind,
				EventName:      toSnakeCase(name),
				Fields:         fields,
				Fingerprint:    fpStr,
				FingerprintU64: fpU64,
				ClientExport:   true,
				ServerExport:   true,
				BridgeMode:     "typed_with_hotfix_fallback",
				HotfixFallback: hotfix,
			})
		}
	}
	sort.Slice(messages, func(i, j int) bool { return messages[i].EKeyValue < messages[j].EKeyValue })

	// Data types: collect everything from data.proto (non-enum messages)
	var dataTypes []DataTypeEntry
	if dataFd := fileByBase(fds, "data.proto"); dataFd != nil {
		pkg := dataFd.GetPackage()
		for _, mt := range dataFd.MessageType {
			full := pkg + "." + mt.GetName()
			dataTypes = append(dataTypes, DataTypeEntry{
				MessageName: mt.GetName(),
				MessageFull: full,
				Fields:      fieldsFromDescriptor(mt, idx),
			})
		}
		sort.Slice(dataTypes, func(i, j int) bool { return dataTypes[i].MessageName < dataTypes[j].MessageName })
	}

	return &ProtocolManifest{
		Version:   1,
		ProtoPkg:  "pb",
		Messages:  messages,
		DataTypes: dataTypes,
	}, nil
}

// ───────────────────────── Field helpers ─────────────────────────

func fieldsFromDescriptor(m *dpb.DescriptorProto, idx *descIndex) []FieldSchema {
	var out []FieldSchema
	for _, f := range m.Field {
		out = append(out, FieldSchema{
			Number:   f.GetNumber(),
			Name:     f.GetName(),
			TypeName: protoFieldTypeName(f),
			Repeated: f.GetLabel() == dpb.FieldDescriptorProto_LABEL_REPEATED,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out
}

// protoFieldTypeName returns a normalised type string:
//   - scalar: "bool", "int32", "sint64", "string", etc.
//   - enum:   "e.ShortName"
//   - message: "m.ShortName"
func protoFieldTypeName(f *dpb.FieldDescriptorProto) string {
	switch f.GetType() {
	case dpb.FieldDescriptorProto_TYPE_MESSAGE:
		return "m." + lastSegment(strings.TrimPrefix(f.GetTypeName(), "."))
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		return "e." + lastSegment(strings.TrimPrefix(f.GetTypeName(), "."))
	default:
		return scalarTypeName(f.GetType())
	}
}

func lastSegment(s string) string {
	if i := strings.LastIndex(s, "."); i >= 0 {
		return s[i+1:]
	}
	return s
}

func scalarTypeName(t dpb.FieldDescriptorProto_Type) string {
	switch t {
	case dpb.FieldDescriptorProto_TYPE_DOUBLE:
		return "double"
	case dpb.FieldDescriptorProto_TYPE_FLOAT:
		return "float"
	case dpb.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case dpb.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case dpb.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case dpb.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case dpb.FieldDescriptorProto_TYPE_SINT32:
		return "sint32"
	case dpb.FieldDescriptorProto_TYPE_SINT64:
		return "sint64"
	case dpb.FieldDescriptorProto_TYPE_FIXED32:
		return "fixed32"
	case dpb.FieldDescriptorProto_TYPE_FIXED64:
		return "fixed64"
	case dpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "sfixed32"
	case dpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "sfixed64"
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case dpb.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		return "bytes"
	default:
		return fmt.Sprintf("unknown_%d", t)
	}
}

// ───────────────────────── Recursive fingerprint ─────────────────────────

// recursiveFingerprintStr builds a canonical schema string for a message,
// recursing into nested message types. This catches schema drift in sub-messages.
func recursiveFingerprintStr(fullMsgName string, idx *descIndex, memo map[string]string) string {
	if s, ok := memo[fullMsgName]; ok {
		return s // already computed (handles cycles with cycle-break via first-seen)
	}
	memo[fullMsgName] = "?" // cycle breaker

	m, ok := idx.byFull[fullMsgName]
	if !ok {
		short := lastSegment(fullMsgName)
		m, ok = idx.byShort[short]
	}
	if !ok {
		return "?"
	}

	var parts []string
	for _, f := range m.Field {
		tn := protoFieldTypeName(f)
		repeated := f.GetLabel() == dpb.FieldDescriptorProto_LABEL_REPEATED
		var part string
		if strings.HasPrefix(tn, "m.") {
			subFull := strings.TrimPrefix(f.GetTypeName(), ".")
			inner := recursiveFingerprintStr(subFull, idx, memo)
			subShort := lastSegment(subFull)
			if repeated {
				part = fmt.Sprintf("%d:r:m:%s{%s}", f.GetNumber(), subShort, inner)
			} else {
				part = fmt.Sprintf("%d:m:%s{%s}", f.GetNumber(), subShort, inner)
			}
		} else {
			if repeated {
				part = fmt.Sprintf("%d:r:%s", f.GetNumber(), tn)
			} else {
				part = fmt.Sprintf("%d:%s", f.GetNumber(), tn)
			}
		}
		parts = append(parts, part)
	}
	sort.Strings(parts)
	result := strings.Join(parts, ";")
	memo[fullMsgName] = result
	return result
}

// fnv1a64 computes FNV-1a 64-bit hash of a string.
func fnv1a64(s string) uint64 {
	const offset uint64 = 14695981039346656037
	const prime uint64 = 1099511628211
	h := offset
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime
	}
	return h
}

// ───────────────────────── File helpers ─────────────────────────

func fileByBase(fds *dpb.FileDescriptorSet, base string) *dpb.FileDescriptorProto {
	for _, fd := range fds.File {
		if filepath.Base(fd.GetName()) == base {
			return fd
		}
	}
	return nil
}

// ───────────────────────── JSON serialisation ─────────────────────────

func writeProtocolManifestJSON(path string, m *ProtocolManifest) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

// ───────────────────────── Shared naming utilities ─────────────────────────
// (used by gen_rust.go for code emission)

func toUpperCamelCase(s string) string {
	parts := splitIdent(s)
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		upper := strings.ToUpper(p)
		if upper == "ZZZ" || upper == "ID" || upper == "OK" {
			// capitalize first letter only, lowercase rest (e.g. "ZZZ" → "Zzz")
			r := []rune(strings.ToLower(p))
			r[0] = unicode.ToUpper(r[0])
			b.WriteString(string(r))
		} else {
			runes := []rune(strings.ToLower(p))
			runes[0] = unicode.ToUpper(runes[0])
			b.WriteString(string(runes))
		}
	}
	return b.String()
}

func toSnakeCase(s string) string {
	parts := splitIdent(s)
	var lower []string
	for _, p := range parts {
		if p != "" {
			lower = append(lower, strings.ToLower(p))
		}
	}
	return strings.Join(lower, "_")
}

func splitIdent(s string) []string {
	var parts []string
	var cur []rune
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prevLower := unicode.IsLower(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || (nextLower && len(cur) > 0) {
				parts = append(parts, string(cur))
				cur = nil
			}
		}
		if r == '_' {
			if len(cur) > 0 {
				parts = append(parts, string(cur))
				cur = nil
			}
			continue
		}
		cur = append(cur, r)
	}
	if len(cur) > 0 {
		parts = append(parts, string(cur))
	}
	return parts
}
