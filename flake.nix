{
  description = "A Nix-flake-based Go 1.23 development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";

  outputs = { self, nixpkgs }:
    let
      goVersion = 23;

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            gnumake
            gh

            go
            gotools
            air
            sqlc
            go-migrate
            sqlite
            litecli

            nodejs_20
            pnpm_8
          ];

          shellHook = ''
            echo 'Installing sqlc completion'
            source <(sqlc completion bash)
          '';
        };
      });
    };
}
