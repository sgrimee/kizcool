<p align="center">
  <img src="doc/homey.svg" alt="logo" width="64"/>
  <h3 align="center">Kizcool</h3>
  <p align="center">A CLI and go package to control devices from Velux, Somfy and other vendos.<p>
  <p align="center"><a href="https://circleci.com/gh/sgrimee/kizcool"><img src="https://circleci.com/gh/sgrimee/kizcool.svg?style=shield" alt="Build Status"></a></p>
</p>

# Get kizcool

```
go get github.com/sgrimee/kizcool/kizcli
kizcli configure
kizcli get devices
...
```

## Supported gateways

This packages interacts with the Overkiz API, as used by Somfy's [Tahoma](https://shop.somfy.co.uk/tahoma) devices. It should work with other controllers but only Tahoma was tested.

## Supported devices

Support is provided for the following devices:
- Velux Integra electric window: [GGL-GGU](https://roofwindows.veluxshop.co.uk/roof-windows/automated)
- Velux Integra electric roller shutter: [SML](https://www.veluxblindsdirect.co.uk/product/velux-blinds/roller-shutters)
- Velux Integra spotlight: [KRA-100](https://www.amazon.fr/VELUX-integra-fen%C3%AAtre-%C3%A9clairage-kRA-100/dp/B00N33FKGA) (hard to find)

However, the Overkiz system supports many more devices from several vendors. Some may work out of the box. Support for others should be easy to add. Please file an issue to report other working devices or request the addition of new devices.

## Package documentation

A static version of the godoc can be found [here](doc/package.md).

## Roadmap

Features I want to add later on include:
- Listener mode to register for events and see changes triggered via other controllers.
- KNX bridge to control velux devices from a KNX system (the main goal for this project).

## Notice of Non-affiliation and Disclaimer

We are not affiliated, associated, authorized, endorsed by, or in any way officially connected with [Overkiz](https://www.overkiz.com/), [Velux](https://www.velux.com/), [Somfy](https://www.somfy.com/), any other trademark mentioned in this project, or any of its subsidiaries or its affiliates. We are grateful for the great products and services they provide.

Please see the [License and Disclaimer notice](LICENSE).
