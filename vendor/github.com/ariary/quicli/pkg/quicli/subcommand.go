package quicli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	stringSlice "github.com/ariary/go-utils/pkg/stringSlice"
)

// Subcommand
type Subcommand struct {
	Name        string
	Description string
	Function    Runner
	// Flags       []Flag
}

type Subcommands []Subcommand

type SubcommandSet []string

// RunWithSubcommand: equivalent of Run function when cli has subcommand defined
func (c *Cli) RunWithSubcommand() {
	var config Config
	usage := new(strings.Builder)
	wUsage := new(tabwriter.Writer)
	wUsage.Init(usage, 2, 8, 1, '\t', 1)
	var shorts []string
	config.Flags = make(map[string]interface{})
	fs := flag.NewFlagSet("parser", flag.ExitOnError)

	//Description
	if isRootCommand(c.Subcommands) {

		if len(c.Subcommands) > 0 {
			subcommandSet := []string{}
			for i := 0; i < len(c.Subcommands); i++ {
				subcommandSet = append(subcommandSet, c.Subcommands[i].Name)
			}
			fmt.Fprintf(wUsage, c.Description+"\n\nUsage: "+c.Usage+"\nAvailable commands: "+strings.Join(subcommandSet, ", ")+"\n\n")
		} else {
			fmt.Fprintf(wUsage, c.Description+"\n\nUsage: "+c.Usage+"\n\n")
		}
	} else {
		//TODO: check if subcommand is misspelled
		sub := getSubcommandByName(c.Subcommands, os.Args[1])
		fmt.Fprintf(wUsage, c.Description+"\n\nUsage: "+c.Usage+"\n"+"Command "+sub.Name+": "+sub.Description+"\n\n")
	}

	//Subcommands preliminary checks
	if len(c.Subcommands) > 0 {
		for i := 0; i < len(c.Subcommands); i++ {
			sub := c.Subcommands[i]
			if sub.Function == nil {
				fmt.Println(QUICLI_ERROR_PREFIX+"subcommand", sub.Name, "does not defined mandatory 'Function' attribute")
				os.Exit(2)
			}
		}
	}

	//flags
	fp := c.Flags
	for i := 0; i < len(fp); i++ {
		f := fp[i]
		// prepation checks
		if len(f.Name) == 0 {
			fmt.Println(QUICLI_ERROR_PREFIX + "empty flag name defintion")
			os.Exit(2)
		}
		//check Default => if no value provided assume it is a bool flag
		if f.Default == nil {
			f.Default = false
		}
		//check if subcommands speicified in flags name are good
		for i := 0; i < len(f.ForSubcommand); i++ {
			subcommandName := f.ForSubcommand[i]
			if getSubcommandByName(c.Subcommands, subcommandName).Name == "" {
				fmt.Println(QUICLI_ERROR_PREFIX+"subcommand", subcommandName, "specified for flag", f.Name, "is not defined")
				os.Exit(2)
			}
		}

		switch f.Default.(type) {
		case int:
			if isRootCommand(c.Subcommands) && !f.NotForRootCommand {
				createIntFlagFs(config, f, &shorts, wUsage, fs)
			} else if len(os.Args) > 1 {
				if f.isForSubcommand(os.Args[1]) {
					createIntFlagFs(config, f, &shorts, wUsage, fs)
				}
			}
		case string:
			if isRootCommand(c.Subcommands) && !f.NotForRootCommand {
				createStringFlagFs(config, f, &shorts, wUsage, fs)
			} else if len(os.Args) > 1 {
				if f.isForSubcommand(os.Args[1]) {
					createStringFlagFs(config, f, &shorts, wUsage, fs)
				}
			}
		case bool:
			if isRootCommand(c.Subcommands) && !f.NotForRootCommand {
				createBoolFlagFs(config, f, &shorts, wUsage, fs)
			} else if len(os.Args) > 1 {
				if f.isForSubcommand(os.Args[1]) {
					createBoolFlagFs(config, f, &shorts, wUsage, fs)
				}
			}
		case float64:
			if isRootCommand(c.Subcommands) && !f.NotForRootCommand {
				createFloatFlagFs(config, f, &shorts, wUsage, fs)
			} else if len(os.Args) > 1 && f.isForSubcommand(os.Args[1]) {
				createFloatFlagFs(config, f, &shorts, wUsage, fs)
			}
		default:
			fmt.Println(QUICLI_ERROR_PREFIX+"Unknown flag type:", f.Default)
			os.Exit(2)
		}
	}
	fmt.Fprintf(wUsage, "\nUse \""+os.Args[0]+" --help\" for more information about the command.\n")

	//cheat sheet pt1
	var cheatSheet bool
	if len(c.CheatSheet) > 0 {
		fmt.Fprintf(wUsage, "\nSee command examples with \""+os.Args[0]+" --cheat-sheet\"\n")
		flag.BoolVar(&cheatSheet, "cheat-sheet", false, "print cheat sheet")
		flag.BoolVar(&cheatSheet, "cs", false, "print cheat sheet")
	}

	wUsage.Flush()
	// Parse
	fs.Usage = func() { fmt.Print(usage.String()) }
	if isRootCommand(c.Subcommands) && len(os.Args) > 1 {
		fs.Parse(os.Args[1:])
	} else if len(os.Args) > 2 {
		fs.Parse(os.Args[2:])
	}
	config.Args = fs.Args()

	//cheat sheet pt2
	if len(c.CheatSheet) > 0 && cheatSheet {
		c.PrintCheatSheet()
		os.Exit(0)
	}

	// Run
	if isRootCommand(c.Subcommands) {
		c.Function(config)
	} else {
		getSubcommandByName(c.Subcommands, os.Args[1]).Function(config)
	}
}

