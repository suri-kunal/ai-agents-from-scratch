# AI Agents from Scratch

## Vision

Build robust, state-of-the-art AI agents from first principles—no agent frameworks, no shortcuts, no crutches. Every agent, every component, every evaluation must be built and understood from the ground up.

**Goal:** Master every step of agent development, from data to deployment, with ruthless honesty and measurable progress.

## What This Project Is (and Isn't)

-   **Is:** A research-grade codebase for building, evaluating, and iterating on AI agents for real-world tasks, using only foundational libraries (LLM APIs, PyTorch, HuggingFace, etc.).
-   **Isn't:** A playground for copying agent frameworks, or a collection of "prompt hacks." No LangChain, CrewAI, AutoGen, smol-agents, or similar abstractions allowed.

## Core Principles

-   **No Frameworks:** All agent logic, memory, planning, and tool use must be implemented from scratch.
-   **Benchmark-Driven:** Every agent must be evaluated on a public, competitive benchmark. No cherry-picking, no hand-waving.
-   **Brutal Honesty:** Every failure, limitation, and blind spot must be documented and analyzed. No self-delusion.
-   **Iterative Mastery:** Each agent is a stepping stone. Master the basics before moving to more complex tasks.

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

1.  **Tool-Use Agent v1 (Single Tool Selection):**
    -   **Task:** Given a scenario, select the correct tool from a list.
    -   **Benchmark:** [Agent Leaderboard Tool Use/Function Calling Dataset](https://huggingface.co/blog/pratikbhavsar/agent-leaderboard)
2.  **Tool-Use Agent v2 (Parallel Tool Use & Parameterization):**
    -   **Task:** Given a scenario, select one or more tools and correctly populate their arguments.
    -   **Benchmark:** To be determined.
3.  **Planning Agent:**
    -   **Task:** Decompose a multi-step problem into a sequence of tool calls.
    -   **Benchmark:** To be determined.

*(Each stage must be benchmarked, analyzed, and documented before moving to the next.)*

## How to Contribute

-   Fork the repo, create a feature branch, and submit a pull request.
-   All new agents must include:
    -   A clear task definition and benchmark.
    -   A baseline implementation.
    -   Automated evaluation scripts.
    -   Error analysis and documentation in a local `README.md`.

## License

This project is released into the public domain under [The Unlicense](https://unlicense.org).
