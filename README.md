# eoxporter

A basic Prometheus exporter for Arista devices. 


## Developing

### To update the vendorHash in the nix package (`eoxporter.nix`)

_This must be done when dependencies change._
Set `vendorHash = lib.fakeHash;` and build the package. This will trigger an error which will provide the correct hash, which can then be set.

## Credits

- https://github.com/aristanetworks/goeapi
- https://github.com/ubccr/arista_exporter (design pattern inspiration)