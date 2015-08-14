package captain

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func Run(cmds []*Command, args []string) {
	if c := findSub(cmds, args); c != nil {
		run(c, args, make([]*Command, 0, 3))
	}
}

func run(cmd *Command, args []string, stack []*Command) {
	var subArgs []string
	if cmd.CustomFlags {
		subArgs = args[1:]
	} else {
		flagSet := cmd.Flag
		if cmd.Flag == nil {
			flagSet = flag.NewFlagSet("", flag.ExitOnError)
		}
		flagSet.Usage = usage(cmd)
		flagSet.Parse(args[1:])
		subArgs = flagSet.Args()
	}
	if sub := findSub(cmd.SubCommands, subArgs); sub != nil {
		run(sub, subArgs, append(stack, cmd))
		return
	}
	cmd.Run(stack, args)
}

func findSub(cmds []*Command, args []string) *Command {
	if len(args) == 0 {
		return nil
	}
	for _, cmd := range cmds {
		if name(cmd) == args[0] {
			return cmd
		}
	}
	return nil
}

func usage(c *Command) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
		if len(c.SubCommands) != 0 {
			fmt.Fprintf(os.Stderr, "\nsubcommands:\n")
			for _, cmd := range c.SubCommands {
				fmt.Fprintf(os.Stderr, "\t%s\t%s\n", name(cmd), cmd.Short)
			}
		}
		os.Exit(2)
	}
}

func name(c *Command) string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}
