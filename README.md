<p align="center">
  <img src="doc/homey.svg" alt="logo" width="64"/>
  <h3 align="center">Kizcool</h3>
  <p align="center">A CLI and go package to control devices from <a href="https://www.velux.com/">Velux</a>, <a href="https://www.somfy.com/">Somfy</a> and other vendors using the <a href="https://www.overkiz.com/">Overkiz</a> platform.<p>
  <p align="center"><a href="https://circleci.com/gh/sgrimee/kizcool"><img src="https://circleci.com/gh/sgrimee/kizcool.svg?style=shield" alt="Build Status"></a></p>
</p>

# Components

- Executable `kizcool/exe/kizcmd` is a command-line utility to control devices.
- Package `kizcool` provides a high-level client with structs for events and devices.
- Package `kizcool/api` provides a lower-level api client that returns the raw responses from the server without parsing.

# Command line tool

## Download and install kizcmd

```
go get -u github.com/sgrimee/kizcool/cmd/kizcmd
go install github.com/sgrimee/kizcool/cmd/kizcmd
```

## Create a configuration file

```
kizcmd configure
```

## Get all devices names

```
kizcmd get labels
```

## Send a command to a device

```
kizcmd on "my light"
kizcmd off "my light"
kizcmd instensity "my light" 50

kizcmd open "my window"
kizcmd close "my window"

kizcmd open "my blind"
kizcmd close "my blind"
kizcmd closure "my blind" 75
```

## Get all devices data

```
kizcmd get devices
```

## Get one device in json format. Hint: jq is a nice json formatter

```
kizcmd get device "parents window" -ojson | jq
```

## Continuously poll for events from the installation and display them on the console

```
kizcmd listen
```

## Environment variables
As an alternative to the config file, configuration items can be given as environment variables:
- KIZ_USERNAME
- KIZ_PASSWORD
- KIZ_BASE_URL

# Requirements

## Go version

Golang >= 1.13

## Supported gateways

This packages interacts with the Overkiz API, as used by Somfy's [Tahoma](https://shop.somfy.co.uk/tahoma) devices. It should work with other controllers but only Tahoma was tested.

## Supported devices

Support is provided for the following devices:
- Velux Integra electric window: [GGL-GGU](https://roofwindows.veluxshop.co.uk/roof-windows/automated)
- Velux Integra electric roller shutter: [SML](https://www.veluxblindsdirect.co.uk/product/velux-blinds/roller-shutters)
- Velux Integra spotlight: [KRA-100](https://www.amazon.fr/VELUX-integra-fen%C3%AAtre-%C3%A9clairage-kRA-100/dp/B00N33FKGA) (hard to find)

However, the Overkiz system supports many more devices from several vendors. Some may work out of the box. Support for others should be easy to add. Please file an issue to report other working devices or request the addition of new devices.

## Supported platforms

This should work wherever go works. Tested on mac and linux. Efforts were made to support windows but it is not tested.

# Go package documentation

A static version of the godoc can be found [here](doc/package.md).

# Roadmap

Upcoming features include:
- KNX bridge to control velux devices from a KNX system (the main goal for this project).

# Notice of Non-affiliation and Disclaimer

We are not affiliated, associated, authorized, endorsed by, or in any way officially connected with [Overkiz](https://www.overkiz.com/), [Velux](https://www.velux.com/), [Somfy](https://www.somfy.com/), any other trademark mentioned in this project, or any of its subsidiaries or its affiliates. We are grateful for the great products and services they provide.

Please see the [License and Disclaimer notice](LICENSE).
