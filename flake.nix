{
  description = "Go Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs
    , flake-utils
    , ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        nativeBuildInputs = with pkgs; [
          delve
          go
          golangci-lint
          golangci-lint-langserver
          gopls
          gotools
        ];
      in
      {
        devShells.default = pkgs.mkShell { inherit nativeBuildInputs; };
      }
    );
}
