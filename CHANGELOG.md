# Changelog

## [0.2.1](https://github.com/joshuar/gokbd/compare/v0.2.0...v0.2.1) (2023-06-26)


### Bug Fixes

* **capabilities:** revert to real UID as needed ([04167f8](https://github.com/joshuar/gokbd/commit/04167f8b444c3b7635fa46e7bce562fc0028cc06))

## [0.2.0](https://github.com/joshuar/gokbd/compare/v0.1.1...v0.2.0) (2023-06-25)


### Features

* use setuid/setgid for uinput device creation ([3547eff](https://github.com/joshuar/gokbd/commit/3547effadd1aa45ebcde73c5df9e45fb0fd1e2ea))

## [0.1.1](https://github.com/joshuar/gokbd/compare/v0.1.0...v0.1.1) (2023-05-03)


### Features

* **keyboard:** add a Grab function to grab a virtual keyboard for exclusive use ([0ae0b1c](https://github.com/joshuar/gokbd/commit/0ae0b1cc30bc591513904aad3be2371194fd8efc))
* **keyboard:** extract grab logic into general function for real/virtual keyboards ([8ef01b8](https://github.com/joshuar/gokbd/commit/8ef01b8e284245c7d9a8a103a3f59b3220b9e783))


### Bug Fixes

* **examples:** utilise exposed device node in examples ([5739deb](https://github.com/joshuar/gokbd/commit/5739deb93ff134e6e3280ddda1dff620c0c37b25))


### Miscellaneous Chores

* release 0.1.1 ([65acb5e](https://github.com/joshuar/gokbd/commit/65acb5ea1c3b21137f618c7df574f56fe5da03da))

## 0.1.0 (2023-05-03)


### Features

* **examples:** rewrite snoop example to use a virtual keyboard ([7b53e19](https://github.com/joshuar/gokbd/commit/7b53e19f6bbbf911c1af0cba760e7c0823e4b6a5))
* **keyboard:** add a SnoopKeyboard (single keyboard) function ([5b0a50d](https://github.com/joshuar/gokbd/commit/5b0a50d936a6d6ab4cfb67451883160ec451d5a2))
* **keyboard:** expose errors for Type* functions ([3d9be9a](https://github.com/joshuar/gokbd/commit/3d9be9aa222b81b999149078eeb569b65dca87dc))
* **keyboard:** OpenKeyboardDevice now exposes an error when it fails ([d21756a](https://github.com/joshuar/gokbd/commit/d21756af7536b47c3b2aa7e65f6c4132cddf1e70))
* **keyboard:** reduce sleep time in sending keys to a virtual keyboard ([262b1bf](https://github.com/joshuar/gokbd/commit/262b1bf85fd84871eb113ddfec6aa8ba60afd869))
* **logging:** switch to zerolog for logging ([42ecac3](https://github.com/joshuar/gokbd/commit/42ecac3201d6446d80665ddeafb4c0affe617c15))


### Bug Fixes

* **examples:** add delays to see the results of the example more clearly ([824dce4](https://github.com/joshuar/gokbd/commit/824dce41baf6e1720ed2dfef16748555e9e5aa7c))
* **examples:** fix logging in snoop example ([f7e8f16](https://github.com/joshuar/gokbd/commit/f7e8f166acdf2e2d63df62791bd3cbaa5e815dd8))
* **keyboard:** adjust timing for initialising virtual keyboard ([909ca1a](https://github.com/joshuar/gokbd/commit/909ca1a877cfa5ca18bbecc18dad5e16dc6c9686))
* **keyboar:** remove unused (experimental) function ([f6d0876](https://github.com/joshuar/gokbd/commit/f6d08768c94083035c8b1217cc5ecc097ed562e6))


### Miscellaneous Chores

* release 0.0.1 ([e27838b](https://github.com/joshuar/gokbd/commit/e27838b9279455382293af448300fb89b46794c7))
* release 0.1.0 ([5a0dd43](https://github.com/joshuar/gokbd/commit/5a0dd43e87d41570c37870c0f14e5359d5bd9946))
