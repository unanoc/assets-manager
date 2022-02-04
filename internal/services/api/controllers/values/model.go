package values

type (
	TagValuesResponse struct {
		Tags []Tag `json:"tags"`
	}

	Tag struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
