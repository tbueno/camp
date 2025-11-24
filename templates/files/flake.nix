{
  description = "Camp Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-24.11-darwin";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-darwin = {
      url = "github:LnL7/nix-darwin/nix-darwin-24.11";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    home-manager = {
      url = "github:nix-community/home-manager/release-24.11";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Custom user-defined flakes
    {{- range .Flakes }}
    {{ .Name }} = {
      url = "{{ .URL }}";
      {{- range $key, $value := .Follows }}
      inputs.{{ $key }}.follows = "{{ $value }}";
      {{- end }}
    };
    {{- end }}
  };

  outputs = { self, nix-darwin, nixpkgs, nixpkgs-unstable, home-manager, {{ range .Flakes }}{{ .Name }}, {{ end }}... }:
  let
    configuration = { pkgs, ... }: {
      # not needed. They will be read from external files.
    };
    hostName = "{{.HostName}}";
    user = "{{.Name}}";
    platform = "{{.Platform}}";
    usersPath = "{{.HomeDir}}";
    architecture = "{{.Architecture}}";

    # Conditional logic to determine the system (darwin or linux)
    isDarwin = platform == "darwin";
    isLinux = !isDarwin;
    system = if isDarwin
      then "aarch64-darwin"
      else if architecture == "arm64"
        then "aarch64-linux"
        else "x86_64-linux";

    # Define variables that will be injected in other templates
    specialArgs = {
      inherit hostName user usersPath;
      customEnvVars = {
        {{- range $key, $value := .EnvVars }}
        "{{ $key }}" = "{{ $value }}";
        {{- end }}
      };
      customPackages = [
        {{- range .Packages }}
        "{{ . }}"
        {{- end }}
      ];
    };

  in
  {
    # macOS setup using nix-darwin (only if isDarwin is true)
    darwinConfigurations = if isDarwin then {
      ${hostName} = nix-darwin.lib.darwinSystem {
        inherit specialArgs system;
        modules = [
          {
            # Set Git commit hash for darwin-version
            system.configurationRevision = self.rev or self.dirtyRev or null;
            # Disable Nix path verification
            system.checks.verifyNixPath = false;
          }

          ./mac.nix

          # Custom system-level flake modules (nix-darwin)
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

          home-manager.darwinModules.home-manager {
            home-manager.useGlobalPkgs = true;
            home-manager.useUserPackages = true;
            home-manager.users.${user} = {
              imports = [
                ./modules/common.nix

                # Custom home-level flake modules
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
              ];
            };
            home-manager.extraSpecialArgs = specialArgs;
            home-manager.backupFileExtension = "backup";
          }
        ];
      };
    } else null;

    # Linux (Ubuntu) setup (only if isLinux is true)
    homeConfigurations = if isLinux then {
      "${user}" = home-manager.lib.homeManagerConfiguration {
        # Use nixpkgs with the correct system
        pkgs = import nixpkgs {
          inherit system;
        };

        # Pass specialArgs to home.nix
        extraSpecialArgs = specialArgs;

        modules = [
          ./linux.nix

          # Custom home-level flake modules
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
        ];
      };
    } else null;

    # Expose the package set for both darwin and linux for convenience
    darwinPackages = if isDarwin then self.darwinConfigurations.${hostName}.pkgs else null;
    linuxPackages = if isLinux then self.homeConfigurations.${user}.pkgs else null;
  };
}