// isRootCommand: return true if the command line is targetting the root command, false if it is targgeting a subcommand
func isRootCommand(subcommands Subcommands) bool {
	if len(os.Args) < 2 {
		return true
	} else {
		sub := getSubcommandByName(subcommands, os.Args[1])
		return sub.Name == ""
	}
}

// getSubcommandByName: return true if the command line is targetting the root command, false if it is targgeting a subcommand
func getSubcommandByName(subcommands Subcommands, subcommandName string) (sub Subcommand) {

	for i := 0; i < len(subcommands); i++ {
		if subcommandName == subcommands[i].Name {
			return subcommands[i]
		}
	}
	return sub
}

// isForSubcommand: return true if the subcommand is concerned by the flag
func (f *Flag) isForSubcommand(subcommandName string) bool {
	for i := 0; i < len(f.ForSubcommand); i++ {
		if subcommandName == f.ForSubcommand[i] {
			return true
		}
	}
	return false
}

func createIntFlagFs(cfg Config, f Flag, shorts *[]string, wUsage *tabwriter.Writer, fs *flag.FlagSet) {
	name := f.Name
	shortName := name[0:1]
	var intPtr int
	fs.IntVar(&intPtr, name, int(reflect.ValueOf(f.Default).Int()), f.Description)
	if !stringSlice.Contains(*shorts, shortName) && !f.NoShortName {
		fs.IntVar(&intPtr, shortName, int(reflect.ValueOf(f.Default).Int()), f.Description)
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, shortName))
		cfg.Flags[shortName] = &intPtr
		*shorts = append(*shorts, shortName)
	} else {
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, ""))
	}
	cfg.Flags[name] = &intPtr
}

func createStringFlagFs(cfg Config, f Flag, shorts *[]string, wUsage *tabwriter.Writer, fs *flag.FlagSet) {
	name := f.Name
	shortName := name[0:1]
	var strPtr string
	fs.StringVar(&strPtr, name, string(reflect.ValueOf(f.Default).String()), f.Description)
	if !stringSlice.Contains(*shorts, shortName) && !f.NoShortName {
		fs.StringVar(&strPtr, shortName, string(reflect.ValueOf(f.Default).String()), f.Description)
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, shortName))
		cfg.Flags[shortName] = &strPtr
		*shorts = append(*shorts, shortName)
	} else {
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, ""))
	}
	cfg.Flags[name] = &strPtr
}

func createBoolFlagFs(cfg Config, f Flag, shorts *[]string, wUsage *tabwriter.Writer, fs *flag.FlagSet) {
	name := f.Name
	shortName := name[0:1]
	var bPtr bool
	fs.BoolVar(&bPtr, name, bool(reflect.ValueOf(f.Default).Bool()), f.Description)
	cfg.Flags[name] = &bPtr
	if !stringSlice.Contains(*shorts, shortName) && !f.NoShortName {
		fs.BoolVar(&bPtr, shortName, bool(reflect.ValueOf(f.Default).Bool()), f.Description)
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, shortName))
		cfg.Flags[shortName] = &bPtr
		*shorts = append(*shorts, shortName)
	} else {
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, ""))
	}
	cfg.Flags[name] = &bPtr
}

func createFloatFlagFs(cfg Config, f Flag, shorts *[]string, wUsage *tabwriter.Writer, fs *flag.FlagSet) {
	name := f.Name
	shortName := name[0:1]
	var floatPtr float64
	fs.Float64Var(&floatPtr, name, float64(reflect.ValueOf(f.Default).Float()), f.Description)
	cfg.Flags[name] = &floatPtr
	if !stringSlice.Contains(*shorts, shortName) && !f.NoShortName {
		fs.Float64Var(&floatPtr, shortName, float64(reflect.ValueOf(f.Default).Float()), f.Description)
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, shortName))
		cfg.Flags[shortName] = &floatPtr
		*shorts = append(*shorts, shortName)
	} else {
		fmt.Fprintf(wUsage, getFlagLine(f.Description, f.Default, name, ""))
	}
	cfg.Flags[name] = &floatPtr
}
