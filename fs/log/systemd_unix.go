// Systemd interface for Unix variants only

//go:build unix
// +build unix

package log

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Evengard/rclone/fs"
	"github.com/coreos/go-systemd/v22/journal"
)

// Enables systemd logs if configured or if auto-detected
func startSystemdLog() bool {
	flagsStr := "," + Opt.Format + ","
	var flags int
	if strings.Contains(flagsStr, ",longfile,") {
		flags |= log.Llongfile
	}
	if strings.Contains(flagsStr, ",shortfile,") {
		flags |= log.Lshortfile
	}
	log.SetFlags(flags)
	// TODO: Use the native journal.Print approach rather than a custom implementation
	fs.LogPrint = func(level fs.LogLevel, text string) {
		text = fmt.Sprintf("<%s>%-6s: %s", systemdLogPrefix(level), level, text)
		_ = log.Output(4, text)
	}
	return true
}

var logLevelToSystemdPrefix = []journal.Priority{
	fs.LogLevelEmergency: journal.PriEmerg,
	fs.LogLevelAlert:     journal.PriAlert,
	fs.LogLevelCritical:  journal.PriCrit,
	fs.LogLevelError:     journal.PriErr,
	fs.LogLevelWarning:   journal.PriWarning,
	fs.LogLevelNotice:    journal.PriNotice,
	fs.LogLevelInfo:      journal.PriInfo,
	fs.LogLevelDebug:     journal.PriDebug,
}

func systemdLogPrefix(l fs.LogLevel) string {
	if l >= fs.LogLevel(len(logLevelToSystemdPrefix)) {
		return ""
	}
	return strconv.Itoa(int(logLevelToSystemdPrefix[l]))
}

func isJournalStream() bool {
	if usingJournald, _ := journal.StderrIsJournalStream(); usingJournald {
		return true
	}

	return false
}
