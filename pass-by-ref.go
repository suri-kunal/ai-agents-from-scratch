package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/codes"
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

var report = `
Q3 Performance Summary:
Our customer satisfaction score rose to 92 points this quarter.
Revenue grew by 45% compared to last year.
Market share is now at 23% in our primary market.
Customer churn decreased to 5% from 8%.
New user acquisition cost is $43 per user.
Product adoption rate increased to 78%.
Employee satisfaction is at 87 points.
Operating margin improved to 34%.
`

func llmApi(ctx context.Context, client *openai.Client, inputMessage string, model shared.ChatModel) (*responses.Response, error) {
	tracer := otel.GetTracerProvider().Tracer("pass-by-ref/llm")
	childCtx, childSpan := tracer.Start(ctx, "Calling LLM API")
	defer childSpan.End()

	childSpan.SetAttributes(
		attribute.String("gen_ai.system","openai"),
		attribute.String("gen_ai.request.model",string(model)),
	)

	resp, err := client.Responses.New(
		childCtx,
		responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(inputMessage)},
			Model: model,
		},
	)

	if err != nil {
		childSpan.SetAttributes(
			attribute.Int64("gen_ai.usage.output_tokens",0),
			attribute.Int64("gen_ai.usage.input_tokens",0),
		)
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error,err.Error())
		return nil, err
	}

	childSpan.SetAttributes(
		attribute.Int64("gen_ai.usage.output_tokens",resp.Usage.OutputTokens),
		attribute.Int64("gen_ai.usage.input_tokens",resp.Usage.InputTokens),
	)
	childSpan.SetStatus(codes.Ok,"")

	return resp, err
}

func promptChain(ctx context.Context, client *openai.Client, input string, prompts []string) {
	tracer := otel.GetTracerProvider().Tracer("pass-by-ref/prompt-chain")
	childCtx, parentSpan := tracer.Start(ctx, "Prompt Chain Begins")
	defer parentSpan.End()

	parentSpan.SetAttributes(
		attribute.String("prompt_chain.input",input),
		attribute.Int("prompt_chain.step_count",len (prompts)),
	)
	
	result := input
	for idx, prompt := range prompts {
		stepCtx, childSpan := tracer.Start(childCtx, fmt.Sprintf("Running Step %d", idx + 1))
		inputMessage := prompt + "\nInput: " + result
		childSpan.SetAttributes(
			attribute.String("step.prompt",prompt),
			attribute.String("step.input_message",inputMessage),
			attribute.String("step.model",openai.ChatModelGPT5Mini),
		)
		resp, err := llmApi(stepCtx, client, inputMessage, openai.ChatModelGPT5Mini)
		if err != nil {
			childSpan.RecordError(err)
			childSpan.SetStatus(codes.Error,err.Error())
			childSpan.End()
			parentSpan.RecordError(err)
			parentSpan.SetStatus(codes.Error,fmt.Sprintf("Chain failed at %d", idx + 1))
			return
		}
		result = resp.OutputText()
		childSpan.SetAttributes(
			attribute.String("step.output",result),
		)
		childSpan.SetStatus(codes.Ok,"")
		childSpan.End()
		childCtx = stepCtx
	}
	parentSpan.SetAttributes(
		attribute.String("prompt_chain.final_output",result),
	)
	parentSpan.SetStatus(codes.Ok,"")
}

func main() {
	fmt.Println("Pass by reference")
	
	ctx := context.Background()

	// 1. Create an exporter
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create a resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("PromptChainPatternPassByRef"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironmentName("staging"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	// 3. Create a TraceProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Why is this good? Conpare it with individual-copy.go
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	client := openai.NewClient()
	
	promptChain(ctx,&client,report,dataProcessingSteps)

}