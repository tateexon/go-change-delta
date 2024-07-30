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
    # basics
    bash
    git
    curl
    gnumake
    jq
    dasel
    github-cli

    # go
    go
    go-mockery
    gotools
    gopls
    delve
    golangci-lint

    # linting tools
    typos
    pre-commit
    python
    shfmt
    shellcheck
  ];

  CGO_ENABLED = "0";

  shellHook = ''
    # Uninstall pre-commit hooks in case they get messed up
    pre-commit uninstall > /dev/null || true
    pre-commit uninstall --hook-type pre-push > /dev/null || true

    # enable pre-commit hooks
    pre-commit install > /dev/null
    pre-commit install -f --hook-type pre-push > /dev/null

    # install gotestloghelper
    go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest
  '';
}
