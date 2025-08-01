package config

import (
	"errors"
	"os"
	"path/filepath"

	rt "github.com/zyedidia/micro/v2/runtime"
)

const (
	RTSyntax       = 0
	RTHelp         = 1
	RTSyntaxHeader = 2
	RTPlugin       = 3 // Stub for tests - plugins removed
)

var (
	NumTypes = 4 // How many filetypes are there (including RTPlugin stub for tests)
)

type RTFiletype int

// RuntimeFile allows the program to read runtime data like colorschemes or syntax files
type RuntimeFile interface {
	// Name returns a name of the file without paths or extensions
	Name() string
	// Data returns the content of the file.
	Data() ([]byte, error)
}

// allFiles contains all available files, mapped by filetype
var allFiles [][]RuntimeFile
var realFiles [][]RuntimeFile

func init() {
	initRuntimeVars()
}

func initRuntimeVars() {
	allFiles = make([][]RuntimeFile, NumTypes)
	realFiles = make([][]RuntimeFile, NumTypes)
}

// NewRTFiletype creates a new RTFiletype
func NewRTFiletype() int {
	NumTypes++
	allFiles = append(allFiles, []RuntimeFile{})
	realFiles = append(realFiles, []RuntimeFile{})
	return NumTypes - 1
}

// some file on filesystem
type realFile string

// some asset file
type assetFile string

// a file with the data stored in memory
type memoryFile struct {
	name string
	data []byte
}

func (mf memoryFile) Name() string {
	return mf.name
}
func (mf memoryFile) Data() ([]byte, error) {
	return mf.data, nil
}

func (rf realFile) Name() string {
	fn := filepath.Base(string(rf))
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

func (rf realFile) Data() ([]byte, error) {
	return os.ReadFile(string(rf))
}

func (af assetFile) Name() string {
	fn := filepath.Base(string(af))
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

func (af assetFile) Data() ([]byte, error) {
	return rt.Asset(string(af))
}

// AddRuntimeFile registers a file for the given filetype
func AddRuntimeFile(fileType RTFiletype, file RuntimeFile) {
	allFiles[fileType] = append(allFiles[fileType], file)
}

// AddRealRuntimeFile registers a file for the given filetype
func AddRealRuntimeFile(fileType RTFiletype, file RuntimeFile) {
	allFiles[fileType] = append(allFiles[fileType], file)
	realFiles[fileType] = append(realFiles[fileType], file)
}

// AddRuntimeFilesFromDirectory registers each file from the given directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromDirectory(fileType RTFiletype, directory, pattern string) {
	files, _ := os.ReadDir(directory)
	for _, f := range files {
		if ok, _ := filepath.Match(pattern, f.Name()); !f.IsDir() && ok {
			fullPath := filepath.Join(directory, f.Name())
			AddRealRuntimeFile(fileType, realFile(fullPath))
		}
	}
}

// AddRuntimeFilesFromAssets registers each file from the given asset-directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromAssets(fileType RTFiletype, directory, pattern string) {
	files, err := rt.AssetDir(directory)
	if err != nil {
		return
	}

assetLoop:
	for _, f := range files {
		if ok, _ := filepath.Match(pattern, f); ok {
			af := assetFile(filepath.Join(directory, f))
			for _, rf := range realFiles[fileType] {
				if af.Name() == rf.Name() {
					continue assetLoop
				}
			}
			AddRuntimeFile(fileType, af)
		}
	}
}

// FindRuntimeFile finds a runtime file of the given filetype and name
// will return nil if no file was found
func FindRuntimeFile(fileType RTFiletype, name string) RuntimeFile {
	for _, f := range ListRuntimeFiles(fileType) {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

// ListRuntimeFiles lists all known runtime files for the given filetype
func ListRuntimeFiles(fileType RTFiletype) []RuntimeFile {
	return allFiles[fileType]
}

// ListRealRuntimeFiles lists all real runtime files (on disk) for a filetype
// these runtime files will be ones defined by the user and loaded from the config directory
func ListRealRuntimeFiles(fileType RTFiletype) []RuntimeFile {
	return realFiles[fileType]
}

// InitRuntimeFiles initializes all assets files and the config directory.
// If `user` is false, InitRuntimeFiles ignores the config directory and
// initializes asset files only.
func InitRuntimeFiles(user bool) {
	add := func(fileType RTFiletype, dir, pattern string) {
		if user {
			AddRuntimeFilesFromDirectory(fileType, filepath.Join(ConfigDir, dir), pattern)
		}
		AddRuntimeFilesFromAssets(fileType, filepath.Join("runtime", dir), pattern)
	}

	initRuntimeVars()

	add(RTSyntax, "syntax", "*.yaml")
	add(RTSyntaxHeader, "syntax", "*.hdr")
	add(RTHelp, "help", "*.md")
}

// InitPlugins is a no-op in micromini since plugins are removed
func InitPlugins() {
	// Plugin system removed in micromini - no-op
}

// Removed unused plugin functions from micromini to reduce code size

// PluginCommand is a no-op stub in micromini since plugins are removed
func PluginCommand(args ...interface{}) error {
	return errors.New("Plugin system removed in micromini")
}

// LoadAllPlugins is a no-op stub in micromini since plugins are removed
func LoadAllPlugins() error {
	return errors.New("Plugin system removed in micromini")
}

// RunPluginFn is a no-op stub in micromini since plugins are removed
func RunPluginFn(name string, args ...interface{}) error {
	return errors.New("Plugin system removed in micromini")
}

// RunPluginFnBool is a no-op stub in micromini since plugins are removed
func RunPluginFnBool(settings interface{}, name string, args ...interface{}) (bool, error) {
	return true, errors.New("Plugin system removed in micromini")
}
