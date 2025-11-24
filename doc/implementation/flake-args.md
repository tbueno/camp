# Flake Arguments Implementation Plan

## Overview
Implement per-flake argument passing with automatic user data (userName, hostName, home) + custom arguments, inferred from YAML types, passed via flake output function parameters.

## Design Decisions

### 1. Argument Scope
**Per-flake arguments** - Each flake can have its own set of arguments in camp.yml. More flexible, allows passing different values to different flakes.

### 2. Data Source
**Automatic + custom arguments** - Auto-pass standard fields (userName, hostName, home) from camp's User struct, and allow additional custom arguments like email.

### 3. Nix Mechanism
**Via flake outputs as function parameters** - External flakes define outputs as functions that accept parameters. More explicit, requires flakes to declare params.

### 4. Type System
**Type inference from YAML** - Support multiple types (string, bool, number, list) by inferring from YAML parser, avoiding verbose type declarations.

## Part A: camp.yml Schema

### New YAML Structure
```yaml
flakes:
  - name: my-personal-config
    url: "github:tbueno/nix-config"
    args:
      email: "tbueno@gmail.com"           # string (inferred)
      enableDevTools: true                 # bool (inferred)
      fontSize: 14                         # number (inferred)
      packages: [vim, git, tmux]           # list (inferred)
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
```

### Automatic Arguments
Always passed to all flake outputs:
- `userName` - from `User.Name`
- `hostName` - from `User.HostName`
- `home` - from `User.HomeDir`

### Custom Arguments
User-defined per-flake in the `args` map. Types inferred from YAML:
- **String**: Quoted values (`"text"`)
- **Bool**: `true` / `false`
- **Number**: `42`, `3.14`
- **List**: `[item1, item2]` or YAML sequence syntax

## Part B: Code Implementation

### Phase 1: Update Data Structures (`internal/system/types.go`)

**Changes:**
- Add `Args map[string]interface{} \`yaml:"args"\`` field to `Flake` struct
- YAML parser will automatically unmarshal values with correct Go types

**Expected behavior:**
- YAML strings → Go `string`
- YAML bools → Go `bool`
- YAML integers → Go `int`
- YAML floats → Go `float64`
- YAML sequences → Go `[]interface{}`

### Phase 2: Add Validation (`internal/system/config.go`)

**New function: `ValidateFlakeArgs()`**

Validation rules:
1. **Valid Nix identifiers**: Arg names must be alphanumeric + hyphens/underscores only
2. **No reserved names**: Cannot use `userName`, `hostName`, `home` (auto-provided)
3. **Supported types only**: string, bool, int, int64, float64, []interface{}
4. **Reject unsupported types**: maps, nil, complex types

**Integration:**
- Call `ValidateFlakeArgs()` from existing `ValidateFlakes()` function
- Return clear error messages for each validation failure

### Phase 3: Template Rendering (`internal/system/template.go`)

**New helper function: `renderNixValue(value interface{}) string`**

Type-based rendering:
- `string` → `"value"` (escape quotes, backslashes, newlines)
- `bool` → `true` / `false` (no quotes)
- `int`, `int64` → `42` (no quotes)
- `float64` → `3.14` (no quotes)
- `[]interface{}` → `[ "item1" "item2" ]` (recurse for each element)

**String escaping:**
- `"` → `\"`
- `\` → `\\`
- `\n` → `\\n`

**Template helper:**
- Add `renderNixValue` to template function map
- Make available in flake.nix template

### Phase 4: Update Flake Template (`templates/files/flake.nix`)

**Modify output imports to pass arguments:**

For **system outputs** (nix-darwin modules):
```nix
{{- range $flake := .Flakes }}
  {{- range .Outputs }}
    {{- if eq .Type "system" }}
({{ $flake.Name }}.{{ .Name }} {
  userName = "{{ $.Name }}";
  hostName = "{{ $.HostName }}";
  home = "{{ $.HomeDir }}";
  {{- range $key, $value := $flake.Args }}
  {{ $key }} = {{ renderNixValue $value }};
  {{- end }}
})
    {{- end }}
  {{- end }}
{{- end }}
```

For **home outputs** (home-manager modules):
```nix
{{- range $flake := .Flakes }}
  {{- range .Outputs }}
    {{- if eq .Type "home" }}
