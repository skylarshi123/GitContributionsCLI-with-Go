// Package main is required for executables in Go
package main

// Import required packages for Git operations, time manipulation, and sorting
import (
    // fmt provides formatted I/O operations
    "fmt"
    // sort provides sorting functionality for slices
    "sort"
    // time provides time-related functions
    "time"
    "strings" // for string manipulation
    // go-git packages for Git operations
    // Note: this is an external package, not part of Go standard library
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
)

// Constants used throughout the program
// const is a Go keyword for declaring constants
const outOfRange = 99999                  // Used as a marker for dates too old
const daysInLastSixMonths = 183          // Approximately 6 months worth of days
const weeksInLastSixMonths = 26          // Number of weeks in 6 months

// New features, File Type Stats (most used file types in commits)
type FileTypeStats struct {
    Extension string
    Count     int
}

// type is a Go keyword for declaring new types
// column is a custom type that's really just a slice of integers
type column []int

func getFileExtension(filename string) string {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return "no_extension"
    }
    return "." + parts[len(parts)-1]
}

func processFileTypes(commit *object.Commit, fileTypes map[string]int) error {
    // Get the commit's files
    files, err := commit.Files()
    if err != nil {
        return err
    }

    // Iterate through changed files
    files.ForEach(func(file *object.File) error {
        ext := getFileExtension(file.Name)
        fileTypes[ext]++
        return nil
    })

    return nil
}


// stats is the main entry function for statistics generation
// Takes an email string parameter to filter commits by author
func stats(email string) {
    // Process all repositories and get commit data
    commits, fileTypes := processRepositories(email)
    // Print the statistics in a formatted way
    printCommitsStats(commits)
	printFileTypeStats(fileTypes)
}

// getBeginningOfDay converts a time.Time to the start of that day (00:00:00)
// time.Time is a type from the time package
func getBeginningOfDay(t time.Time) time.Time {
    // Get the date components using the Date() method from time.Time
    year, month, day := t.Date()
    
    // time.Date is a function from time package
    // Creates a new time.Time at 00:00:00 for the given date
    startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
    return startOfDay
}

// countDaysSinceDate counts days between a given date and today
func countDaysSinceDate(date time.Time) int {
    days := 0
    // Get start of today using our helper function
    now := getBeginningOfDay(time.Now())
    
    // Loop until we reach today's date
    for date.Before(now) {
        // Add 24 hours to our date
        // time.Hour is a constant from time package
        date = date.Add(time.Hour * 24)
        days++
        
        // If we're beyond 6 months, return our outOfRange constant
        if days > daysInLastSixMonths {
            return outOfRange
        }
    }
    return days
}

