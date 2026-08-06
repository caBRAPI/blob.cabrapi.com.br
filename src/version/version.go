package version

// V is the build version, overridden at link time via
// -X blob/src/version.V=<version>. main.go seeds it from its own ldflags var.
var V = "N/A"
