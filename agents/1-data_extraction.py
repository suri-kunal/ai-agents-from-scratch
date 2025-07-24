import marimo

__generated_with = "0.14.12"
app = marimo.App(width="full", auto_download=["ipynb"])


@app.cell
def _():
    import marimo as mo
    return (mo,)


@app.cell
def _():
    from datasets import load_dataset, load_from_disk
    import os
    from tqdm import tqdm
    import json
    return load_dataset, tqdm


@app.cell
def _(mo):
    mo.md(r"""# This notebook to meant to extract data from [galileo-ai/agent-leaderboard](https://huggingface.co/datasets/galileo-ai/agent-leaderboard)""")
    return


@app.cell
def _(load_dataset, tqdm):
    for dataset_name in tqdm(['BFCL_v3_irrelevance', 'BFCL_v3_multi_turn_base_multi_func_call', 'BFCL_v3_multi_turn_base_single_func_call', 'BFCL_v3_multi_turn_composite', 'BFCL_v3_multi_turn_long_context', 'BFCL_v3_multi_turn_miss_func', 'BFCL_v3_multi_turn_miss_param', 'tau_long_context', 'toolace_single_func_call_1', 'toolace_single_func_call_2', 'xlam_multiple_tool_multiple_call', 'xlam_multiple_tool_single_call', 'xlam_single_tool_multiple_call', 'xlam_single_tool_single_call', 'xlam_tool_miss']):
        ds = load_dataset("galileo-ai/agent-leaderboard",dataset_name,)
        ds.save_to_disk(f"../data/tool_usage/{dataset_name}.hf")
    return


if __name__ == "__main__":
    app.run()
