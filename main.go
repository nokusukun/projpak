package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	SyntaxStart = "^<!FILE! path=(.*?)>$"
	SyntaxEnd   = "</!FILE!>"
	SyntaxWrite = "<!FILE! path=%s>\n%s\n</!FILE!>\n"
)

func flattenFiles(directory string, extensions []string, outputFile string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && hasValidExtension(path, extensions) {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(out, SyntaxWrite, path, string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func reconstructFiles(flattenedFile string, outputDir string) error {
	file, err := os.Open(flattenedFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var contentBuilder strings.Builder
	var filePath string
	fileRegex := regexp.MustCompile(SyntaxStart)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := fileRegex.FindStringSubmatch(line); matches != nil {
			// Start of a new file
			if filePath != "" && contentBuilder.Len() > 0 {
				// Write the previously collected file content
				err := writeToFile(filepath.Join(outputDir, filePath), contentBuilder.String())
				if err != nil {
					return err
				}
				contentBuilder.Reset()
			}
			filePath = matches[1]
		} else if strings.TrimSpace(line) == SyntaxEnd {
			// End of the current file
			if filePath != "" {
				err := writeToFile(filepath.Join(outputDir, filePath), contentBuilder.String())
				if err != nil {
					return err
				}
				contentBuilder.Reset()
				filePath = ""
			}
		} else {
			// Accumulate the content of the file
			contentBuilder.WriteString(line + "\n")
		}
	}

	// Write the last file if exists
	if filePath != "" && contentBuilder.Len() > 0 {
		err := writeToFile(filepath.Join(outputDir, filePath), contentBuilder.String())
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func hasValidExtension(filePath string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filePath, ext) {
			return true
		}
	}
	return false
}

func writeToFile(filePath string, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func main() {
	flattenCmd := flag.NewFlagSet("f", flag.ExitOnError)
	reconstructCmd := flag.NewFlagSet("r", flag.ExitOnError)

	flattenDir := flattenCmd.String("dir", "", "Directory to search for files.")
	flattenExtensions := flattenCmd.String("ext", "", "Comma-separated list of file extensions to include, e.g., .go,.txt")
	flattenOutput := flattenCmd.String("output", "", "Output file to write the flattened contents to.")

	reconstructFile := reconstructCmd.String("file", "", "The flattened file to read and reconstruct files from.")
	reconstructOutputDir := reconstructCmd.String("directory", "", "Directory to unpack the reconstructed files to.")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'flatten' or 'reconstruct' command.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "flatten":
		flattenCmd.Parse(os.Args[2:])
		if *flattenDir == "" || *flattenExtensions == "" || *flattenOutput == "" {
			flattenCmd.PrintDefaults()
			os.Exit(1)
		}
		extensions := strings.Split(*flattenExtensions, ",")
		err := flattenFiles(*flattenDir, extensions, *flattenOutput)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "reconstruct":
		reconstructCmd.Parse(os.Args[2:])
		if *reconstructFile == "" {
			reconstructCmd.PrintDefaults()
			os.Exit(1)
		}
		err := reconstructFiles(*reconstructFile, *reconstructOutputDir)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Expected (f)latten or (r)econstruct command.")
		os.Exit(1)
	}
}
