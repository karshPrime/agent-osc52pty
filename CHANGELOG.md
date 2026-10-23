# Changelog

## added/modified:

* chore/modernise-toolchain: builds with current Go, installs with `go install github.com/karnull/osc52pty@latest`, and lints with go vet and staticcheck in place of golint.
* chore/ci: pull requests are checked in CI with go vet, staticcheck and the tests on macOS and Linux, with pins kept current by Renovate.
* feature/release-binaries: a new version in VERSION, merged to master, publishes macOS binaries for Intel and Apple silicon with the pull request's bullet points as release notes.
* fix/output-lost-at-exit: the last output of a command, and any clipboard sequence in it, is no longer cut off when the command exits.
