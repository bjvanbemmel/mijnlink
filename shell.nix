{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    go
    httpie
  ];
  name = "mijnlink";
  shellHook = ''
    exec zsh
  '';
}
