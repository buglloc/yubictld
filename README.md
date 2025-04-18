## yubictld

**yubictld** is a lightweight daemon for managing the lifecycle of YubiKeys in test environments.

## Key Features
  - Prevent concurrent operations on the same YubiKey by using locks to ensure safe and consistent access
  - Supports "soft" Yubikeys reboot via FIDO interface using [fidoctl](https://github.com/buglloc/fidoctl)
  - Supports  simulating user "touch" on a YubiKey via [H4ptiX](https://github.com/buglloc/H4ptiX)
  - Automatically detects YubiKeys connected to the same USB hub as the H4ptiX controller, simplifying setup and dynamic environments
