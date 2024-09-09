{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    go
  ];
  name = "mijnlink";
  shellHook = ''
    exec zsh
  '';
}
