# Media Backup for Childcare Services

[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT)
![Go version for branch](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup/main)
![CI Status](https://img.shields.io/github/actions/workflow/status/leocov-dev/tadpoles-backup/ci.yml)

## About
This tool will allow you to save all your child's images and videos at full resolution from various service providers. Comments and timestamp info will be applied as EXIF image metadata where possible.

Providers:
* Tadpoles

---
## Install
Download the exe from above and run

---
## Usage

```
# Print help with command details:
$ tadpoles-backup --help

# Get account statistics
$ tadpoles-backup --provider <service-provider> stat

# Download media (only new files not present in the target dir are downloaded)
$ tadpoles-backup --provider <service-provider> backup <a-local-directory>

# Clear Saved Login
$ tadpoles-backup --provider <service-provider> clear login
```

### Provider Notes

#### Tadpoles

You **MUST** have a _www.tadpoles.com_ account with a tadpoles specific password.
You **CAN NOT** log in to this tool with Google Auth.
If you normally log into _tadpoles.com_ with Google/Gmail account verification you will need to
request a password reset with the command:
```shell
# this simply requests a reset email be sent to you
# it does not change or access your password
$ tadpoles-backup --provider tadpoles reset-password <email>
```

The tool stores your _www.tadpoles.com_ authentication cookie for future use so that you don't need to enter your password every time.
This cookie lasts for about 2 weeks. Your email and password are never stored.


## Development

See the contributing guide [here](CONTRIBUTING.md).

### Basic Setup

Install the Go version defined in [go.mod](go.mod) or use [goenv](https://github.com/syndbg/goenv) to manage Go (as set by [.go-version](.go-version)).

### Dev build
```shell
# build for your platform only and run.
$ make && bin/tadpoles-backup --help
```

### Testing

Run all unit tests with helper utility. This will build a coverage report as
`coverage.html`
```shell
make test
```


---
## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [ChuckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)

## Thanks to
* @arthurnn - for assistance with Docker image
* @AndyRPH - for assistance with Bright Horizons support
* @s0rcy - for assistance with Tadpoles password reset
