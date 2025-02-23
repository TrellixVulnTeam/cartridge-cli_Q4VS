# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Updated `cartridge` to `2.7.6` and `metrics` to `0.15.1` in application template.
- Now file times are preserved when copying.

## [2.12.3] - 2022-11-09

### Changed

- A unix domain socket pathname (used in cartridge-cli to connect to instances) longer
  than the value defined by the system resulted in a connection error.
  It now works fine or throws a human readable error if we can't work around this
  limitation by connecting to an instance with a relative socket path.

### Added

- Tarantool benchmark tool update (select and update operations):
  * option --insert has been added - sets percentage of insert operations to bench space.
  * option --select has been added - sets percentage of select operations from bench space.
  * option --update has been added - sets percentage of update operations in bench space.
  * option --fill" has been added - sets number of records to pre-fill the space.
- An ability to choose Tarantool version to install in the result docker image:
  * command line option --tarantool-version
  * tarantool.txt config file with TARANTOOL option

### Fixed

- Fixed sporadic "Failed to dial" errors on cartridge cli different commands. Only
  local instances will be used for commands execution.

## [2.12.2] - 2022-06-05

### Changed

- Updated `cartridge` to `2.7.4` and `metrics` to `0.13.0` in application template.

## [2.12.1] - 2022-05-26

### Changed

- Loosed ``cartridge pack`` `--version` and `--suffix` verification,
  so it will log a warning instead of returning an error if non-valid string is passed
  (in terms or RPM/DEB standards).
### Added

- Ability to explicitly set a full name of the bundle created by ``cartridge pack``.
  (flag `--filename`)

### Fixed

- Bug that caused a typed command not to be displayed on the terminal after
  exiting the "connect" console.

## [2.12.0] - 2022-03-30

### Changed

- **Breaking**: change ``cartridge pack`` naming strategy:
  * RPM: `<app-name>-<version>[.<suffix>]-1.<arch>.rpm`,
  * DEB: `<app-name>_<version>[.<suffix>]-1_<arch>.deb`,
  * TGZ: `<app-name>-<version>[.<suffix>].<arch>.tar.gz`.
  `<version>` is `git describe --tags --long` output transformed to `X.Y.Z.N` or non-transformed ``--version`` value.
  (Previously it was transformed to `X.Y.Z-N-gHASH` with ability to set `X.Y.Z-N`, `X.Y.Z-gHASH` with `--version`.)
- **Breaking**: ``cartridge pack --version`` value is used as provided without transformation.
  (Previously it had the same restrictions and transformations as `git describe --tags --long`
  output.)
