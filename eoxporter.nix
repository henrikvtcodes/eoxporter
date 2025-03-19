{
  lib,
  buildGoModule,
  pkgs,
  ...
}:
buildGoModule rec {
  pname = "eoxporter";
  version = "0.0.0";

  src = ./.;
  doCheck = false;
  proxyVendor = true;

  preBuild = ''
    export GOPROXY='https://proxy.golang.org,direct'
    ${pkgs.go}/bin/go mod tidy
    ${pkgs.go}/bin/go mod vendor
  '';

#  vendorHash = "sha256-R7dBP5kR+oRypF2ppCF5CrZiqNzKmQZPxQ3ycQbxjq0=";
vendorHash = null;  
  # vendorHash = lib.fakeHash;
}