// fillCommits processes a Git repository and counts commits per day and file types
// Parameters:
//   - email: string to filter commits by author email
//   - path: string path to the Git repository
//   - commits: map[int]int to store days-ago -> commit-count mapping
// Returns: 
//   - map[int]int: the updated commits map
//   - map[string]int: counts of file types modified
func fillCommits(email string, path string, commits map[int]int) (map[int]int, map[string]int) {
	// Create a new map to store file extension counts
	// map[string]int where key is file extension (e.g., ".go") and value is count
	fileTypes := make(map[string]int)
 
	// git.PlainOpen comes from go-git package
	// Opens an existing repository at the given path
	// Returns a *git.Repository and error if any
	repo, err := git.PlainOpen(path)
	if err != nil {
		// panic is a built-in Go function that stops program execution
		// Used here because we can't continue without repository access
		panic(err)
	}
 
	// repo.Head() gets the HEAD reference of repository
	// HEAD typically points to the latest commit of current branch
	// Returns a *plumbing.Reference and error if any
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
 
	// repo.Log gets commit history starting from HEAD
	// &git.LogOptions{From: ref.Hash()} specifies starting point
	// ref.Hash() gets the commit hash from the reference
	// Returns a object.CommitIterator for walking through commits
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
 
	// Calculate offset for proper calendar alignment
	// This adjusts commit dates to match GitHub's contribution graph
	offset := calcOffset()
 
	// iterator.ForEach comes from go-git
	// Walks through each commit in history
	// Takes a function to process each commit
	err = iterator.ForEach(func(c *object.Commit) error {
		// Get number of days between commit date and today
		// c.Author.When is the commit timestamp
		daysAgo := countDaysSinceDate(c.Author.When) + offset
 
		// Skip if commit author email doesn't match filter
		// c.Author.Email comes from the commit metadata
		if c.Author.Email != email {
			// Return nil to continue to next commit
			return nil
		}
 
		// If commit is within our time range (not outOfRange)
		if daysAgo != outOfRange {
			// Increment commit count for that day
			commits[daysAgo]++
			
			// Process file types for this commit
			// processFileTypes is our helper function that counts file extensions
			// Pass the current commit and our fileTypes map to update
			if err := processFileTypes(c, fileTypes); err != nil {
				// If there's an error processing files, return it
				// This will stop the commit iteration
				return err
			}
		}
 
		// Return nil to continue processing commits
		return nil
	})
 
	// Check for errors during commit processing
	if err != nil {
		panic(err)
	}
 
	// Return both the commits map and fileTypes map
	// This allows the caller to aggregate statistics across repositories
	return commits, fileTypes
 }
 
 // processRepositories scans all repositories and processes commit data
 // Parameter:
 //   - email: string to filter commits by author
 // Returns: 
 //   - map[int]int: days-ago to commit count mapping
 //   - []FileTypeStats: sorted slice of file extension statistics
 func processRepositories(email string) (map[int]int, []FileTypeStats) {
	// Get path to our repository list file
	// getDotFilePath() is defined in scan.go
	filePath := getDotFilePath()
	
	// Read repository paths from file
	// parseFileLinesToSlice() is defined in scan.go
	repos := parseFileLinesToSlice(filePath)
 
	// Store number of days we're tracking
	daysInMap := daysInLastSixMonths
 
	// Create map for commit counts
	// make is a built-in Go function to create maps
	commits := make(map[int]int, daysInMap)
	
	// Create map for aggregating file type counts across all repositories
	// Key is file extension, value is total count
	allFileTypes := make(map[string]int)
 
	// Initialize all days with zero commits
	// Using reverse loop: daysInMap down to 1
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}
 
	// Process each repository in our list
	// range is a Go keyword for iterating over slices
	for _, path := range repos {
		// Process this repository and get its statistics
		// newCommits: updated commit counts
		// newFileTypes: file type counts from this repo
		newCommits, newFileTypes := fillCommits(email, path, commits)
		
		// Update our commits map with results from this repo
		commits = newCommits
		
		// Merge file type counts from this repo into our total counts
		// range over map of new file types
		for ext, count := range newFileTypes {
			// Add counts to our running totals
			allFileTypes[ext] += count
		}
	}
 
	// Create slice to hold sorted file type statistics
	var fileTypeStats []FileTypeStats
	
	// Convert our map of counts into a slice of FileTypeStats
	// This makes it easier to sort
	for ext, count := range allFileTypes {
		// Append creates a new FileTypeStats struct for each extension
		fileTypeStats = append(fileTypeStats, FileTypeStats{ext, count})
	}
	
	// sort.Slice comes from sort package
	// Sorts our slice of FileTypeStats based on Count field
	// The function provided returns true if element i should come before element j
	sort.Slice(fileTypeStats, func(i, j int) bool {
		// Sort in descending order (higher counts first)
		return fileTypeStats[i].Count > fileTypeStats[j].Count
	})
 
	// Return both the commit counts and sorted file type statistics
	return commits, fileTypeStats
 }

// calcOffset determines how many days to offset for calendar alignment
// Returns: int representing number of days to offset
func calcOffset() int {
    // Declare offset variable we'll calculate
    var offset int

    // time.Now() gets current time
    // Weekday() returns the day of the week (time.Sunday, time.Monday, etc.)
    weekday := time.Now().Weekday()

    // switch is a Go keyword for control flow
    // Similar to if-else but cleaner for multiple cases
    switch weekday {
    case time.Sunday:    // time.Sunday is a constant from time package
        offset = 7
    case time.Monday:
        offset = 6
    case time.Tuesday:
        offset = 5
    case time.Wednesday:
        offset = 4
    case time.Thursday:
        offset = 3
    case time.Friday:
        offset = 2
    case time.Saturday:
        offset = 1
    }

    return offset
}

