{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
  let
    forAllSystems = nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed;
    nixpkgsFor = forAllSystems(system: import nixpkgs {
      inherit system;
      overlays = [ self.overlays.default ];
    });
  in
  {
    overlays.default = final: prev: { };

    devShells = forAllSystems(system:
    let
      pkgs = nixpkgsFor.${system};
    in
    {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [
          nixpkgs-fmt
          go_1_21
          nodejs_20
          postgresql_15
          gotools # The official go tools (https://go.googlesource.com/tools)
          go-tools # silly nix (tooling from https://staticcheck.dev/)
          gopls
        ];
        shellHook = ''
          mkdir -p .data
          export LANG=C.UTF-8
        '';
      };
    });
  };
}

# vim: set expandtab tabstop=2 shiftwidth=2:
