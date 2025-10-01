package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/diiyw/dw/template"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

var projectName = ""
var appId *string

// initProjectName finds and initializes the project name from _wasm.js files
func initProjectName(source string, targetDir string) string {
	// Find all _wasm.js files in the current directory
	files, err := filepath.Glob(filepath.Join(source, "*_wasm.js"))
	if err != nil {
		fmt.Printf(Red+"Error finding _wasm.js files: %v\n"+Reset, err)
		return ""
	}

	if len(files) == 0 {
		fmt.Println(Red + "No _wasm.js files found" + Reset)
		return ""
	}

	// Use the first found file
	filename := files[0]
	// Remove the "_wasm.js" suffix
	projectName := strings.TrimSuffix(filename, "_wasm.js")
	projectName = filepath.Base(projectName)

	fmt.Printf(Cyan+"Project name initialized to: %s\n"+Reset, projectName)
	return projectName
}

// ensureDir ensures that a directory exists, creating it if necessary
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// copyGameAssetFiles copies a file from source to destination
func copyGameAssetFiles(sourcePath, destPath string) error {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", sourcePath, err)
	}

	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", destPath, err)
	}

	return nil
}

func main() {
	// Define command line arguments
	sourceDir := flag.String("source", "", "Source directory path")
	targetDir := flag.String("target", "", "Target directory path")
	appId = flag.String("appid", "", "WeChat mini-game appId")

	// Parse command line arguments
	flag.Parse()

	// Check required arguments
	if *sourceDir == "" || *targetDir == "" || *appId == "" {
		fmt.Println(Red + "Please provide source directory, target directory, and WeChat mini-game appId" + Reset)
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf(Blue+"Source directory: %s\n"+Reset, *sourceDir)
	fmt.Printf(Blue+"Target directory: %s\n"+Reset, *targetDir)
	fmt.Printf(Blue+"AppId: %s\n"+Reset, *appId)

	// Initialize projectName using the filename without _wasm.js suffix
	projectName = initProjectName(*sourceDir, *targetDir)

	// Construct the final target directory path: targetDir + projectName
	finalTargetDir := filepath.Join(*targetDir, projectName)
	fmt.Printf(Yellow+"Final target directory: %s\n"+Reset, finalTargetDir)

	// Ensure the final target directory exists
	if err := ensureDir(finalTargetDir); err != nil {
		fmt.Printf(Red+"Error creating final target directory: %v\n"+Reset, err)
		return
	}

	// Copy template files to the final target directory
	copyTemplate(finalTargetDir)

	// Call each processing function using the final target directory
	modifyWasmJS(*sourceDir, finalTargetDir)
	renameArchiveFiles(*sourceDir, finalTargetDir)
	modifyArchiveJSON(*sourceDir, finalTargetDir)
	compressWasm(*sourceDir, finalTargetDir)

	fmt.Println(Green + "All file processing completed" + Reset)
}

// copyTemplate copies template files to the target directory
func copyTemplate(targetDir string) {
	// Define the list of files to copy
	files := []string{
		"game.json",
		"project.config.json",
		"project.private.config.json",
		"dmloader.js",
		"game.js",
	}

	// Copy individual files
	for _, file := range files {
		data, err := template.Files.ReadFile(file)
		if err != nil {
			fmt.Printf(Red+"Error reading embedded file %s: %v\n"+Reset, file, err)
			continue
		}

		// If the file is project.config.json or project.private.config.json, replace the {projectName} placeholder
		data = []byte(strings.ReplaceAll(string(data), "{projectName}", projectName))
		data = []byte(strings.ReplaceAll(string(data), "{wxAppId}", *appId))

		targetPath := filepath.Join(targetDir, file)
		err = os.WriteFile(targetPath, data, 0755)
		if err != nil {
			fmt.Printf(Red+"Error writing file %s: %v\n"+Reset, targetPath, err)
			continue
		}

		fmt.Printf(Green+"Successfully copied file: %s -> %s\n"+Reset, file, targetPath)
	}

	// Copy the weapp-adapter directory and its contents
	copyDir("weapp-adapter", filepath.Join(targetDir, "weapp-adapter"))
}

// copyDir recursively copies a directory
func copyDir(srcDir, destDir string) {
	entries, err := template.Files.ReadDir(srcDir)
	if err != nil {
		fmt.Printf(Red+"Error reading embedded directory %s: %v\n"+Reset, srcDir, err)
		return
	}

	// Create target directory
	if err := ensureDir(destDir); err != nil {
		fmt.Printf(Red+"Error creating directory %s: %v\n"+Reset, destDir, err)
		return
	}

	for _, entry := range entries {
		srcPath := fmt.Sprintf("%s/%s", srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			copyDir(srcPath, destPath)
		} else {
			// Copy files
			data, err := template.Files.ReadFile(srcPath)
			if err != nil {
				fmt.Printf(Red+"Error reading embedded file %s: %v\n"+Reset, srcPath, err)
				continue
			}

			err = os.WriteFile(destPath, data, 0755)
			if err != nil {
				fmt.Printf(Red+"Error writing file %s: %v\n"+Reset, destPath, err)
				continue
			}

			fmt.Printf(Green+"Successfully copied file: %s -> %s\n"+Reset, srcPath, destPath)
		}
	}
}

// modifyWasmJS processes _wasm.js files, replacing global variables with GameGlobal
func modifyWasmJS(sourceDir, targetDir string) {
	// Find all _wasm.js files in the source directory
	pattern := filepath.Join(sourceDir, "*_wasm.js")
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf(Red+"Error finding _wasm.js files: %v\n"+Reset, err)
		return
	}

	// Process each _wasm.js file
	for _, file := range files {
		// Read file content
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf(Red+"Error reading file %s: %v\n"+Reset, file, err)
			continue
		}

		// Perform text replacements
		newContent := string(content)
		newContent = strings.ReplaceAll(newContent, "var DMSYS", "GameGlobal.DMSYS")
		newContent = strings.ReplaceAll(newContent, "var FS=", "GameGlobal.FS=")
		newContent = strings.ReplaceAll(newContent, "var Module=typeof Module!=\"undefined\"?Module:{};", "var Module=GameGlobal.Module;")

		// Get filename
		filename := filepath.Base(file)
		// Write to target directory
		targetPath := filepath.Join(targetDir, filename)
		err = os.WriteFile(targetPath, []byte(newContent), 0755)
		if err != nil {
			fmt.Printf(Red+"Error writing file %s: %v\n"+Reset, targetPath, err)
			continue
		}

		fmt.Printf(Green+"Successfully processed file: %s -> %s\n"+Reset, file, targetPath)
	}
}

