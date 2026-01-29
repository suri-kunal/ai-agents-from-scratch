package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
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

func llmApi(ctx context.Context, tracer trace.Tracer, inputMessage string, model shared.ChatModel, idx int) (*responses.Response, error) {

	childCtx, childSpan := tracer.Start(
		ctx, "Invoking LLM API",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer childSpan.End()

	childSpan.SetAttributes(
		attribute.String("gen_ai.system", "openai"),
		attribute.String("gen_ai.request.model", string(model)),
	)

	// if idx == 3 {
	// 	childSpan.SetAttributes(
	// 		attribute.Int("gen_ai.usage.input_tokens", 0),
	// 		attribute.Int("gen_ai.usage.output_tokens", 0),
	// 	)
	// 	simErr := fmt.Errorf("Simulated error at step %d", idx)
	// 	childSpan.RecordError(simErr)
	// 	childSpan.SetStatus(codes.Error, simErr.Error())
	// 	return nil, simErr
	// }

	client := openai.NewClient()
	resp, err := client.Responses.New(
		childCtx,
		responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(inputMessage)},
			Model: model,
		},
	)

	if err != nil {
		childSpan.SetAttributes(
			attribute.Int64("gen_ai.usage.input_tokens", 0),
			attribute.Int64("gen_ai.usage.output_tokens", 0),
		)
		childSpan.RecordError(err)
		childSpan.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	childSpan.SetAttributes(
		attribute.Int64("gen_ai.usage.input_tokens", resp.Usage.InputTokens),
		attribute.Int64("gen_ai.usage.output_tokens", resp.Usage.OutputTokens),
	)
	childSpan.SetStatus(codes.Ok, "")
	return resp, err
}

func promptChain(ctx context.Context, tracer trace.Tracer, input string, prompts []string) {
	childCtx, parentSpan := tracer.Start(ctx, "Prompt Chain Begins")
	defer parentSpan.End()

	parentSpan.SetAttributes(
		attribute.String("prompt_chain.input", input),
		attribute.Int("prompt_chain.step_count", len(prompts)),
	)

	result := input
	for idx, prompt := range prompts {
		stepCtx, childSpan := tracer.Start(childCtx, fmt.Sprintf("Step %d", idx))
		inputMessage := prompt + "\nInput: " + result

		// Adding attributes to describe the request
		childSpan.SetAttributes(
			attribute.String("step.input_message", inputMessage),
			attribute.String("step.prompt", prompt),
			attribute.String("step.model", openai.ChatModelGPT5Mini),
		)
		resp, err := llmApi(stepCtx, tracer, inputMessage, openai.ChatModelGPT5Mini, idx)
		if err != nil {
			childSpan.RecordError(err)
			childSpan.SetStatus(codes.Error, err.Error())
			childSpan.End()
			parentSpan.RecordError(err)
			parentSpan.SetStatus(codes.Error, fmt.Sprintf("Chain failed at %d", idx))
			return
		}
		result = resp.OutputText()
		childSpan.SetAttributes(
			attribute.String("step.output", result),
		)
		childSpan.SetStatus(codes.Ok, "")
		childSpan.End()
		childCtx = stepCtx
	}
	parentSpan.SetAttributes(
		attribute.String("prompt_chain.final_output", result),
	)
	parentSpan.SetStatus(codes.Ok, "")
}

func main() {
	fmt.Println("Individual copy")

	// Initialize the Background context
	ctx := context.Background()

	// Step 1. Create an exporter where telemetry data goes
	exporter, err := otlptracehttp.New( 
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("prompt-chaining-agent"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "dev"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: Create a traceprovider (the "Factory" that creates the tracers)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer tp.Shutdown(ctx)

	// Step 3: Register it globally so that any code can use it
	otel.SetTracerProvider(tp)

	// Step 4: Get a tracer
	tracer := otel.Tracer("individual-copy/main")

	// trace Prompt Chaining
	// promptChain(ctx,report,dataProcessingSteps)
	promptChain(ctx, tracer, report, dataProcessingSteps)
}