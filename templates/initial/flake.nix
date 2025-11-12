# This is an empty template just used to setup the basic nix installation
# It should be overriden by the actual configuration in the setup scripts
{
  description = "Initial system flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";  # Use the same nixpkgs as the flake
    };
  };

  outputs = { self, nixpkgs, home-manager }: let
    system = "__SYSTEM__";  # Will be replaced with actual system
  in {
    homeConfigurations = {
      "__USER__" = home-manager.lib.homeManagerConfiguration {
        pkgs = import nixpkgs {
          inherit system;
        };
        modules = [
          ./home.nix  # Your home.nix configuration
        ];
      };
    };
  };
}
