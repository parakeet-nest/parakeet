package content

import (
	"os"
	"path/filepath"
)

// FindFiles searches for files with a specific extension in the given root directory and its subdirectories.
//
// Parameters:
// - root: The root directory to start the search from.
// - ext: The file extension to search for.
//   examples: ".md", ".html", ".txt", "*.*"
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func FindFiles(dirPath string, ext string) ([]string, error) {
	var textFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			textFiles = append(textFiles, path)
		}
		return nil
	})
	return textFiles, err
}
/*
This is a Go function named FindFiles that searches for files with a specific extension in a given root directory and its subdirectories. 
It takes two parameters: dirPath (the directory path to start the search from) and ext (the file extension to search for). 
It returns a slice of file paths that match the given extension and an error if the search encounters any issues.

The function uses the filepath.Walk function to iterate over all files in the directory and its subdirectories. 
For each file found, it checks if it is not a directory and if its extension matches the given extension. 
If it does, it appends the file path to the textFiles slice.

If there is an error during the search, it is returned. 
Otherwise, the textFiles slice and any error encountered during the search are returned.
*/

// ForEachFile iterates over all files with a specific extension in a directory and its subdirectories.
//
// Parameters:
// - dirPath: The root directory to start the search from.
// - ext: The file extension to search for.
// - callback: A function to be called for each file found.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func ForEachFile(dirPath string, ext string, callback func(string) error) ([]string, error) {
	var textFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			textFiles = append(textFiles, path)
			err = callback(path)
			// generate an error to stop the walk
			if err != nil {
				return err
			}
		}
		return nil
	})
	return textFiles, err
}
/*
This code snippet defines a function called ForEachFile in Go. 
It takes three parameters: dirPath (the root directory to start the search from), ext (the file extension to search for), and callback (a function to be called for each file found).

The function uses the filepath.Walk function to iterate over all files in the directory and its subdirectories. 
For each file found, it checks if it is not a directory and if its extension matches the given extension. 
If it does, it appends the file path to the textFiles slice and calls the callback function with the file path.

The function returns a slice of file paths that match the given extension and an error if the search encounters any issues.
*/

// ReadTextFile reads the contents of a text file at the given path and returns the contents as a string.
//
// Parameters:
// - path: the path to the text file.
//
// Returns:
// - string: the contents of the text file as a string.
// - error: an error if the file cannot be read.
func ReadTextFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteTextFile(path, content string) error {
	// Create a new file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
