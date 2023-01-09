The resource folder is an example of an embedded resource. The Go language supports embedded resources, just
like a Windows program, that are loaded into application types at application startup

In this case, the entire resource directory will be accessible as a file system in the application. So, this
is a way to mount an in-memory file system.

Here is the link to the package : https://pkg.go.dev/embed

Be careful on formatting if returning a http response. There needs to be blank lines between the header and body.
Look at http-503.txt, which ends with a blank line and no body. The last blank line is required