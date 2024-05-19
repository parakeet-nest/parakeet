package content

import (
	"os"
)

// GetArrayOfContentFiles searches for files with a specific extension in the given directory and its subdirectories.
//
// Parameters:
// - dirPath: The directory path to start the search from.
// - ext: The file extension to search for.
//
// Returns:
// - []string: A slice of file paths that match the given extension.
// - error: An error if the search encounters any issues.
func GetArrayOfContentFiles(dirPath string, ext string) ([]string, error) {
	content := []string{}
	_, err := ForEachFile(dirPath, ext, func(path string) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content = append(content, string(data))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return content, nil
}

/*
This is a Go function named GetArrayOfContentFiles that searches for files with a specific extension in a given directory and its subdirectories.
It takes two parameters: dirPath (the directory path to start the search from) and ext (the file extension to search for).
It returns a slice of file paths that match the given extension and an error if the search encounters any issues.
The function uses the ForEachFile function to iterate over all files with the given extension in the directory and its subdirectories.
For each file found, it reads the file's content using os.ReadFile and appends it to the content slice.
If there is an error reading the file, it returns the error.
Finally, it returns the content slice and any error encountered during the search.
*/

// GetMapOfContentFiles searches for files with a specific extension in the given directory and its subdirectories.
//
// Parameters:
// - dirPath: The directory path to start the search from.
// - ext: The file extension to search for.
//
// Returns:
// - map[string]string: A map of file paths to their contents, where the keys are the base names of the files and the values are the file contents as strings.
// - error: An error if the search encounters any issues.
func GetMapOfContentFiles(dirPath string, ext string) (map[string]string, error) {
	content := map[string]string{}
	_, err := ForEachFile(dirPath, ext, func(path string) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content[path] = string(data)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return content, nil
}

/*
This is a Go function named GetMapOfContentFiles that searches for files with a specific extension in a given directory and its subdirectories.
It takes two parameters: dirPath (the directory path to start the search from) and ext (the file extension to search for).
It returns a map of file paths to their contents, where the keys are the path of the files and the values are the file contents as strings.
If the search encounters any issues, it returns an error.

The function uses the ForEachFile function to iterate over all files with the given extension in the directory and its subdirectories.
For each file found, it reads the file's content using os.ReadFile and adds it to the content map with the path of the file as the key and the file content as the value.
If there is an error reading the file, it returns the error.
Finally, it returns the content map and any error encountered during the search.
*/
