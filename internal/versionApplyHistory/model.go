package versionApplyHistory

type CreateVersionApplyHistory struct {
	ID             string `json:"id"`
	VersionId      string `json:"version_id"`
	ApplyDate      string `json:"apply_date"`
	IsValid        string `json:"is_valid"`
	IsApplySuccess string `json:"is_apply_success"`
	Errors         string `json:"errors"`
}
