package version

import (
	"github.com/go-logr/logr"
)

var (
	Version = "0.0.1"
)

// These fields are set during an official build
// Global vars set from command-line arguments
var (
	BuildVersion = "--"
	BuildHash    = "--"
	BuildTime    = "--"
)

//PrintVersionInfo displays the kyverno version - git version
func PrintVersionInfo(log logr.Logger) {
	log.Info("Oldmonk", "Version", BuildVersion)
	log.Info("Oldmonk", "BuildHash", BuildHash)
	log.Info("Oldmonk", "BuildTime", BuildTime)
}
