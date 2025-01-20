# File Explorer

## Description

Showcases the following:

-   Go backend serving the files in this repository
-   React frontend to browse/filter/search the files
-   TLS encryption
-   Authentication
-   Protection against the following attacks:
    -   Directory traversal
    -   Timing attacks

### Development Notes

The Go backend in this repository uses the [`embed`](https://pkg.go.dev/embed)
package to embed the React app inside the Go binary. Running `go build` in the
root will capture whatever is present in the `web/build` subdirectory.

To ensure you have an up to date copy of the web app in your binary, you should:

-   `cd web`
-   `pnpm install`
-   `pnpm build`
-   `cd ..`
-   `go build`

The Go app is hardcoded to listen on port 8080.

For a faster feedback loop and more developer friendly process, you can run
the webapp's dev server alongside the Go backend:

```
$ cd web
$ pnpm start
```
