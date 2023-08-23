package mainpackage

import (
	"fmt"
	"strings"
)

import (
	"github.com/your-repo/utils" // replace with your actual repo path
)

type Document struct {
	PageContent string
	Metadata    map[string]interface{}
}

func DocToPoints(docs []Document, chatTool *utils.ChatGptTool) ([]map[string]interface{}, error) {
	points := make([]map[string]interface{}, len(docs))

	// Iterate over the docs slice
	for i, doc := range docs {
		// Combine the PageContent and Metadata into a single string
		metadataStrs := make([]string, 0, len(doc.Metadata))
		for k, v := range doc.Metadata {
			metadataStrs = append(metadataStrs, fmt.Sprintf("%s: %v", k, v))
		}
		fullText := fmt.Sprintf("%s\nMetadata:\n%s", doc.PageContent, strings.Join(metadataStrs, "\n"))

		// Get the embedding of the combined text
		embeddingResponse, err := chatTool.GetEmbedding(fullText, "text-davinci-002")
		if err != nil {
			return nil, fmt.Errorf("Failed to get embedding for document %d: %v", i, err)
		}
		embedding := embeddingResponse.Data[0].Embedding

		// Prepare the point data
		points[i] = map[string]interface{}{
			"id": i + 1,
			"payload": map[string]interface{}{
				"text": fullText,
			},
			"vector": embedding,
		}
	}

	return points, nil
}
