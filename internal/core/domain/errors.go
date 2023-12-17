package domain

type AssetNotFoundError struct {
}

func (e *AssetNotFoundError) Error() string {
	return "Asset not found"
}
