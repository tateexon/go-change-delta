{ pkgs, scriptDir }:
with pkgs;
let
  go = pkgs.go_1_22;
in
mkShell {
  nativeBuildInputs = [
    bash
    git
    gnumake
    go
  ];

  CGO_ENABLED = "0";
  GOROOT = "${go}/share/go";

  shellHook = ''
    # setup go bin
    export GOBIN=$HOME/.nix-go/bin
    mkdir -p $GOBIN
    export PATH=$GOBIN:$PATH
  '';
}
