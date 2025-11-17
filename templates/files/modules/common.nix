{ config, lib, pkgs, user, usersPath, customEnvVars, ... }:

{
  programs.home-manager.enable = true;

  programs.direnv = {
    enable = true;
    enableZshIntegration = true;   # Enable zsh integration for direnv
    enableBashIntegration = false;
    nix-direnv.enable = true;
  };

  home = {
    homeDirectory = "${usersPath}";
    packages = with pkgs; [
      devbox
      direnv
      git  # Add git from Nix to ensure it's available
    ];
    stateVersion = "24.05";
    username = user;

    # Session variables managed by home-manager through zsh
    sessionVariables = customEnvVars;
  };

  # Enable zsh management with dotDir approach
  programs.zsh = {
    enable = true;
    dotDir = "${config.home.homeDirectory}/.camp";

    # Source user's original .zshrc after camp's config loads
    initExtra = ''
      [ -f ~/.zshrc ] && source ~/.zshrc
    '';
  };
}
