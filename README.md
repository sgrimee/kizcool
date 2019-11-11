![](doc/homey.svg | width=100)

# kizcool
[![CircleCI](https://circleci.com/gh/sgrimee/kizcool.svg?style=svg)](https://circleci.com/gh/sgrimee/kizcool)

**A go package and CLI tool to control IO Velux devices.**

## Supported gateways

This packages interacts with the Overkiz API, as used by Somfy's [Tahoma](https://shop.somfy.co.uk/tahoma) devices. It should work with other controllers but only Tahoma was tested.

## Supported devices

Support is provided for the following devices:
- Velux Integra electric window: [GGL-GGU](https://roofwindows.veluxshop.co.uk/roof-windows/automated)
- Velux Integra electric roller shutter: [SML](https://www.veluxblindsdirect.co.uk/product/velux-blinds/roller-shutters)
- Velux Integra spotlight: [KRA-100](https://www.amazon.fr/VELUX-integra-fen%C3%AAtre-%C3%A9clairage-kRA-100/dp/B00N33FKGA) (hard to find)

However, the Overkiz system supports many more devices from several vendors. Some may work out of the box. Support for others should be easy to add. Please file an issue to report other working devices or request the addition of new devices.

## Command-line application

An executable app called `kizcli` allows control of the devices from the command line.

```
go get github.com/sgrimee/kizcool/kizcli
kizcli configure
kizcli get devices
...
```

## Package documentation

A static version of the godoc can be found [here](doc/package.md).

## Roadmap

Features I want to add later on include:
- Listener mode to register for events and see changes triggered via other controllers.
- KNX bridge to control velux devices from a KNX system (the main goal for this project).

## Notice of Non-affiliation and Disclaimer

We are not affiliated, associated, authorized, endorsed by, or in any way officially connected with [Overkiz](https://www.overkiz.com/), [Velux](https://www.velux.com/), [Somfy](https://www.somfy.com/), any other trademark mentioned in this project, or any of its subsidiaries or its affiliates. We are grateful for the great products and services they provide.

Please see the [License and Disclaimer notice](LICENSE).