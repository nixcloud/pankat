{ pkgs ? import <nixpkgs>{}, ... }:

pkgs.buildGoPackage rec {
  name = "nixcloud-backend-${version}";
  version = "0.0.1";

  #src = ./src/nixcloud/backend;
  goPackagePath = "github.com/nixcloud/websocket";
  buildInputs = [ pkgs.pandoc ];

  # HACK: check that ./src/nixcloud/backend/backend does not exist or it will fail
  #preConfigure = ''
  #  rm -f backend
  #'';

  #goPackagePath = "nixcloud/backend";
  goDeps = ./deps.nix;

  #extraSrcPaths = [ "${leapsLibSrc}" ];
  #buildInputs = if pkgs?compile-daemon then [ pkgs.compile-daemon ] else [ ];

}
