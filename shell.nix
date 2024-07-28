{ pkgs, scriptDir }:
with pkgs;
let
  go = pkgs.go_1_22;

  mkShell' = mkShell.override {
    # The current nix default sdk for macOS fails to compile go projects, so we use a newer one for now.
    stdenv = if stdenv.isDarwin then overrideSDK stdenv "11.0" else stdenv;
  };
in
mkShell' {
  nativeBuildInputs = [
    git
    go
    curl
    go-mockery
    gotools
    gopls
    delve
    golangci-lint
    github-cli
    jq
    dasel
    typos
  ];

  GOROOT = "${go}/share/go";
  CGO_ENABLED = "0";

  shellHook = ''
    # enable pre-commit hooks
    # pre-commit install > /dev/null
    # enable pre-push hooks
    # pre-commit install -f --hook-type pre-push > /dev/null
    # setup go bin for nix
    export GOBIN=$HOME/.nix-go/bin
    mkdir -p $GOBIN
    export PATH=$GOBIN:$PATH
    # install gotestloghelper
    go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest
  '';
}
