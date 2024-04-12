package convert

import "github.com/chaewonkong/msa-link-api/link"

// MapToUpdatePayload converts a map of Open Graph tags to a link.UpdatePayload
func MapToUpdatePayload(ogtags map[string]string) link.UpdatePayload {
	updatePayload := link.UpdatePayload{}

	if img, exists := ogtags["og:image"]; exists {
		updatePayload.ThumbnailImage = img
	}

	if title, exists := ogtags["og:title"]; exists {
		updatePayload.Title = title
	}

	if description, exists := ogtags["og:description"]; exists {
		updatePayload.Description = description
	}

	return updatePayload
}
