{
  description = "development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs@{ self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ ];
        };

        # set go version here that will be used in all shells
        go = pkgs.go_1_23;

        scriptDir = toString ./.;  # Converts the flake's root directory to a string

        # Importing the shell environments from separate files
        fullEnv = pkgs.callPackage ./nix/devshell.nix {
          inherit pkgs scriptDir go;
        };

        ciEnv = pkgs.callPackage ./nix/ci-runtests.nix {
          inherit pkgs scriptDir go;
        };
      in rec {
        devShell = fullEnv;
        devShells = {
          ci-runtests = ciEnv;
        };

        formatter = pkgs.nixpkgs-fmt;
      });
}
