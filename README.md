# docker2exe

This tool can be used to convert a Docker image to an executable that you can send to your friends!

## Installation

Minimal dependencies are:

- A working go environment with version >=1.18
- [pigz](https://zlib.net/pigz/)

To compile, just run `make` and executables for each platform will be generated in `dist/` folder.

## Usage

To create a new binary:

    $ docker2exe --name alpine --image alpine:3.9

This will create the following files:

    dist
    ├── alpine-darwin-amd64
    ├── alpine-linux-amd64
    ├── alpine-windows-amd64

Now, you can run the executable:

    $ dist/alpine-darwin-amd64 cat /etc/alpine-release
    3.9.5

When the executable is run, we'll check for the `alpine:3.9.5` image on the user's system. If it doesn't exist, the executable will automatically run:

    $ docker pull alpine:3.9.5

### Embedded Mode

In this mode, if the specified image doesn't exist, we'll attempt to load it from a tarball that is embeddded in the executable.

    $ docker2exe --name alpine --image alpine:3.9 --embed

When creating the executable above, the image was dumped to a tarball and baked into the resulting executable:

    $ docker save alpine:3.9 | gzip > alpine.tar.gz

When the executable runs, we'll check for the `alpine:3.9` image on the user's system. If it doesn't exist, the executable will automatically run:

    $ docker load alpine.tar.gz

For small images, this approach works great. In the example above, the resulting executable was under 10MB.