({{ $flake.Name }}.{{ .Name }} {
  userName = "{{ $.Name }}";
  hostName = "{{ $.HostName }}";
  home = "{{ $.HomeDir }}";
  {{- range $key, $value := $flake.Args }}
  {{ $key }} = {{ renderNixValue $value }};
  {{- end }}
})
    {{- end }}
  {{- end }}
{{- end }}
```

### Phase 5: Add Tests

**`internal/system/config_test.go`:**
- Test YAML unmarshaling preserves types
- Test arg name validation (valid names, invalid chars)
- Test reserved name detection (userName, hostName, home)
- Test type validation (supported vs unsupported types)
- Test validation error messages

**`internal/system/template_test.go`:**
- Test `renderNixValue()` for each type:
  - String: basic, with quotes, with special chars
  - Bool: true, false
  - Numbers: int, float
  - List: empty, strings, numbers, mixed
- Test automatic args (userName, hostName, home) injection
- Test custom args rendering in template
- Test full integration: config → template → rendered Nix

### Phase 6: Documentation

**Update `CLAUDE.md`:**
- Add "Flake Arguments" section after "Flake System"
- Document YAML schema with examples
- Explain automatic vs custom arguments
- Show type inference from YAML
- Provide external flake examples

**Create example template:**
- Add `templates/flakes/parameterized-example.yml`
- Show real-world usage with all arg types

## Part C: External Flake Pattern

### Example: Parameterized Flake Output

**Your flake** (`github:tbueno/nix-config/flake.nix`):

```nix
{
  description = "Thiago Bueno's Nix configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-darwin.url = "github:LnL7/nix-darwin";
    home-manager.url = "github:nix-community/home-manager";
  };

  outputs = { self, nixpkgs, nix-darwin, home-manager }:
  {
    # System-level module (nix-darwin)
    darwinModules.default = { userName, hostName, home, email, ... }@args: {
      networking.hostName = hostName;

      users.users.${userName} = {
        name = userName;
        home = home;
      };

      home-manager.users.${userName} = {
        home.sessionVariables = {
          EMAIL = email;
        };
      };
    };

    # User-level module (home-manager)
    homeManagerModules.default = { userName, home, email, ... }@args: {
      home.username = userName;
      home.homeDirectory = home;

      programs.git = {
        enable = true;
        userEmail = email;
      };
    };
  };
}
```

**Your camp.yml:**
```yaml
flakes:
  - name: personal-config
    url: "github:tbueno/nix-config"
    args:
      email: "tbueno@gmail.com"
    outputs:
      - name: darwinModules.default
        type: system
      - name: homeManagerModules.default
        type: home
```

**Generated flake.nix** (excerpt):
```nix
modules = [
  ./mac.nix

  # Custom system-level flake modules
  (personal-config.darwinModules.default {
    userName = "bueno";
    hostName = "macbook";
    home = "/Users/bueno";
    email = "tbueno@gmail.com";
  })
];
```

## Implementation Order

1. ✅ **Phase 1**: Update `types.go` (add Args field)
2. ✅ **Phase 2**: Update `config.go` (validation)
3. ✅ **Phase 3**: Update `template.go` (renderNixValue helper)
4. ✅ **Phase 4**: Update `flake.nix` template
5. ✅ **Phase 5**: Add comprehensive tests
6. ✅ **Phase 6**: Update documentation

## Files to Modify

- `internal/system/types.go`
- `internal/system/config.go`
- `internal/system/template.go`
- `templates/files/flake.nix`
- `internal/system/config_test.go`
- `internal/system/template_test.go`
- `CLAUDE.md`
- `templates/flakes/parameterized-example.yml` (new file)

## Testing Strategy

**Unit tests:**
- Each validation rule independently
- Each type rendering case
- Edge cases (empty strings, special chars, empty lists)

**Integration tests:**
- Full flow: YAML → validation → template → rendered Nix
- Multiple flakes with different arg types
- Flakes with and without args

**Manual testing:**
- Create test camp.yml with args
- Run `camp env rebuild`
- Verify generated `~/.camp/nix/flake.nix`
- Test with actual external flake
