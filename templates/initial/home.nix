# This is an empty template just used to setup the basic nix for linux.
# It should be overriden by the actual configuration in the setup scripts
{ config, pkgs, ... }:

{
  # Home Manager needs a bit of information about you and the paths it should
  # manage.
  home.username = "__USER__";
  home.homeDirectory = "__HOME__";

  home.stateVersion = "24.05"; # Please read the comment before changing.

  # The home.packages option allows you to install Nix packages into your
  # environment.
  home.packages = [
    pkgs.git
  ];

  home.sessionVariables = {
    CAMP_HOME = "˜/.camp";
  };

  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

}
