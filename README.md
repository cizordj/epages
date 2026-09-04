# epages

A Cloudflare pages' deployer written in Go.

## Use case

`epages` exists to deploy to Cloudflare Pages from low-end, resource-constrained devices
that can't comfortably run the standard Node.js/`wrangler` toolchain — the Brume 2 being
a prime example.

The Brume 2 technically supports Docker, but it's an older version of it, and the device
doesn't have enough memory to pull and run containers reliably. That rules out the common
workaround of containerizing a proper Node environment to get a modern `wrangler` running.

Running `wrangler` natively on the device isn't a great option either: the Node.js version
that ships with Brume 2 is quite old, and `wrangler` complains about (or outright refuses
to run on) outdated Node runtimes. A vendorized Node environment capable of running
`wrangler` is also heavy for a small, low-memory device.

`epages` solves this by being:

- **Native, no containers** — runs directly on the device, no Docker required at all.
- **Written entirely in Go** — no dependency on OS-provided libraries or a Node.js runtime.
- **Compiled and lightweight** — a single static binary with a low memory footprint,
  making it practical to run on constrained hardware.
- **Fast** — starts up and deploys quicker than spinning up a full Node/`wrangler`
  environment.

In short, `epages` lets you build a deployment pipeline for Cloudflare Pages that runs
natively on devices like the Brume 2, without Docker, without Node.js, and without the
overhead that comes with either.

### Building for ARM64 Statically Linked

To build the project for the ARM64 architecture with static linking, use the
following command:

```sh
env GOOS=linux GOARCH=arm64 make
```
