package go_minion

import (
  "github.com/bobziuchkovski/writ"
  "os"
)

type Arguments struct {
  HelpFlag bool   `flag:"help" description:"Display this help message and exit"`
}

// ParseArguments should be calldd with something that embeds the
// Arguments type and adds more command line arguments.
// Be aware: If an error occurs during parsing or the HelpFlag was specified,
// this method will terminate the application!
func ParseArguments(args *Arguments, name string) {
  cmd := writ.New(name, args)

  // Parse command line arguments
  _, _, err := cmd.Decode(os.Args[1:])
  if err != nil || args.HelpFlag {
    cmd.ExitHelp(err)
  }
}