// compressWasm compresses .wasm files using brotli compression
func compressWasm(sourceDir, targetDir string) {
	// Read all files in the source directory
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf(Red+"Error reading source directory: %v\n"+Reset, err)
		return
	}

	// Iterate through all files
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Only process .wasm files
		if filepath.Ext(file.Name()) == ".wasm" {
			// Construct source file path
			sourceFilePath := filepath.Join(sourceDir, file.Name())

			// Read .wasm file content
			data, err := os.ReadFile(sourceFilePath)
			if err != nil {
				fmt.Printf(Red+"Error reading file %s: %v\n"+Reset, sourceFilePath, err)
				continue
			}

			// Compress using brotli
			compressedData, err := compressBrotli(data, 11) // Use high compression quality
			if err != nil {
				fmt.Printf(Red+"Error compressing file %s: %v\n"+Reset, sourceFilePath, err)
				continue
			}

			// Create target file path (add .br suffix)
			targetFileName := file.Name() + ".br"
			targetFilePath := filepath.Join(targetDir, targetFileName)

			// Write compressed file to target directory
			err = os.WriteFile(targetFilePath, compressedData, 0755)
			if err != nil {
				fmt.Printf(Red+"Error writing compressed file %s: %v\n"+Reset, targetFilePath, err)
				continue
			}

			fmt.Printf(Green+"Successfully compressed file: %s -> %s\n"+Reset, sourceFilePath, targetFilePath)
		}
	}
}

