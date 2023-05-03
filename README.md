<!--
 Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# gokbd

![GitHub](https://img.shields.io/github/license/joshuar/gokbd)
[![Go Reference](https://pkg.go.dev/badge/github.com/joshuar/gokbd.svg)](https://pkg.go.dev/github.com/joshuar/gokbd)
[![gokbd](https://goreportcard.com/badge/github.com/joshuar/gokbd?style=flat-square)](https://goreportcard.com/report/github.com/joshuar/gokbd)
[![codecov](https://codecov.io/gh/joshuar/gokbd/branch/main/graph/badge.svg?token=2BDVOTORZB)](https://codecov.io/gh/joshuar/gokbd)

## About

gokbd is a package that uses
[libevdev](https://www.freedesktop.org/wiki/Software/libevdev/) to talk to a
keyboard on Linux. It allows snooping the keys pressed as well as typing out
keys.

## Usage

```go
import gokbd "github.com/joshuar/gokbd"
```

Examples for reading what keys are being typed (snooping) and writing to a
virtual keyboard are available under the `examples/` directory. To run them:

```shell
go run examples/snoop/main.go # snooping
go run examples/type/main.go # typing
```

## Permissions

You may need to grant additional permissions to the user running any program
using `gokbd`.

- To read (snoop) from keyboards, the user will need to be part of the `input`
  group. Typically, the user can be added with the following command:

```shell
sudo gpasswd -a $USER input
```

- To create a virtual keyboard and write to it, the user will need access to the
  [kernel uinput
  device](https://kernel.org/doc/html/latest/input/uinput.html). Typically, this
  can be granted with a [udev rule](https://opensource.com/article/18/11/udev)
  like the following:

```shell
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
```
