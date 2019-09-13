package repository

type Repository struct {
  Name string `json:"name"`
  FilePath string `json:"file_path"`
  Organization bool `json:"is_organization"`
  Owner string `json:"owner"`
  Provider string `json:"provider"`
  Id string `json:"id"`
}
