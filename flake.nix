{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "utils";
    };
  };
  outputs = {
    self,
    nixpkgs,
    utils,
    gomod2nix,
  }:
    utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
        
        pkg = pkgs.buildGoApplication {
          pname = "icali-tui";
          version = "0.1.0";
          src = ./.;
          module = ./gomod2nix.toml;
        };
      in {
        packages.default = pkg;
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            gnumake
            gomod2nix.packages.${system}.default
          ];
        };
      }
    );
}
