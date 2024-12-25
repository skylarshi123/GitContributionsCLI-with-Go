// Every Go file must declare which package it belongs to
// The 'main' package is special - it's for executable programs
package main

// Import statements declare which external packages we need
// Each import gives us access to different functionality
import (
    "bufio"      // Package for buffered I/O - helps read files efficiently
    "fmt"        // Package for formatted I/O - like Printf, Println
    "io"         // Package providing basic I/O interfaces like EOF
    "io/ioutil"  // Package with simplified I/O utilities
    "log"        // Package for logging functionality
    "os"         // Package for operating system functionality like file operations
    "os/user"    // Package for user account information
    "strings"    // Package for string manipulation functions
)

// getDotFilePath returns the path where we'll store our repository list
// It creates a hidden file in the user's home directory
func getDotFilePath() string {
    // user.Current() comes from the os/user package
    // It returns information about the current user account
    // Returns a *user.User struct containing user info and an error if something goes wrong
    usr, err := user.Current()
    
    // If there was an error getting user information
    if err != nil {
        // log.Fatal comes from the log package
        // It prints the error and exits the program immediately
        log.Fatal(err)
    }

    // usr.HomeDir comes from the *user.User struct
    // It contains the path to the current user's home directory
    // We add /.gogitlocalstats to create our hidden file path
    // The dot prefix makes it hidden on Unix-like systems
    dotFile := usr.HomeDir + "/.gogitlocalstats"

    // Return the complete file path
    return dotFile
}

// openFile attempts to open a file at the given path
// If the file doesn't exist, it creates it
// Returns a pointer to os.File which we can use for reading/writing
func openFile(filePath string) *os.File {
    // os.OpenFile comes from the os package
    // Parameters:
    // - filePath: where to find/create the file
    // - os.O_APPEND|os.O_RDWR: flags saying we want to both read and append
    // - 0755: Unix permissions (owner: rwx, others: rx)
    f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0755)
    
    // If there was an error opening the file
    if err != nil {
        // os.IsNotExist comes from os package
        // Checks if the error was because the file doesn't exist
        if os.IsNotExist(err) {
            // os.Create comes from os package
            // Creates a new file at filePath
            // Returns a file handle and possibly an error
            _, err = os.Create(filePath)
            if err != nil {
                // If we couldn't create the file, panic
                // panic stops program execution immediately
                panic(err)
            }
        } else {
            // If there was any other kind of error, panic
            panic(err)
        }
    }

    // Return the file handle, will be null if we created the file
    return f
}

// parseFileLinesToSlice reads a file and returns each line as an element in a string slice
func parseFileLinesToSlice(filePath string) []string {
    // Open the file using our helper function
    f := openFile(filePath)
    
    // defer comes from Go itself
    // It ensures f.Close() is called when this function returns
    // This is important for cleaning up system resources
    defer f.Close()

    // Create an empty slice of strings to store our lines, a slice is a dynamic array, if arr you would put a number in the bracket
    var lines []string
    
    // bufio.NewScanner comes from bufio package
    // Creates a Scanner that can read the file line by line efficiently
    scanner := bufio.NewScanner(f)
    
    // scanner.Scan() reads next line, returns false when done
    // This is a common Go pattern for reading things
    for scanner.Scan() {
        // scanner.Text() gets the string content of the current line
        // append adds it to our slice of lines
        lines = append(lines, scanner.Text())
    }
    
    // Check if there were any errors during scanning
    // scanner.Err() returns any error that occurred
    if err := scanner.Err(); err != nil { //go allows you to initialize and check in one line, the check is the thing after the semi colon
        // io.EOF is a special error from the io package
        // It means we reached the end of the file (not a real error)
        if err != io.EOF {
            panic(err)
        }
    }

    // Return all the lines we read
    return lines
}

// sliceContains checks if a string exists in a slice of strings
// Returns true if found, false if not
func sliceContains(slice []string, value string) bool {
    // Range over the slice
    // _ means we don't care about the index
    // v gets each value in the slice
    for _, v := range slice {
        // If we found the value, return true immediately
        if v == value {
            return true
        }
    }
    // If we get here, we didn't find the value
    return false
}

