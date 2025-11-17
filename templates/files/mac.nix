# Mac system specific configurations go here.
{ config, pkgs, user, usersPath,  ... }:

{
  users.users.${user} = {
    home = usersPath;
  };
  nix.enable = false;
  nix.settings.experimental-features = "nix-command flakes";
  nix.nixPath = [];  # Disable channel lookups (using flakes instead)
  system.stateVersion = 5;
  nixpkgs.hostPlatform = "aarch64-darwin";
  system.primaryUser = user;
}
