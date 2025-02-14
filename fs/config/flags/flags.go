// Package flags contains enhanced versions of spf13/pflag flag
// routines which will read from the environment also.
package flags

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Evengard/rclone/fs"
	"github.com/spf13/pflag"
)

// Groups of Flags
type Groups struct {
	// Groups of flags
	Groups []*Group

	// GroupsMaps maps a group name to a Group
	ByName map[string]*Group
}

// Group related flags together for documentation purposes
type Group struct {
	Groups *Groups
	Name   string
	Help   string
	Flags  *pflag.FlagSet
}

// NewGroups constructs a collection of Groups
func NewGroups() *Groups {
	return &Groups{
		ByName: make(map[string]*Group),
	}
}

// NewGroup to create Group
func (gs *Groups) NewGroup(name, help string) *Group {
	group := &Group{
		Name:  name,
		Help:  help,
		Flags: pflag.NewFlagSet(name, pflag.ExitOnError),
	}
	gs.ByName[group.Name] = group
	gs.Groups = append(gs.Groups, group)
	return group
}

// Filter makes a copy of groups filtered by flagsRe
func (gs *Groups) Filter(flagsRe *regexp.Regexp) *Groups {
	newGs := NewGroups()
	for _, g := range gs.Groups {
		newG := newGs.NewGroup(g.Name, g.Help)
		g.Flags.VisitAll(func(f *pflag.Flag) {
			matched := flagsRe == nil || flagsRe.MatchString(f.Name) || flagsRe.MatchString(f.Usage)
			if matched {
				newG.Flags.AddFlag(f)
			}
		})
	}
	return newGs
}

// Include makes a copy of groups only including the ones in the filter string
func (gs *Groups) Include(groupsString string) *Groups {
	if groupsString == "" {
		return gs
	}
	want := map[string]bool{}
	for _, groupName := range strings.Split(groupsString, ",") {
		_, ok := All.ByName[groupName]
		if !ok {
			log.Fatalf("Couldn't find group %q in command annotation", groupName)
		}
		want[groupName] = true
	}
	newGs := NewGroups()
	for _, g := range gs.Groups {
		if !want[g.Name] {
			continue
		}
		newG := newGs.NewGroup(g.Name, g.Help)
		newG.Flags = g.Flags
	}
	return newGs
}

// Add a new flag to the Group
func (g *Group) Add(flag *pflag.Flag) {
	g.Flags.AddFlag(flag)
}

// AllRegistered returns all flags in a group
func (gs *Groups) AllRegistered() map[*pflag.Flag]struct{} {
	out := make(map[*pflag.Flag]struct{})
	for _, g := range gs.Groups {
		g.Flags.VisitAll(func(f *pflag.Flag) {
			out[f] = struct{}{}
		})
	}
	return out
}

// All is the global stats Groups
var All *Groups

// Groups of flags for documentation purposes
func init() {
	All = NewGroups()
	All.NewGroup("Copy", "Flags for anything which can Copy a file.")
	All.NewGroup("Sync", "Flags just used for `rclone sync`.")
	All.NewGroup("Important", "Important flags useful for most commands.")
	All.NewGroup("Check", "Flags used for `rclone check`.")
	All.NewGroup("Networking", "General networking and HTTP stuff.")
	All.NewGroup("Performance", "Flags helpful for increasing performance.")
	All.NewGroup("Config", "General configuration of rclone.")
	All.NewGroup("Debugging", "Flags for developers.")
	All.NewGroup("Filter", "Flags for filtering directory listings.")
	All.NewGroup("Listing", "Flags for listing directories.")
	All.NewGroup("Logging", "Logging and statistics.")
	All.NewGroup("Metadata", "Flags to control metadata.")
	All.NewGroup("RC", "Flags to control the Remote Control API.")
}

// installFlag constructs a name from the flag passed in and
// sets the value and default from the environment if possible
// the value may be overridden when the command line is parsed
//
// Used to create non-backend flags like --stats.
//
// It also adds the flag to the groups passed in.
func installFlag(flags *pflag.FlagSet, name string, groupsString string) {
	// Find flag
	flag := flags.Lookup(name)
	if flag == nil {
		log.Fatalf("Couldn't find flag --%q", name)
	}

	// Read default from environment if possible
	envKey := fs.OptionToEnv(name)
	if envValue, envFound := os.LookupEnv(envKey); envFound {
		err := flags.Set(name, envValue)
		if err != nil {
			log.Fatalf("Invalid value when setting --%s from environment variable %s=%q: %v", name, envKey, envValue, err)
		}
		fs.Debugf(nil, "Setting --%s %q from environment variable %s=%q", name, flag.Value, envKey, envValue)
		flag.DefValue = envValue
	}

	// Add flag to Group if it is a global flag
	if flags == pflag.CommandLine {
		for _, groupName := range strings.Split(groupsString, ",") {
			if groupName == "rc-" {
				groupName = "RC"
			}
			group, ok := All.ByName[groupName]
			if !ok {
				log.Fatalf("Couldn't find group %q for flag --%s", groupName, name)
			}
			group.Add(flag)
		}
	}
}

