package application

import (
	"fmt"
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type Application struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	debian   []base.FileInfo
	misc     []base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (a Application) String() string {
	return fmt.Sprintf("Application project")
}

func createApplicationDirtree(path string, options base.ProjectOptions) error {
	var subdirs []string
	var prefix string

	if options.PackageProject {
		prefix = options.ProjectName
		subdirs = append(subdirs, "pkg_install/misc")
		subdirs = append(subdirs, "pkg_install/debian")
	}

	subdirs = append(subdirs, prefix+"/src")
	subdirs = append(subdirs, prefix+"/include")

	for _, dir := range subdirs {
		err := os.MkdirAll(path+"/"+dir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a Application) Build() error {
	// create root path and subdirs
	if err := createApplicationDirtree(a.rootPath, a.ProjectOptions); err != nil {
		return err
	}

	// create sources
	for _, f := range a.sources {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range a.headers {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create debian scripts
	for _, f := range a.debian {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create Makefile (future CMakeLists.txt)

	return nil
}

func createSources(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo

	sources := []string{
		"main",
	}

	for _, s := range sources {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/src/"+s, ".c"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewSource(fileOptions),
		})
	}

	return files
}

func createHeaders(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo
	var headers []string

	for _, suffix := range []string{"_def", "_prt", "_struct"} {
		headers = append(headers, options.ProjectName+suffix)
	}

	headers = append(headers, options.ProjectName)

	for _, h := range headers {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/include/"+h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, ""),
		})
	}

	return files
}

func New(options base.ProjectOptions) (base.Project, error) {
	var rootPath string
	var prefix string
	cwd, err := os.Getwd()

	if err != nil {
		return &Application{}, err
	}

	if options.PackageProject {
		prefix = options.ProjectName
		rootPath = cwd + "/package-" + options.ProjectName
	} else {
		rootPath = cwd + "/" + options.ProjectName
	}

	application := &Application{
		rootPath:       rootPath,
		sources:        createSources(options, rootPath, prefix),
		headers:        createHeaders(options, rootPath, prefix),
		debian:         common.CreateDebianScripts(options, rootPath, prefix),
		ProjectOptions: options,
	}

	return application, nil
}
