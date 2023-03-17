# `zctl`

My personal CLI tool. Use AS-IS and check this README.md for updates to the commands.

- [`zctl`](#zctl)
- [Commands](#commands)
  - [`create`](#create)
    - [`gpgkey`](#gpgkey)
    - [`helm`](#helm)
    - [`k8s`](#k8s)
    - [`sshkey`](#sshkey)
    - [`tlscert`](#tlscert)
  - [`get`](#get)
    - [`base64`](#base64)
    - [`deps`](#deps)
    - [`md5`](#md5)
    - [`sha256`](#sha256)

# Commands

## `create`

### `gpgkey`

Creates a GPG key for use with Github/Gitlab/Bitbucket *etc*.

### `helm`

Creates Kubernetes deployment manifests for use in Helm charts.

### `k8s`

Creates Kubernetes deployment manifests.

### `sshkey`

Creates an SSH key

### `tlscert`

Creates a set of certificates and keys for use in development for TLS.

## `get`

### `base64`

Gets the Base64 encoding of a file or string.

### `deps`

Detects the type of project and automagically runs the standard commands for that language for installing dependencies.

### `md5`

Gets the MD-5 hash of a file or string.

### `sha256`

Gets the SHA-256 hash of a file or string.
