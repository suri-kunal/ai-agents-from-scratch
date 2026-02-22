package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

var dataProcessingSteps = []string{
	`Extract only the numerical values and their associated metrics from the text.
    Format each as 'value: metric' on a new line.
    Example format:
    92: customer satisfaction
    45%: revenue growth`,
	`Convert all numerical values to percentages where possible.
    If not a percentage or points, convert to decimal (e.g., 92 points -> 92%).
    Keep one number per line.
    Example format:
    92%: customer satisfaction
    45%: revenue growth`,
	`Sort all lines in descending order by numerical value.
    Keep the format 'value: metric' on each line.
    Example:
    92%: customer satisfaction
    87%: employee satisfaction`,
	`Format the sorted data as a markdown table with columns:
    | Metric | Value |
    |:--|--:|
    | Customer Satisfaction | 92% |`,
}

var inputReport = `Q3 Performance Summary:
Our customer satisfaction score rose to 92 points this quarter.
Revenue grew by 45% compared to last year.
Market share is now at 23% in our primary market.
Customer churn decreased to 5% from 8%.
New user acquisition cost is $43 per user.
Product adoption rate increased to 78%.
Employee satisfaction is at 87 points.
Operating margin improved to 34%.`

func llmApi(ctx context.Context, client *openai.Client, inputMessage string, modelName string) (resp *responses.Response, err error) {
	resp, err = (*client).Responses.New(
		ctx,
		responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(inputMessage),
			},
			Model: modelName,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func promptChain(ctx context.Context, client *openai.Client, promptList []string, initialMessage string, modelName string) (output string, err error) {
	var resp *responses.Response
	output = initialMessage
	fmt.Printf("======Input Message======\n %s\n", initialMessage)
	for idx, prompt := range promptList {
		fmt.Printf("======Step %d Prompt======\n%s\n", idx+1, prompt)
		inputString := fmt.Sprintf("%s\nInput: %s", prompt, output)
		resp, err = llmApi(ctx, client, inputString, modelName)
		if err != nil {
			return "", fmt.Errorf("Step %d failed due to %w", idx+1, err)
		}
		output = resp.OutputText()
		fmt.Printf("======Step %d Output Message======\n%s\n", idx+1, output)
	}
	return output, nil
}

func main() {
	ctx := context.Background()
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		panic("OPENAI_API_KEY is not set")
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	output, err := promptChain(ctx, &client, dataProcessingSteps, inputReport, "gpt-5-mini-2025-08-07")
	if err != nil {
		fmt.Printf("Error occurred\n %v\n", err.Error())
	} else {
		fmt.Printf("Final Output: %s", output)
	}

}
