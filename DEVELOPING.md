# Developing

## Getting started

In order to test and develop in this repo you will need the following dependencies installed:
- docker
- make
- git-lfs 

After cloning you will need to do the following:
1. run `git lfs install && git lfs pull` to grab the latest test assets in LFS 
2. run `make bootstrap` to download go mod dependencies, create the `/.tmp` dir, and download helper utilities.
3. run `make` to run linting, tests, and other verifications to make certain everything is working alright.

## Running tests

The main make tasks for common static analysis and testing are:

- `static-analysis`: runs the linter and license checks
- `lint-fix`: let the linter auto-fix as many issues as it can
- `unit`: plain-ol'e unit tests
- `cli`: run trait assertions against snapshot builds

Checkout `make help` to see what other actions you can take.

## Testing with Docker

To create a snapshot and install into the local ` ~/.docker/cli-plugins` directory:

```
make clean-snapshot snapshot install-snapshot
```