// SetDefaultFromEnv constructs a name from the flag passed in and
// sets the default from the environment if possible
//
// Used to create backend flags like --skip-links
func SetDefaultFromEnv(flags *pflag.FlagSet, name string) {
	envKey := fs.OptionToEnv(name)
	envValue, found := os.LookupEnv(envKey)
	if found {
		flag := flags.Lookup(name)
		if flag == nil {
			log.Fatalf("Couldn't find flag --%q", name)
		}
		fs.Debugf(nil, "Setting default for %s=%q from environment variable %s", name, envValue, envKey)
		//err = tempValue.Set()
		flag.DefValue = envValue
	}
}

// StringP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.StringP
func StringP(name, shorthand string, value string, usage string, groups string) (out *string) {
	out = pflag.StringP(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// StringVarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.StringVarP
func StringVarP(flags *pflag.FlagSet, p *string, name, shorthand string, value string, usage string, groups string) {
	flags.StringVarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// BoolP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.BoolP
func BoolP(name, shorthand string, value bool, usage string, groups string) (out *bool) {
	out = pflag.BoolP(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// BoolVarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.BoolVarP
func BoolVarP(flags *pflag.FlagSet, p *bool, name, shorthand string, value bool, usage string, groups string) {
	flags.BoolVarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// IntP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.IntP
func IntP(name, shorthand string, value int, usage string, groups string) (out *int) {
	out = pflag.IntP(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// Int64P defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.IntP
func Int64P(name, shorthand string, value int64, usage string, groups string) (out *int64) {
	out = pflag.Int64P(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// Int64VarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.Int64VarP
func Int64VarP(flags *pflag.FlagSet, p *int64, name, shorthand string, value int64, usage string, groups string) {
	flags.Int64VarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// IntVarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.IntVarP
func IntVarP(flags *pflag.FlagSet, p *int, name, shorthand string, value int, usage string, groups string) {
	flags.IntVarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// Uint32VarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.Uint32VarP
func Uint32VarP(flags *pflag.FlagSet, p *uint32, name, shorthand string, value uint32, usage string, groups string) {
	flags.Uint32VarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// Float64P defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.Float64P
func Float64P(name, shorthand string, value float64, usage string, groups string) (out *float64) {
	out = pflag.Float64P(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// Float64VarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.Float64VarP
func Float64VarP(flags *pflag.FlagSet, p *float64, name, shorthand string, value float64, usage string, groups string) {
	flags.Float64VarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// DurationP defines a flag which can be set by an environment variable
//
// It wraps the duration in an fs.Duration for extra suffixes and
// passes it to pflag.VarP
func DurationP(name, shorthand string, value time.Duration, usage string, groups string) (out *time.Duration) {
	out = new(time.Duration)
	*out = value
	pflag.VarP((*fs.Duration)(out), name, shorthand, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// DurationVarP defines a flag which can be set by an environment variable
//
// It wraps the duration in an fs.Duration for extra suffixes and
// passes it to pflag.VarP
func DurationVarP(flags *pflag.FlagSet, p *time.Duration, name, shorthand string, value time.Duration, usage string, groups string) {
	flags.VarP((*fs.Duration)(p), name, shorthand, usage)
	installFlag(flags, name, groups)
}

// VarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.VarP
func VarP(value pflag.Value, name, shorthand, usage string, groups string) {
	pflag.VarP(value, name, shorthand, usage)
	installFlag(pflag.CommandLine, name, groups)
}

// FVarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.VarP
func FVarP(flags *pflag.FlagSet, value pflag.Value, name, shorthand, usage string, groups string) {
	flags.VarP(value, name, shorthand, usage)
	installFlag(flags, name, groups)
}

// StringArrayP defines a flag which can be set by an environment variable
//
// It sets one value only - command line flags can be used to set more.
//
// It is a thin wrapper around pflag.StringArrayP
func StringArrayP(name, shorthand string, value []string, usage string, groups string) (out *[]string) {
	out = pflag.StringArrayP(name, shorthand, value, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// StringArrayVarP defines a flag which can be set by an environment variable
//
// It sets one value only - command line flags can be used to set more.
//
// It is a thin wrapper around pflag.StringArrayVarP
func StringArrayVarP(flags *pflag.FlagSet, p *[]string, name, shorthand string, value []string, usage string, groups string) {
	flags.StringArrayVarP(p, name, shorthand, value, usage)
	installFlag(flags, name, groups)
}

// CountP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.CountP
func CountP(name, shorthand string, usage string, groups string) (out *int) {
	out = pflag.CountP(name, shorthand, usage)
	installFlag(pflag.CommandLine, name, groups)
	return out
}

// CountVarP defines a flag which can be set by an environment variable
//
// It is a thin wrapper around pflag.CountVarP
func CountVarP(flags *pflag.FlagSet, p *int, name, shorthand string, usage string, groups string) {
	flags.CountVarP(p, name, shorthand, usage)
	installFlag(flags, name, groups)
}
