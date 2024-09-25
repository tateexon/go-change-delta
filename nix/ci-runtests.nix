{ pkgs, scriptDir, go }:
with pkgs;
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
