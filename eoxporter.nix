{lib, buildGoModule, ...} :

buildGoModule rec {
  src = ./.;


  vendorHash = lib.fakeHash;
}