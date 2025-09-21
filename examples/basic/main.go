package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pluggedin "github.com/veriteknik/pluggedinkit-go"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("PLUGGEDIN_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set PLUGGEDIN_API_KEY environment variable")
	}

	// Initialize client
	client := pluggedin.NewClient(apiKey)

	ctx := context.Background()

	// List documents
	fmt.Println("Listing documents...")
	docs, err := client.Documents.List(ctx, &pluggedin.DocumentFilters{
		Source: pluggedin.SourceAll,
		Limit:  5,
	})
	if err != nil {
		log.Printf("Error listing documents: %v", err)
	} else {
		fmt.Printf("Found %d documents:\n", docs.Total)
		for _, doc := range docs.Documents {
			fmt.Printf("  - %s (%d bytes)\n", doc.Title, doc.FileSize)
		}
	}

	// Perform a search
	fmt.Println("\nSearching documents...")
	results, err := client.Documents.Search(ctx, "API", nil, 5, 0)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
	} else {
		fmt.Printf("Found %d search results:\n", results.Total)
		for _, result := range results.Results {
			fmt.Printf("  - %s (relevance: %.2f)\n", result.Title, result.RelevanceScore)
		}
	}

	// Query RAG
	fmt.Println("\nQuerying knowledge base...")
	answer, err := client.RAG.AskQuestion(ctx, "What is Plugged.in?")
	if err != nil {
		log.Printf("Error querying RAG: %v", err)
	} else {
		fmt.Printf("Answer: %s\n", answer)
	}
}