package web

import "gin-boilerplate/comm/swagger/gen"

func GenerageSwaggerDoc(searchDir string) {
	swagger, cfg := gen.New(), gen.Config{
		SearchDir:          searchDir,
		MainAPIFile:        "../main.go",
		PropNamingStrategy: "camelcase",
		MarkdownFilesDir:   "",
		OutputDir:          "./",
		ParseVendor:        true,
		ParseDependency:    true,
	}
	go swagger.Build(&cfg)
}