// renameArchiveFiles copies archive files to target directory, adding .bin suffix to non-JSON files
func renameArchiveFiles(sourceDir, targetDir string) {
	// Construct source and target archive directory paths
	sourceArchiveDir := filepath.Join(sourceDir, "archive")
	targetArchiveDir := filepath.Join(targetDir, "archive")

	// Ensure source and target archive directories exist
	if err := ensureDir(sourceArchiveDir); err != nil {
		fmt.Printf(Red+"Error creating source archive directory: %v\n"+Reset, err)
		return
	}

	if err := ensureDir(targetArchiveDir); err != nil {
		fmt.Printf(Red+"Error creating target archive directory: %v\n"+Reset, err)
		return
	}

	// Read all files in the source archive directory
	files, err := os.ReadDir(sourceArchiveDir)
	if err != nil {
		fmt.Printf(Red+"Error reading source archive directory: %v\n"+Reset, err)
		return
	}

	// Process each file
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Get filename and extension
		filename := file.Name()
		ext := filepath.Ext(filename)

		// Construct source file path
		sourcePath := filepath.Join(sourceArchiveDir, filename)

		// If it's a JSON file, copy directly
		if ext == ".json" {
			targetPath := filepath.Join(targetArchiveDir, filename)
			if err := copyGameAssetFiles(sourcePath, targetPath); err != nil {
				fmt.Printf(Red+"Error copying file %s: %v\n"+Reset, filename, err)
				continue
			}
			fmt.Printf(Green+"Successfully copied file: %s -> %s\n"+Reset, filename, targetPath)
		} else {
			// If it's not a JSON file, add .bin suffix and copy
			targetPath := filepath.Join(targetArchiveDir, filename+".bin")
			if err := copyGameAssetFiles(sourcePath, targetPath); err != nil {
				fmt.Printf(Red+"Error copying file %s: %v\n"+Reset, filename, err)
				continue
			}
			fmt.Printf(Green+"Successfully copied and renamed file: %s -> %s\n"+Reset, filename, filename+".bin")
		}
	}
}

// modifyArchiveJSON modifies archive JSON file to ensure all piece names have .bin suffix
func modifyArchiveJSON(sourceDir, targetDir string) {
	// Construct source and target archive directory paths
	sourceArchiveDir := filepath.Join(sourceDir, "archive")
	targetArchiveDir := filepath.Join(targetDir, "archive")
	sourceJSONFilePath := filepath.Join(sourceArchiveDir, "archive_files.json")
	targetJSONFilePath := filepath.Join(targetArchiveDir, "archive_files.json")

	// Ensure source and target archive directories exist
	if err := ensureDir(sourceArchiveDir); err != nil {
		fmt.Printf(Red+"Error creating source archive directory: %v\n"+Reset, err)
		return
	}

	if err := ensureDir(targetArchiveDir); err != nil {
		fmt.Printf(Red+"Error creating target archive directory: %v\n"+Reset, err)
		return
	}

	// Check if source file exists
	if _, err := os.Stat(sourceJSONFilePath); os.IsNotExist(err) {
		fmt.Printf(Yellow+"Source archive_files.json file does not exist: %s\n"+Reset, sourceJSONFilePath)
		return
	}

	// Read source JSON file
	data, err := os.ReadFile(sourceJSONFilePath)
	if err != nil {
		fmt.Printf(Red+"Error reading source archive_files.json file: %v\n"+Reset, err)
		return
	}

	// Parse JSON
	var archiveJSON ArchiveJSON
	err = json.Unmarshal(data, &archiveJSON)
	if err != nil {
		fmt.Printf(Red+"Error parsing source archive_files.json file: %v\n"+Reset, err)
		return
	}

	// Modify each name in pieces to ensure they have .bin suffix
	for i := range archiveJSON.Content {
		for j := range archiveJSON.Content[i].Pieces {
			name := archiveJSON.Content[i].Pieces[j].Name
			// If name doesn't end with .bin, add .bin suffix
			if !strings.HasSuffix(name, ".bin") {
				archiveJSON.Content[i].Pieces[j].Name = name + ".bin"
			}
		}
	}

	// Write modified data to target file
	modifiedData, err := json.MarshalIndent(archiveJSON, "", "    ")
	if err != nil {
		fmt.Printf(Red+"Error serializing modified data: %v\n"+Reset, err)
		return
	}

	err = os.WriteFile(targetJSONFilePath, modifiedData, 0755)
	if err != nil {
		fmt.Printf(Red+"Error writing target archive_files.json file: %v\n"+Reset, err)
		return
	}

	fmt.Printf(Green+"Successfully modified and wrote archive_files.json file: %s\n"+Reset, targetJSONFilePath)
}
