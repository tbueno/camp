# This is an empty template just used to setup the basic nix for linux.
# It should be overriden by the actual configuration in the setup scripts
{ config, pkgs, ... }:

{
  home.stateVersion = "24.05"; # Please read the comment before changing.

  # The home.packages option allows you to install Nix packages into your
  # environment.
  home.packages = [
    pkgs.git
    pkgs.direnv
    pkgs.devbox
  ];

  home.sessionVariables = {
    CAMP_HOME = "˜/.camp";
  };

  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

}