// joinSlices combines two string slices while avoiding duplicates
// 'new' and 'existing' are just parameter names we chose, not special words
func joinSlices(new []string, existing []string) []string {
    // Range through each item in the new slice
    // 'for' is a Go keyword
    // 'i' is just a variable name we chose for each item
    for _, i := range new {
        // sliceContains is our own function defined above
        // '!' is the built-in NOT operator
        if !sliceContains(existing, i) {
            // append is a built-in Go function for adding to slices
            // It's not from any package, it's part of Go itself
            existing = append(existing, i)
        }
    }
    // Return the modified slice
    return existing
}

// dumpStringsSliceToFile writes a slice of strings to a file
// Each string becomes a line in the file
func dumpStringsSliceToFile(repos []string, filePath string) {
    // strings.Join comes from the strings package
    // It combines all strings in the slice with \n between them
    content := strings.Join(repos, "\n")
    
    // ioutil.WriteFile comes from the ioutil package
    // []byte() is a built-in Go type conversion
    // 0755 is the Unix file permissions
    ioutil.WriteFile(filePath, []byte(content), 0755)
}

// addNewSliceElementsToFile combines existing and new repository paths
// and writes them to the tracking file
func addNewSliceElementsToFile(filePath string, newRepos []string) {
    // Call our own functions defined above
    existingRepos := parseFileLinesToSlice(filePath)
    repos := joinSlices(newRepos, existingRepos)
    dumpStringsSliceToFile(repos, filePath)
}

// recursiveScanFolder starts the repository scanning process
func recursiveScanFolder(folder string) []string {
    // make is a built-in Go function for creating slices
    // Here we make a string slice with initial size 0
    return scanGitFolders(make([]string, 0), folder)
}

// scan is the main scanning function that users will call
func scan(folder string) {
    // fmt.Printf comes from fmt package
    // \n is the newline character
    fmt.Printf("Found folders:\n\n")
    
    // Call our own functions defined in this file
    repositories := recursiveScanFolder(folder)
    filePath := getDotFilePath()
    addNewSliceElementsToFile(filePath, repositories)
    
    fmt.Printf("\n\nSuccessfully added\n\n")
}

// scanGitFolders recursively looks for Git repositories in directories
func scanGitFolders(folders []string, folder string) []string {
    // strings.TrimSuffix comes from strings package
    // Removes the trailing "/" if it exists
    folder = strings.TrimSuffix(folder, "/")

    // os.Open comes from os package
    // Returns a file handle and possibly an error
    f, err := os.Open(folder)
    if err != nil {
        // log.Fatal comes from log package
        // Prints error and exits program
        log.Fatal(err)
    }
    
    // f.Readdir comes from os package's File type
    // -1 means read all entries
    // Returns slice of FileInfo and possibly an error
    files, err := f.Readdir(-1)
    
    // Close the directory handle
    // This is from the os.File type
    f.Close()
    
    if err != nil {
        log.Fatal(err)
    }

    // var is a Go keyword for declaring variables
    // string is a built-in Go type
    var path string

    // Iterate through all files/directories found
    for _, file := range files {
        // file.IsDir() is a method from os.FileInfo interface
        // Checks if the entry is a directory
        if file.IsDir() {
            // + is the built-in string concatenation operator
            path = folder + "/" + file.Name()
            
            // file.Name() comes from os.FileInfo interface
            // Gets the name of the file/directory
            if file.Name() == ".git" {
                path = strings.TrimSuffix(path, "/.git")
                fmt.Println(path)
                folders = append(folders, path)
                // continue is a Go keyword
                // Skips rest of loop, starts next iteration
                continue
            }
            
            // Skip vendor and node_modules directories
            if file.Name() == "vendor" || file.Name() == "node_modules" {
                continue
            }
            
            // Recursive call to scan subdirectories
            folders = scanGitFolders(folders, path)
        }
    }

    return folders
}

/*
scan(folder)
    │
    ├──► recursiveScanFolder(folder)
    │       │
    │       └──► scanGitFolders(empty_slice, folder)
    │               │
    │               ├──► Recursively searches directories
    │               └──► Returns list of found Git repos
    │
    ├──► getDotFilePath()
    │       └──► Gets/creates ~/.gogitlocalstats
    │
    └──► addNewSliceElementsToFile()
            │
            ├──► parseFileLinesToSlice (reads existing repos)
            ├──► joinSlices (combines new & existing repos)
            └──► dumpStringsSliceToFile (writes back to file)
*/