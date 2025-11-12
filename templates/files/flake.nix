{
  description = "Camp Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-darwin = {
      url = "github:LnL7/nix-darwin";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nix-darwin, nixpkgs, nixpkgs-unstable, home-manager, ... }:
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
    };

  in
  {
    # macOS setup using nix-darwin (only if isDarwin is true)
    darwinConfigurations = if isDarwin then {
      ${hostName} = nix-darwin.lib.darwinSystem {
        inherit specialArgs;
        system = {
          configurationRevision = self.rev or self.dirtyRev or null;
          checks.verifyNixPath = false;
        };

        modules = [
          ./mac.nix

          home-manager.darwinModules.home-manager {
            home-manager.useGlobalPkgs = true;
            home-manager.useUserPackages = true;
            home-manager.users.${user} = {
              imports = [
                ./modules/common.nix
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
        ];
      };
    } else null;

    # Expose the package set for both darwin and linux for convenience
    darwinPackages = if isDarwin then self.darwinConfigurations.${hostName}.pkgs else null;
    linuxPackages = if isLinux then self.homeConfigurations.${user}.pkgs else null;
  };
}
