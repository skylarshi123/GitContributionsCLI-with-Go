// Package declaration - every Go file must start with a package name
// 'main' is a special package name that tells Go this is an executable program
package main

// Import block - brings in external packages needed by this file
// The 'flag' package is used to handle command-line arguments
import (
    "flag"
)

// main() function is the entry point of the program
// Every executable Go program must have exactly one main() function
func main() {
    // Declare variables to store command-line arguments
    // These are string variables that will hold the folder path and email
    var folder string
    var email string
    
    // flag.StringVar sets up command-line flags with default values
    // Parameters are:
    // 1. Pointer to variable where flag value will be stored
    // 2. Flag name to use on command line (e.g., -add)
    // 3. Default value if flag is not provided
    // 4. Help text describing the flag
    flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
    flag.StringVar(&email, "email", "your@email.com", "the email to scan")
    
    // Parse the command-line flags
    // This must be called after flags are defined but before they are accessed
    flag.Parse()
    
    // If a folder was provided (flag -add was used)
    if folder != "" {
        // Call scan() function with the folder path and return
        scan(folder)
        return
    }
    
    // If no folder was provided, call stats() with the email
    // This is the default behavior when run without the -add flag
    stats(email)
}