// printCell formats and prints a single cell of the commit calendar
// Parameters:
//   - val: int representing number of commits
//   - today: bool indicating if this cell represents today
func printCell(val int, today bool) {
	// Initialize default escape code for empty cells
	// \033 is the escape character for terminal formatting
	// [0;37;30m sets foreground and background colors
	escape := "\033[0;37;30m"
 
	// switch with cases for different commit counts
	// Uses Go's special switch syntax with boolean conditions
	switch {
	// 1-4 commits: light color
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	// 5-9 commits: medium color
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	// 10+ commits: dark color
	case val >= 10:
		escape = "\033[1;30;42m"
	}
 
	// Override color if cell represents today
	if today {
		escape = "\033[1;37;45m"
	}
 
	// If no commits, print empty cell with dash
	if val == 0 {
		// fmt.Printf comes from fmt package
		// Prints formatted string with escape codes for color
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}
 
	// Format string for cells with commits
	// Controls spacing based on number of digits
	str := "  %d "
	switch {
	case val >= 10:  // Two-digit numbers
		str = " %d "
	case val >= 100: // Three-digit numbers
		str = "%d "
	}
 
	// Print cell with commit count and proper formatting
	fmt.Printf(escape+str+"\033[0m", val)
 }
 
 // printCommitsStats prints the full commit calendar visualization
 // Parameter:
 //   - commits: map[int]int where key is days-ago and value is commit count
 func printCommitsStats(commits map[int]int) {
	// Get sorted list of day indices
	keys := sortMapIntoSlice(commits)
	// Organize commits into columns (weeks)
	cols := buildCols(keys, commits)
	// Print the formatted calendar
	printCells(cols)
 }
 
 // sortMapIntoSlice converts map keys to sorted slice
 // Parameter:
 //   - m: map[int]int to get keys from
 // Returns: []int slice of sorted keys
 func sortMapIntoSlice(m map[int]int) []int {
	// Create slice to hold keys
	var keys []int
	
	// range over map to extract keys
	// _ ignores the values, we only want keys
	for k := range m {
		keys = append(keys, k)
	}
	
	// sort.Ints comes from sort package
	// Sorts slice of integers in ascending order
	sort.Ints(keys)
 
	return keys
 }
 
 // buildCols organizes commits into a column-based structure
 // Parameters:
 //   - keys: []int slice of sorted day indices
 //   - commits: map[int]int of commit counts
 // Returns: map[int]column where key is week number
 func buildCols(keys []int, commits map[int]int) map[int]column {
	// Create map to store columns
	// Each column represents a week
	cols := make(map[int]column)
	
	// Create empty column slice
	// column is our custom type defined at top
	col := column{}
 
	// Iterate through sorted days
	for _, k := range keys {
		// Calculate week number (0-26)
		// Integer division by 7 gives week number
		week := int(k / 7)
		
		// Calculate day within week (0-6)
		// Modulo 7 gives day of week
		dayinweek := k % 7
 
		// If it's start of week (Sunday)
		// Reset column to empty
		if dayinweek == 0 {
			col = column{}
		}
 
		// Add this day's commit count to current column
		col = append(col, commits[k])
 
		// If it's end of week (Saturday)
		// Save completed column to map
		if dayinweek == 6 {
			cols[week] = col
		}
	}
 
	return cols
 }

 // printCells prints the entire commit calendar visualization
