# eoxporter

A basic Prometheus exporter for Arista devices. 

## Usage

### Create a view-only user to collect eAPI data
In the interest of security, it is best to create a user for this collector that cannot make changes to the device configuration.
You can use the following commands:
```shell
# Log into your switch (or whatever device)

enable # Enter priveledged mode
configure terminal 

username monitor secret monitor # Creates a user `monitor` whos password is `monitor`
write # copy the running config to startup

exit
exit

# Logout
```
See the `eapi.ini` file in this repository to get an idea how this is configured.

## Developing

### To update the vendorHash in the nix package (`eoxporter.nix`)

_This must be done when dependencies change._
Set `vendorHash = lib.fakeHash;` and build the package. This will trigger an error which will provide the correct hash, which can then be set.

## Credits

- https://github.com/aristanetworks/goeapi
- https://github.com/ubccr/arista_exporter (design pattern inspiration)