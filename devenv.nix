{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  languages.go.enable = true;
  pre-commit.hooks = {
    shellcheck.enable = true;
    alejandra.enable = true;
    shfmt.enable = true;
    gofmt.enable = true;
    staticcheck.enable = true;
  };
}
