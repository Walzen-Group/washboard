let
  pkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-25.11.tar.gz") {};
in
pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs
    pkgs.nodePackages.pnpm
    pkgs.go
  ];
}