// Parameters:
//   - cols: map[int]column containing organized commit data by weeks
func printCells(cols map[int]column) {
	// First print the month names row at top of calendar
	printMonths()
	
	// Iterate through days of week (top to bottom)
	// 6 to 0 represents Saturday to Sunday
	for j := 6; j >= 0; j-- {
		// Iterate through weeks (right to left)
		// weeksInLastSixMonths+1 to 0 for all weeks plus current
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			// If we're at the start of a row
			// Print the day name (Mon, Wed, etc.)
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			
			// Check if we have data for this week
			// ok is a bool that's true if key exists in map
			if col, ok := cols[i]; ok {
				// Check if this cell represents today
				// Uses calcOffset() to align with GitHub's display
				if i == 0 && j == calcOffset()-1 {
					// Print cell with today's formatting
					printCell(col[j], true)
					continue
				} else {
					// If we have data for this day
					if len(col) > j {
						// Print regular cell
						printCell(col[j], false)
						continue
					}
				}
			}
			// If no data exists, print empty cell
			printCell(0, false)
		}
		// Print newline at end of each row
		// fmt.Printf comes from fmt package
		fmt.Printf("\n")
	}
 }
 
 // printMonths prints the month labels at top of calendar
 // Uses no parameters as it calculates based on current date
 func printMonths() {
	// Calculate start date (6 months ago)
	// time.Now() gets current time
	// getBeginningOfDay converts to start of day
	// Subtract days to get to start date
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	
	// Get initial month to track changes
	// Month() returns time.Month type
	month := week.Month()
	
	// Print initial spacing for alignment
	fmt.Printf("         ")
	
	// Loop through weeks until we reach current date
	for {
		// If month has changed
		if week.Month() != month {
			// Print abbreviated month name (e.g., "Jan")
			// String() converts month to string
			// [:3] takes first 3 characters
			fmt.Printf("%s ", week.Month().String()[:3])
			// Update tracking month
			month = week.Month()
		} else {
			// Print spaces for weeks within same month
			fmt.Printf("    ")
		}
 
		// Add 7 days to move to next week
		week = week.Add(7 * time.Hour * 24)
		
		// If we've passed current date, exit loop
		if week.After(time.Now()) {
			break
		}
	}
	
	// Print newline after month row
	fmt.Printf("\n")
 }
 
 // printDayCol prints the day labels on left side of calendar
 // Parameter:
 //   - day: int representing day of week (0=Sunday, 6=Saturday)
 func printDayCol(day int) {
	// Default to spaces (used for Sun/Tues/Thurs/Sat)
	out := "     "
	
	// switch on day number to determine what to print
	switch day {
	case 1:  // Monday
		out = " Mon "
	case 3:  // Wednesday
		out = " Wed "
	case 5:  // Friday
		out = " Fri "
	}
 
	// Print the day label
	// Note: Some days intentionally left blank for spacing
	fmt.Printf(out)
 }

 func printFileTypeStats(stats []FileTypeStats) {
    fmt.Printf("\nFile Type Statistics:\n")
    fmt.Printf("===================\n")
    
    // Print top 10 or all if less than 10
    limit := 10
    if len(stats) < limit {
        limit = len(stats)
    }
    
    for i := 0; i < limit; i++ {
        stat := stats[i]
        fmt.Printf("%-15s %5d files\n", stat.Extension, stat.Count)
    }
    fmt.Println()
}

 /*
 stats(email)
    │
    ├──► processRepositories(email)
    │       │
    │       ├──► Gets repo list from ~/.gogitlocalstats
    │       │
    │       ├──► For each repository:
    │       │    └──► fillCommits(email, path, commits)
    │       │           │
    │       │           ├──► Opens Git repo
    │       │           ├──► Gets commit history
    │       │           └──► Counts commits per day
    │       │
    │       └──► Returns map[days_ago]commit_count
    │
    └──► printCommitsStats(commits)
            │
            ├──► sortMapIntoSlice (orders days)
            ├──► buildCols (organizes into weeks)
            └──► printCells
                  │
                  ├──► printMonths (top row)
                  ├──► printDayCol (left column)
                  └──► printCell (commit data)
 */

 //With file extension stats
 /*
 stats(email)
    │
    ├──► processRepositories(email)
    │       │
    │       ├──► Gets repo list from ~/.gogitlocalstats
    │       │
    │       ├──► Initialize maps for:
    │       │    - commits (days → count)
    │       │    - allFileTypes (extension → count)
    │       │
    │       ├──► For each repository:
    │       │    └──► fillCommits(email, path, commits)
    │       │           │
    │       │           ├──► Opens Git repo
    │       │           ├──► Gets commit history
    │       │           ├──► For each commit:
    │       │           │    ├──► Count commits per day
    │       │           │    └──► processFileTypes:
    │       │           │         ├──► Get changed files
    │       │           │         ├──► Extract extensions
    │       │           │         └──► Update file type counts
    │       │           │
    │       │           └──► Returns: 
    │       │                ├──► commits map
    │       │                └──► fileTypes map
    │       │
    │       └──► Post-process file stats:
    │           ├──► Convert to FileTypeStats slice
    │           └──► Sort by frequency
    │
    ├──► printCommitsStats(commits)
    │       │
    │       ├──► sortMapIntoSlice (orders days)
    │       ├──► buildCols (organizes into weeks)
    │       └──► printCells
    │             │
    │             ├──► printMonths (top row)
    │             ├──► printDayCol (left column)
    │             └──► printCell (commit data)
    │
    └──► printFileTypeStats(fileTypeStats)
            │
            └──► Displays top 10 file types:
                 Extension    Count
                 .go         123
                 .js          89
                 etc...
 */