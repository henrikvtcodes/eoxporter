# eoxporter

A basic Prometheus exporter for Arista EOS devices.

## Usage

```
usage: eoxporter [<flags>]

Flags:
  --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  --eapiConf="~/.eapi.conf"    Path to Arista eAPI config file. If the flag isn't provided, $EAPI_CONF is checked
  --collectors="version,cooling,power,temperature"
                             Comma-separated list of collectors to enable
  --listen="localhost:9396"  Address to listen on for HTTP
```

### Create a view-only user to collect eAPI data

In the interest of security, it is best to create a user for this collector that cannot make changes to the device configuration.
You can use the following commands:

```shell
# Log into your switch (or whatever device)

enable # Enter priveledged mode
configure terminal

# create a readonly role
# this will allow the users with this role to only run the show command and subcommands
role readonly
10 permit mode exec command show
20 deny mode exec command .*
30 deny mode config-all command .*
exit

write # saves the new role

username monitor secret monitor # Creates a user `monitor` with password `monitor`
username monitor role readonly # Add our readonly role to this user
write # copy the running config to startup

exit
exit

# Logout
```

See the [`eapi.ini`](#prometheus--eapi-example-config) file in this repository to get an idea how this is configured.

### Requirements for running on your device

**Prerequisites**

- an eAPI configuration file ([example](https://github.com/aristanetworks/goeapi?tab=readme-ov-file#example-eapiconf-file))
- The ability to build a Go app (for my Nix people there's a flake that provides a package here)

**Collector Configuration**
All configuration can be supplied using CLI flags. None are required unless otherwise specified

| Flag           | Format                    | Default                             | notes                                                                               |
|----------------|---------------------------|-------------------------------------|-------------------------------------------------------------------------------------|
| `--collectors` | comma separated list      | `version,cooling,power,temperature` | the metrics that should be provided if the request lacks a `collectors` query param |
| `-eapiConf`    | Relative or Absolute path | if no flag, `$EAPI_CONF` is checked | Required in flag or env form                                                        |
| `--listen`     | `address:port` format     | `localhost:9396`                    |                                                                                     |

**Valid Collectors**

- `version`: provides metadata about the switch, and the following
- - Switch memory usage (and total installed memory)
- - Uptime
- `cooling`: status for all fans in the system
- - current speed
- - set speed
- - maximum speed
- `temperature`: temperature readings in degrees celsius from all sensors in the system. for each sensor, provides:
- - current temperature
- - critical temperature (useful for alerting)
- - overheat temperature (useful for alerting)
- - target temperature
- `power`: power consumption info
- - PSU metadata (model, capacity)
- - power consumption in watts
- - input AC current
- - output DC current
- - uptime
- `interfaces`: traffic in/out per interface
- - in/out multicast/unicast/broadcast traffic
- - in/out octets
- - in/out octets

### Prometheus & eAPI Example Config

```ini
[connection:<name of configured device>]
host=<device management host/ip>
username=monitor
password=monitor
transport=http
```

This config allows you to provide the names of your Arista EOS devices in the target section of the config (as though they were exporting the metrics themselves)
but moves the labels around so that it works properly. The

```yaml
scrape_configs:
  - job_name: Arista EOS Status
    static_configs:
      - labels:
          collectors: version,power,temperature,cooling,interfaces
        targets:
          - <name of configured device>
    relabel_configs:
      - source_labels:
          - __address__
        target_label: __param_target
      - source_labels:
          - __param_target
        target_label: instance
      - source_labels:
          - collectors
        target_label: __param_collectors
      - replacement: <host:port of your eoxporter instance> # REPLACE THIS
        target_label: __address__
    scrape_interval: 15s
```

## Roadmap

### Support more metrics collectors

_Open to contributions only if the contributor can test on an actual arista device._

I currently run an Arista DCS-7050TX-64 and that is what I am using to test. In the future I would like to add collectors for MLAG and other things, but I want to
actually set that stuff up on my switch so that I can test my metrics implementation with actual data from a real device.

### More robust logging and CLI parsing

I am currently using the Go std library for flags and logging right now but they aren't great.

### Utilize switch extensions to run on Arista devices directly

I think this is theoretically possible and it would be cool to collect metrics directly from the device but there's some experimentation to be done with how eAPI would work there (sockets?)

## Developing & Building

### Building the nix package with `gomod2nix`
The tool `gomod2nix` is used because there are some weird issues with the Nixpkgs standard way of building Go modules; 
particularly centered around hashing of dependencies and whatnot. `gomod2nix` consumes the `go.mod` and `go.sum` files to 
create a file containing the versions and hashing of Go dependencies in a format that Nix can understand.

**TL;DR: Whenever dependencies change in any way, the `gomod2nix.toml` file must be regenerated.**

The default dev shell defined by this flake contains the tool to do this; it is simply the command `gomod2nix`.

## Credits

- https://github.com/aristanetworks/goeapi
- https://github.com/ubccr/arista_exporter (inspiration)