- **Breaking**: ``cartridge pack`` ``--version`` and ``--suffix`` is validated
  after DEB/RPM conventions for corresponding type (see debian policy [
  [1](https://www.debian.org/doc/debian-policy/ch-controlfields.html#version),
  [2](https://www.debian.org/doc/manuals/debian-reference/ch02.en.html#_debian_package_file_names)],
  rpm policy [[1](http://ftp.rpm.org/max-rpm/ch-rpm-file-format.html)]).
  For example, dashes in RPM version (like `1.2.3-0`) is no longer supported.
- Bump `Go` requirement to `1.18`.
- Replace vfsgen with embed for ``cartridge create`` template.
- Updated ``containerd`` version to 1.5.10 to fix the vulnerability bugs:
  bug was found in containerd where containers launched through containerd’s
  CRI implementation with a specially-crafted image configuration could gain
  access to read-only copies of arbitrary files and directories on the host.
  This may bypass any policy-based enforcement on container setup (including
  a Kubernetes Pod Security Policy) and expose potentially sensitive
  information. Kubernetes and crictl can both be configured to use containerd’s
  CRI implementation.
  CVE ID: CVE-2022-23648
  GHSA ID: GHSA-crp2-qrr5-8pq7
- Updated ``golang.org/x/text`` version to 0.3.7 to fix the vulnerability bug:
  due to improper index calculation, an incorrectly formatted
  language tag can cause Parse to panic, due to an out of bounds read.
  If Parse is used to process untrusted user inputs, this may be used
  as a vector for a denial of service attack.
  CVE ID: CVE-2021-38561
- Updated ``golang.org/x/crypto`` version to 0.0.0-20211202192323-5770296d904e
  to fix the vulnerability bug:
  there's an input validation flaw in ``golang.org/x/crypto``
  ``readCipherPacket()`` function. An unauthenticated attacker who sends
  an empty plaintext packet to a program linked with ``golang.org/x/crypto/ssh``
  could cause a panic, potentially leading to denial of service.
  CVE ID: CVE-2021-43565

## [2.11.0] - 2022-01-26

### Changed
- Updated ``containerd`` version to 1.5.8 to fix the vulnerability bugs:
  * https://github.com/advisories/GHSA-c2h3-6mxw-7mvq (container root
    directories and some plugins had insufficiently restricted permissions,
    allowing otherwise unprivileged Linux users to traverse directory contents
    and execute programs.)
  * https://github.com/opencontainers/image-spec/security/advisories/GHSA-77vh-xpmg-72qh
    https://github.com/opencontainers/distribution-spec/security/advisories/GHSA-mc8v-mgrf-8f4m
    (manifest and index documents are ambiguous without an accompanying Content-Type
    HTTP header. Versions of containerd prior to 1.4.12 and 1.5.8 treat the Content-Type
    header as trusted and deserialize the document according to that header.
    If the Content-Type header changed between pulls of the same ambiguous
    document (with the same digest), the document may be interpreted differently,
    meaning that the digest alone is insufficient to unambiguously identify
    the content of the image.)
- Updated ``image-spec`` version to 1.0.2 to fix the vulnerability bug
  https://github.com/opencontainers/distribution-spec/security/advisories/GHSA-mc8v-mgrf-8f4m
  (in the OCI Image Specification version 1.0.1 and prior, manifest and index
  documents are not self-describing and documents with a single digest could be
  interpreted as either a manifest or an index.)
- Updated `cartridge` to `2.7.3` in application template.
  Incompatibility with `tarantool` `2.10.0~beta2` or greater is fixed in `2.7.3`,
  see https://github.com/tarantool/cartridge/commit/1649b40743bd05f2e71f2ce1ad5b5caf19f864d4
- Updated `metrics` to `0.12.0` in application template.
- Updated `luatest` to `0.5.6` in application template.
- Updated `luacheck` to `0.26.0` in application template.
- ``cartridge pack`` command now ignores the contents of ``run_dir``, ``data_dir``,
  ``log_dir`` directories (if default or set with ``.cartridge.yml``).
- Change package release policy:
  * drop EL6 (CentOS 6.x, RHEL 6.x, CloudLinux 6.x);
  * drop Fedora 29;
  * drop Ubuntu 14.04 (Trusty);
  * drop Debian 8 (Jessie);
  * add Fedora 31, 32, 33, 34, 35;
  * add Ubuntu 21.04 (Hirsute).

### Fixed

- Cartridge errors in the ``replicasets`` command are now more readable.
- Removed unnecessary flags (``--rocks``, ``--project-path``) from ``cartridge help`` command.
- Fixed project build with capital letters in the project name.
- Fixed display of Docker image pull (``cartridge pack`` command with ``--verbose`` flag).
- Added support for new Tarantool release policy.
- ``mage clean`` removes all generated code.
- All mage test commands now depend on code generation step.

### Added

- Tarantool benchmark tool (early alpha, API can be changed in the near future).
- Ability to reverse search in ``cartridge enter`` and ``cartridge connect`` commands.
- Added support for functionality from Golang 1.17.
- Describe scenario of local test run.

## [2.10.0] - 2021-07-28

### Added

- Ability to specify pre and post install scripts for the RPM and DEB (command
  ``cartridge pack``), using the ``--preinst`` and ``--postinst`` flags.
- ``cartridge pack`` generates ``VERSION.lua`` file with the current
  version of project.
- Ability to caching any paths specified in ``pack-cache-config.yml`` file
  when packaging application via ``cartridge pack`` command.
- Ability to specify fd limit in the systemd unit template
  (command ``cartridge pack``) in the ``systemd-unit-params.yml`` file.
- ``cartridge pack`` now uses the VERSION file from the
  ``TARANTOOL_SDK_PATH`` environment variable on building in Docker
- Ability to specify in the `systemd-unit-params.yml` file arguments passing
  by env with systemd unit file.
- ``cartridge failover`` command to manage failover

### Fixed

- Improved tests in ``cartridge create`` template:
  * Tests are reduced to the form corresponding to luatest documentation
  * ``before_suit`` now remove ``.xlog`` and ``.snap`` files
- Now ``stateboard: true`` specified in ``.cartridge.yml`` affects
  only ``cartridge start/stop/status/log/clean`` calls without
  arguments, e.g. ``cartridge stop router`` doesn't lead to
  stopping stateboard too, but ``cartridge stop`` stops all
  instances includes stateboard.
- Fixed incorrect error message when trying to
  ``cartridge replicasets bootstrap-vshard`` without a configured cluster

### Changed

- Removed setting `cluster_cookie` on `cartridge.cfg` in application template
- Updated `metrics` to `0.9.0` in application template
- ``cartridge version`` and ``cartridge --version`` commands
  now show Cartridge CLI and Cartridge versions. Moreover, they
  can show version of the project rocks.

## [2.9.1] - 2021-05-13

### Fixed

- Allowed to pack RPM and DEB packages with ``cartridge pack``
  command without default dependencies file.

## [2.9.0] - 2021-04-26

### Changed

- Updated `cartridge` to `2.6.0` in application template
- Updated `metrics` to `0.8.0` in application template

### Added

- Ability to specify dependencies for the RPM and DEB (command
  ``cartridge pack``), using the ``--deps`` and ``--deps-file`` flags.

### Fixed

- Improved ``cartridge create`` template:
  * Removed extra http-endpoint ``metrics`` in ``app.roles.custom.lua.``
  * Removed mix of spaces and tabs in ``.rockspec`` file.

## [2.8.0] - 2021-04-07

### Added

- `--spec` option for `build` and `pack` commands to specify a path to rockspec for current build.

### Fixed

- It is possible to run an image generated with the ``cartridge pack docker``
  command in an unprivileged Kubernetes container. It became possible, because
  tarantool user now always has ``UID = 1200`` and ``GID = 1200``.
- Correct display of insertion of multi-line code snippets
  in ``cartridge enter`` command.
- Improved responsiveness of the ``cartridge enter`` and ``cartridge enter``
  commands. Requests that work with a large amount of data have become faster.
- Make pack type case insensitive.

## [2.7.2] - 2021-03-24

### Changed

- Updated `cartridge` to `2.5.1` in application template
- Updated `metrics` to `0.7.1` in application template
- Updated `cartridge-cli-extensions` to `1.1.1` in application template

### Added

- Variables ``TARANTOOL_WORKDIR``, ``TARANTOOL_PID_FILE`` and
  ``TARANTOOL_CONSOLE_SOCK`` can be customized when packing in docker via
  ``cartridge pack`` command. Variables ``CARTRIDGE_RUN_DIR`` and
  ``CARTRIDGE_DATA_DIR`` have also been added.

## [2.7.1] - 2021-03-18

### Changed

- Updated `cartridge` to `2.5.0` in application template

### Fixed

- Added interruption of an incomplete expression when pressing
  ``Ctrl-C`` in ``cartridge enter`` command.
- Packing applications that use `cartridge 2.5.0` and higher

## [2.7.0] - 2021-03-11

### Fixed

- Connector crashing on using `cartridge admin` with binary port
- ``cartridge pack docker`` consumes all disk space

### Added

- `--no-log-prefix` option for `cartridge start` command to disable instance name prefix in logs when running interactively.
- default metrics and health check endpoint in template
- added ability to specify stateboard flag in .cartridge.yml

### Changed

- Update `metrics` to `0.7.0`
- Reworked ``cartridge create`` template:
  * Registration ``admin`` function and path configuration moved to separated files
  * Update and improve ``helper`` in tests
  * Changed standard port of the ``stateboard``
  * Added README.md
  * And other cosmetic fixes
- Allowed to use any base docker image for the ``cartridge pack`` command.
- Executable with any name (not only ``tarantool``) can run processes.
- Updated Go to version 1.16

## [2.6.0] - 2021-01-27

### Added

- `--conn, -c` option for `cartridge admin` command to specify instance address
- `-l` shortcut for `--list` flag of `cartridge admin` command

### Fixed

- Parsing cluster-wide config with empty roles list
- Parsing 'number' args on calling admin functions

### Changed

- Update `metrics` to `0.6.1`

## [2.5.0] - 2020-12-29

### Fixed

- Using Tarantool console socket
  * end of Tarantool output data and read timeout are handled properly
  * Tarantool greeting is read once on connection creation
- Logs writer used on interactive start: it become waiting forever on big output
  received (such as curl verbose log)

### Changed

- Improved error message on building in docker fail on GitLab CI
- `cartridge pack` fails for RPM and DEB if `--use-docker` isn't specified
- Refactored verbosity flags:
  * `--quiet`: no logs (only errors are shown)
  * no flags: logs + spinner instead of commands/docker output
  * `--verbose`: logs + commands/docker output
- Spinner is started only for a terminal
- Update `cartridge` to `2.4.0` (and `checks` to `3.1.0`)
- Update `metrics` to `0.6.0`

### Added

- `cartridge replicasets` command to manage replicasets on local running
- `cartridge enter` command to connect to local running instance
- `cartridge connect` command to connect to instance by address
- Messages from `print` in admin functions are displayed on `cartridge admin` call

## [2.4.0] - 2020-10-26

### Fixed

- Bash completion file mode discarding

### Added

- `cartridge admin` command to call admin functions provided by application

### Changed

- Updated tarantool/metrics to version 0.5.0

## [2.3.0] - 2020-08-31

### Added

- `cartridge repair` command to patch and reload cluster configuration files

### Changed

- Updated Go to version 1.15
- Updated `cartridge` to 2.3.0 in default application template

## [2.2.1] - 2020-08-19

### Fixed

- Now instance process is ran from the application directory

### Added

- `cartridge gen completion` command to generate shell completion scripts
- Bash completion is delivered with RPM and DEB packages

## [2.2.0] - 2020-08-12

### Added

- `--force` option for `stop` command to send SIGKILL to instances
- `cartridge clean` command to remove instance(s) files
- `--from` option for `create` command to use custom application templates

### Changed

- Update `cartridge` to 2.2.1

## [2.1.0] - 2020-07-17

### Added

- `cartridge log` command to get logs of instances running in background
- `--timeout` option for `start` command
- `--version, -v` command to print `cartridge-cli` version

## [2.0.1] - 2020-07-07

### Fixed

- Error on packing application without build hooks
- Unexpected end of JSON input error on image build

## [2.0.0] - 2020-07-02

### Changed
  
- Completely rewritten in Go
- `CARTRIDGE_BUILDDIR` is renamed to `CARTRIDGE_TEMPDIR`;
  now it can be project subdirectory
- `centos:7` is allowed to be used as a base image al well as `centos:8`
- `--tag` option for `pack` command is an array of strings now
- `start`, `stop`, `status` commands requires only instance names,
  application name is discovered from a rockspec or passed by `--name` option
- `cartridge.{pre,post}-build` hooks should be executable
- `cartridge-cli` can't be installed as a rock module
- build requires rockspec in the application directory


### Added

- `--quiet` flag to hide build output
- `--verbose` flag to make output more verbose
- `--data-dir` and `--log-dir` options for `start` command
- `--cache-from` and `--co-cache` options for `pack` command on building in docker
- `--stateboard-unit-template` option for `pack` command

### Removed

- `TARANTOOL_DOCKER_BUILD_ARGS` env variable for `pack` command
- deprecated `.cartridge.pre` and `.cartridge.ignore` support

## [1.8.3] - 2020-06-29

### Added

- New metrics role in template project

### Changed

- Update `cartridge` to 2.2.0

## [1.8.2] - 2020-05-28

### Fixed

- Project version isn't detected by git when `--tag` option is specified
- Fixed a bug in normalizing the version passed with `--version`.
  The expected pattern is `major.minor.patch[-count][-commit]`,
  but previously the normalization failed: `--version xxx1.2.3xxx`
  resolved to `1.2.3` instead an error.

## [1.8.1] - 2020-05-06

### Fixed

- Fixed docker image fullname

## [1.8.0] - 2020-04-27

### Added

- Cartridge Stateboard support:
  * Application template contains stateboard entrypoint script and configuration
  * Unit file for stateboard `systemd` service is delivered in RPM/DEB
  * Added `--stateboard` and `--stateboard-only` options for `start` and `stop`
    commands to start/stop stateboard locally
- Warning on running `cartridge start` without `cartridge build` before
- Checking notify socket length on `cartridge start -d`
- `cartridge status` command to check instances status

### Changed

- Prettified `start` and `stop` logs
- `start` and `stop` commands try to start/stop all instances and accumulate
  errors
- If instance is already stopped, `stop` command doesn't fail, only warning
  message is shown
- Update `cartridge` to 2.1.2

### Fixed

- `Not enough memory` error on running `cartridge pack`
- Error on project files and build directory containing symbols that
  should be escaped

## [1.7.0] - 2020-04-10

### Added

- Option `--suffix` for `pack` command to specify the result artifact name suffix.

## [1.6.0] - 2020-04-03

### Added

- Packing in docker. Added a new option `--use-docker` for `cartridge pack` command.
  This option allows to build application in docker image.

## [1.5.0] - 2020-03-27

### Changed

- Git errors aren't fatal, if `git clean` command fails (in the project root or
  for sumbodules), it just prints warning message
- Project build: one Lua-executable is compiled

## [1.4.2] - 2020-03-17

### Added

- Option `--build-from` to specify build image base layers.
- Add `TARANTOOL_DIR` to rockspec build.variables.

### Changed

- Refactored packing to docker: `--download-token` option is replaced with
  `--sdk-local` and `--sdk-path` options.
- Refactored RPM and DEB scripts (pre- and post- install, ExecStartPre in systemd
units) to be platform independent.

## [1.4.1] - 2020-03-06

### Changed

- Improved arguments parsing:
  * boolean flags `--flag` shouldn't be passed after all other options;
  * Both `--long_opt` and `--long-opt` patterns can be used, it will be parsed
    as `long_opt` option

### Fixed

- Docker error on placing dockerfile not within the build context
- Creating files owned by root on local machine when building application in docker

## [1.4.0] - 2020-02-05

### Added

- Allow to pass directory to build application in using `CARTRIDGE_BUILDDIR`
  environment variable
- `cartridge build` command to build application locally

### Changed

- By default, temporary directory for application building is created in
  `~/.cartridge/tmp`
- Commands usage messages are prettified
- `path` argument for `cartridge pack` command isn't required.
  By default, current directory is used.

### Fixed

- Delayed messages on application packing

## [1.3.2] - 2020-01-23

### Changed

- Common packing flow parameters are stored in the global `pack_state`
- Update cartridge to 2.0.1 in template
- Update luatest to 0.5.0
- Add luacov to app template

### Fixed

- Error on runnning `git clean` for submodules on `cartridge pack`

## [1.3.1] - 2020-01-20

### Added

- Allow to pass `--version` in format `major.minor.patch[-count][-hash]`

### Changed

- RPM header: added `PAYLOADDIGEST` and `PAYLOADDIGESTALGO` flags,
  removed `RPMVERSION`.

## [1.3.0] - 2020-01-13

### Added

- Packing to Docker image
- Check filemodes before packing
- `--from` option for `docker pack` command to specify base image Dockerfile path
- `cartridge.pre-build` and `cartridge.post-build` hooks
  to be ran before and after `rocks make`
- Deprecated build flow (`.cartridge.ignore` + `.cartridge.pre`) is supported
  for all distribution types except `docker`
- Recursively cleaning all submodules on application packing

### Changed

- `docker pack` log messages are coloured
- Pre-build, build and post-build actions are grouped in one RUN directive
  on packing to Docker image
- Update luatest to 0.4.0
- Freeze cartridge 2.0.0 in template

### Fixed

- Error on using environment variables in base Dockerfile
- Error on using COPY instruction in base Dockerfile
- Added missing environment variable `TARANTOOL_APP_NAME`
- Fix parsing options priority to match `cartridge.argparse`. Current parsing priority:
  firstly commandline options, then environment variables, then options from
  `.config.yml`, in the end default options. Options from `.config.yml` are
  overriden by corresponding to them environment variables.
- Error on rocks manifest processing

## [1.2.1] - 2019-11-25

- Fix building RPM package on CentOS 8
- Fix starting foreground apps with current Tarantool

## [1.2.0] - 2019-11-15

### Added

- luacheck in examples and templates
- `--version` option to display version
- Default cartridge-cli configuration in getting-started template
- Use current tarantool executable to start instance

### Changed

- Warnings in log are shown with yellow color
- `cartridge start` starts instances in foreground, `--foreground` is replaced with `--daemonize`

### Removed

- `plain` template

## [1.1.0] - 2019-10-24

### Added
- Start and stop all instances
- Start/stop instances defined in multiple files
- Colorized logs and prefix with instance name for multiple foreground instances
- Packing DEB

### Changed
- Disabled jit in tests until tarantool/tarantool#4476 is fixed
- Getting started app READMEs improved

### Fixed
- Luacheck warnings
- Missing setsearchroot in 1.10
- /var/run dir removal after reboot

## [1.0.0] - 2019-09-02

### Added

- Basic functionality
- Integration tests
- End-to-end tests
- Gitlab CI tests for opensource and enterprise Tarantool
- Packing RPM with Tarantool dependency for opensource
- Loading templates from .rocks
- Configuring systemd units using `cartrigde.argparse` way
- Getting started app
- Start and stop commands
- Cache downloaded sdk between ci jobs
