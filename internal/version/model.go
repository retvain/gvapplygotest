package version

type Version struct {
	ID              int    `json:"id"`
	IsApplied       bool   `json:"is_applied"`
	VersionFilePath string `json:"version_file_path"`
}
