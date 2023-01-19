package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getResp = func(client gpt3.Client, ctx context.Context, question string) {

	err := client.CompletionStreamWithEngine(
		ctx,
		gpt3.TextDavinci001Engine,

		gpt3.CompletionRequest{
			Prompt:      []string{question},
			MaxTokens:   gpt3.IntPtr(300),
			Temperature: gpt3.Float32Ptr(0),
		},

		func(cr *gpt3.CompletionResponse) {
			fmt.Printf(cr.Choices[0].Text)
		},
	)
	fmt.Println()

	if err != nil {
		panic(err)
	}
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	apiKey := viper.GetString("API_KEY")

	ctx := context.Background()

	client := gpt3.NewClient(apiKey)

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Connect with ChatGPT through console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)

			quit := false

			for !quit {
				fmt.Print("Ask me something ('quit' to end): ")
				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				switch question {
				case "quit":
					quit = true

				default:
					getResp(client, ctx, question)
				}
			}
		},
	}

	rootCmd.Execute()
}
