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
          config = {
            permittedInsecurePackages = [
              "python-2.7.18.8"
            ];
          };
        };

        scriptDir = toString ./.;  # Converts the flake's root directory to a string

        # Importing the shell environments from separate files
        fullEnv = pkgs.callPackage ./nix/devshell.nix {
          inherit pkgs scriptDir;
        };

        ciEnv = pkgs.callPackage ./nix/ci-runtests.nix {
          inherit pkgs scriptDir;
        };
      in rec {
        devShell = fullEnv;
        devShells = {
          ci-runtests = ciEnv;
        };

        formatter = pkgs.nixpkgs-fmt;
      });
}
