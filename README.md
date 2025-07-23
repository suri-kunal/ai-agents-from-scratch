# AI Agents from Scratch

## Vision

I have been awestruck by capabilities of LLMs/VLMs and how they can perform tasks with agency. Dabbled with a lot of LLM frameworks (Langchain, LangGraph, and whatnot) but somehow felt something missing. One day I came across [From Deep Learning Foundations to Stable Diffusion](https://course.fast.ai/Lessons/part2.html) by Jeremy Howard and his methodology for this course hit a nerve - Don't use an API until you cannot recreate it yourself.

Following principles of this course, my vision for this repo is to build robust, state-of-the-art AI agents from first principles — no agent frameworks, no shortcuts, no crutches. Every agent, every component, every evaluation must be built and understood from the ground up.

**Goal:** Master every step of agent development - from data to deployment - with measurable progress.

## Scope of this Project

- My vision for this project is purely pedagogical. I plan on learning how to build, evaluate, and iterate every part of AI agents - tool usage, memory, planning, and reasoning. Since my goal is learning with measurable progress, I would be using a lot of publicly available benchmarks. As part of this journey, I would be using only foundational libraries (Python, LLM APIs, PyTorch, HuggingFace, etc.) and would be actively avoiding agent frameworks like LangChain, CrewAI, AutoGen, smol-agents.

## Project Structure

-   `agents/` — Implementations of different agent architectures (tool-use, planning, etc.)
-   `data/` — Scripts and utilities for downloading, processing, and inspecting benchmark datasets.
-   `baselines/` — Trivial and non-trivial baselines for each task.
-   `eval/` — Evaluation scripts, metrics, and leaderboard tracking.
-   `docs/` — Design docs, research notes, and deep dives into agent failures and improvements.
-   `README.md` — This file.

## Allowed Technologies

-   **LLMs:** OpenAI, Gemini, or any public LLM API.
-   **ML Libraries:** PyTorch, HuggingFace Transformers.
-   **Infra:** Any standard software for queues, caches, databases, logging, etc.
-   **Prohibited:** Any agent framework or abstraction that handles agent logic for you.

## Roadmap

1.  **Tool-Use Agent:**
    -   **Task:** Given a scenario, select the correct tool from a given list.
    -   **Benchmark:** [Agent Leaderboard Tool Use/Function Calling Dataset](https://huggingface.co/blog/pratikbhavsar/agent-leaderboard)

*(Each stage must be benchmarked, analyzed, and documented before moving to the next - I really cannot emphasize this enough)*

## License

This project is released into the public domain under [The Unlicense](https://unlicense.org).
