{ config, lib, pkgs, user, usersPath, customEnvVars, ... }:

{
  programs.home-manager.enable = true;

  programs.direnv = {
    enable = true;
    enableZshIntegration = true;
    enableBashIntegration = true;
    nix-direnv.enable = true;
  };

  home = {
    homeDirectory = "${usersPath}";
    packages = with pkgs; [
      devbox
      direnv
    ];
    stateVersion = "24.05";
    username = user;

    sessionVariables = customEnvVars;
  };

  programs.zsh = {
    enable = true;
    dotDir = ".camp";

    # These commands will be added to .zshrc file. They are executed after the shell is initialized with .zshenv
    initExtra = ''
        [ -f ~/.zshrc ] && source ~/.zshrc
    '';
  };
}
