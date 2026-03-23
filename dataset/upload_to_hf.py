"""
Upload SWE-Bench 5G dataset to HuggingFace Hub.

Prerequisites:
    pip install datasets huggingface_hub
    huggingface-cli login
"""

from datasets import Dataset, Features, Value
import json

# Define schema
features = Features({
    "instance_id": Value("string"),
    "dataset_id": Value("string"),
    "task": Value("string"),
    "nf_type": Value("string"),
    "user": Value("string"),
    "repo": Value("string"),
    "language": Value("string"),
    "workdir": Value("string"),
    "image_url": Value("string"),
    "patch": Value("string"),
    "commit_id": Value("string"),
    "parent_commit": Value("string"),
    "problem_statement": Value("string"),
    "f2p_patch": Value("string"),
    "f2p_script": Value("string"),
    "FAIL_TO_PASS": Value("string"),
    "PASS_TO_PASS": Value("string"),
    "github": Value("string"),
    "pre_commands": Value("string"),
    "spec_reference": Value("string"),
    "difficulty": Value("string"),
    "affected_function": Value("string"),
    "lines_changed": Value("int32"),
    "files_changed": Value("int32"),
    "test_suite_num": Value("string"),
    "split": Value("string"),
})

# Load data
records = []
with open("swebench5g.jsonl", "r", encoding="utf-8") as f:
    for line in f:
        records.append(json.loads(line.strip()))

# Create dataset
ds = Dataset.from_list(records, features=features)

# Push to HuggingFace Hub
# Change "tianjiang/SWEBench5G" to your desired repo ID
REPO_ID = "tianjiang/SWEBench5G"

ds.push_to_hub(
    REPO_ID,
    split="test",
    commit_message="Initial release: 1 pilot task (PCF Issue #879)",
)

print(f"Dataset uploaded to https://huggingface.co/datasets/{REPO_ID}")
