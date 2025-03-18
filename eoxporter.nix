{
  lib,
  buildGoModule,
  ...
}:
buildGoModule rec {
  pname = "eoxporter";
  version = "0.0.0";

  src = ./.;

  vendorHash = "sha256-R7dBP5kR+oRypF2ppCF5CrZiqNzKmQZPxQ3ycQbxjq0=";
  # vendorHash = lib.fakeHash;
}
