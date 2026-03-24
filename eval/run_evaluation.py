#!/usr/bin/env python3
"""
run_evaluation.py — SWE-Bench 5G Evaluation Harness

Evaluates AI coding agents on SWE-Bench 5G task instances.

Workflow per (agent, instance):
  1. Start Docker container from task image
  2. Inject agent with the problem statement
  3. Agent reads code, diagnoses bug, generates patch
  4. Extract the agent's patch (git diff)
  5. Run test suite (existing + fail-to-pass)
  6. Record results

Usage:
    python eval/run_evaluation.py --agent claude-code --model opus-4.6
    python eval/run_evaluation.py --agent claude-code --instances pcf_pr65 amf_pr179
    python eval/run_evaluation.py --dataset dataset/swebench5g.jsonl --agent all
"""

import argparse
import json
import os
import subprocess
import sys
import time
from dataclasses import dataclass, asdict
from datetime import datetime
from pathlib import Path
from typing import Optional


@dataclass
class EvalResult:
    """Result of evaluating one agent on one task instance."""
    instance_id: str
    agent: str
    model: str
    resolved: bool              # All FAIL_TO_PASS pass AND all PASS_TO_PASS still pass
    existing_tests_pass: bool   # PASS_TO_PASS tests still pass
    fail_tests_pass: bool       # FAIL_TO_PASS tests now pass
    patch: str                  # Agent's generated patch
    time_seconds: float         # Wall-clock time
    error: str                  # Error message if any
    timestamp: str              # ISO timestamp


def docker_exec(container_id: str, cmd: str, timeout: int = 600) -> tuple:
    """Execute a command inside a running Docker container.
    Returns (exit_code, stdout, stderr)."""
    try:
        result = subprocess.run(
            ["docker", "exec", container_id, "bash", "-c", cmd],
            capture_output=True, text=True, timeout=timeout
        )
        return result.returncode, result.stdout, result.stderr
    except subprocess.TimeoutExpired:
        return -1, "", "TIMEOUT"
    except Exception as e:
        return -1, "", str(e)


def start_container(image_name: str) -> str:
    """Start a Docker container and return its ID."""
    result = subprocess.run(
        ["docker", "run", "-d", "--rm", image_name, "sleep", "3600"],
        capture_output=True, text=True
    )
    if result.returncode != 0:
        raise RuntimeError(f"Failed to start container: {result.stderr}")
    return result.stdout.strip()


def stop_container(container_id: str):
    """Stop and remove a Docker container."""
    subprocess.run(["docker", "stop", container_id],
                   capture_output=True, timeout=30)


def extract_patch(container_id: str, workdir: str) -> str:
    """Extract the agent's changes as a unified diff."""
    code, stdout, stderr = docker_exec(container_id, f"cd {workdir} && git diff")
    return stdout if code == 0 else ""


def run_agent_in_container(
    container_id: str,
    agent: str,
    model: str,
    workdir: str,
    timeout: int = 1800
) -> tuple:
    """Run an AI agent inside the container to fix the bug.

    Returns (success: bool, error: str).
    Each agent type has its own invocation method.
    """

    # Read problem statement
    code, problem, _ = docker_exec(container_id, "cat /opt/task/problem_statement.md")
    if code != 0:
        return False, "Failed to read problem statement"

    # Agent-specific invocation
    if agent == "claude-code":
        return _run_claude_code(container_id, model, workdir, problem, timeout)
    elif agent == "aider":
        return _run_aider(container_id, model, workdir, problem, timeout)
    elif agent == "codex-cli":
        return _run_codex_cli(container_id, model, workdir, problem, timeout)
    elif agent == "qwen":
        return _run_qwen(container_id, model, workdir, problem, timeout)
    elif agent == "manual-patch":
        # For testing: apply a pre-prepared patch
        return _run_manual_patch(container_id, workdir)
    else:
        return False, f"Unknown agent: {agent}"


def _run_claude_code(container_id, model, workdir, problem, timeout):
    """Run Claude Code CLI inside the container."""
    prompt = (
        f"You are inside a Docker container with a Go project at {workdir}. "
        f"There is a bug described below. Fix the bug by editing the source code. "
        f"Do NOT modify test files.\n\n{problem}"
    )
    # Claude Code CLI invocation
    cmd = (
        f'cd {workdir} && claude --model {model} --yes '
        f'--prompt "{prompt.replace(chr(34), chr(92)+chr(34))}"'
    )
    code, stdout, stderr = docker_exec(container_id, cmd, timeout=timeout)
    if code == 0:
        return True, ""
    else:
        return False, stderr[:500]


def _run_aider(container_id, model, workdir, problem, timeout):
    """Run Aider inside the container."""
    cmd = (
        f'cd {workdir} && echo "{problem}" | '
        f'aider --model {model} --yes --no-git'
    )
    code, stdout, stderr = docker_exec(container_id, cmd, timeout=timeout)
    return code == 0, stderr[:500] if code != 0 else ""


def _run_codex_cli(container_id, model, workdir, problem, timeout):
    """Run OpenAI Codex CLI inside the container."""
    cmd = (
        f'cd {workdir} && codex --model {model} '
        f'--approval-mode full-auto "{problem[:500]}"'
    )
    code, stdout, stderr = docker_exec(container_id, cmd, timeout=timeout)
    return code == 0, stderr[:500] if code != 0 else ""


def _run_qwen(container_id, model, workdir, problem, timeout):
    """Run Qwen model via DashScope OpenAI-compatible API to fix the bug.

    Workflow:
    1. Read the buggy source files from the container
    2. Send problem + source to Qwen API
    3. Parse the response for code changes
    4. Apply changes back to the container

    Env vars:
        DASHSCOPE_API_KEY: DashScope API key (required)
        DASHSCOPE_BASE_URL: API base URL (default: https://dashscope.aliyuncs.com/compatible-mode/v1)
    """
    import openai

    api_key = os.environ.get("DASHSCOPE_API_KEY")
    if not api_key:
        return False, "DASHSCOPE_API_KEY not set"

    base_url = os.environ.get(
        "DASHSCOPE_BASE_URL",
        "https://dashscope.aliyuncs.com/compatible-mode/v1"
    )

    client = openai.OpenAI(api_key=api_key, base_url=base_url)

    # Step 1: List Go source files in the project
    code, file_list, _ = docker_exec(
        container_id,
        f"find {workdir} -name '*.go' -not -path '*/vendor/*' -not -name '*_test.go' | head -50",
        timeout=30
    )
    go_files = [f.strip() for f in file_list.strip().split('\n') if f.strip()]

    # Step 2: Read the most relevant source files (from problem statement hints)
    source_contents = {}
    for fpath in go_files[:20]:  # limit to avoid token overflow
        rc, content, _ = docker_exec(container_id, f"cat {fpath}", timeout=10)
        if rc == 0 and len(content) < 15000:  # skip huge files
            rel_path = fpath.replace(workdir + "/", "")
            source_contents[rel_path] = content

    # Step 3: Build prompt
    files_text = ""
    for path, content in source_contents.items():
        files_text += f"\n--- {path} ---\n{content}\n"

    system_prompt = (
        "You are an expert Go developer fixing bugs in the free5GC 5G core network. "
        "Read the bug description and source code, then output ONLY the corrected file(s). "
        "For each file you modify, output in this exact format:\n\n"
        "=== FILE: <relative/path/to/file.go> ===\n"
        "<complete file content>\n"
        "=== END FILE ===\n\n"
        "Do NOT modify test files. Only fix the bug described."
    )

    user_prompt = f"## Bug Description\n\n{problem}\n\n## Source Files\n{files_text}"

    # Step 4: Call Qwen API
    # Qwen3.5-Flash uses thinking mode by default. We disable it for
    # structured output, or handle both content and reasoning_content.
    try:
        response = client.chat.completions.create(
            model=model or "qwen3.5-flash",
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_prompt[:60000]},  # token limit safety
            ],
            temperature=0.7,
            max_tokens=16000,
            extra_body={"enable_thinking": False},  # disable thinking for direct output
        )
        reply = response.choices[0].message.content or ""
        # If content is empty, try reasoning_content (thinking mode fallback)
        if not reply.strip():
            reasoning = getattr(response.choices[0].message, 'reasoning_content', '') or ''
            if reasoning:
                reply = reasoning
    except Exception as e:
        return False, f"Qwen API error: {str(e)[:300]}"

    if not reply.strip():
        return False, "Qwen returned empty response"

    # Step 5: Parse response and apply file changes
    import re
    file_pattern = r'=== FILE: (.+?) ===\n(.*?)\n=== END FILE ==='
    matches = re.findall(file_pattern, reply, re.DOTALL)

    if not matches:
        # Try alternative format: ```go blocks with filename comments
        alt_pattern = r'(?:// File: |// )(.+\.go)\n```go\n(.*?)```'
        matches = re.findall(alt_pattern, reply, re.DOTALL)

    if not matches:
        return False, f"Could not parse file changes from Qwen response (length={len(reply)})"

    applied = 0
    for rel_path, content in matches:
        rel_path = rel_path.strip()
        full_path = f"{workdir}/{rel_path}"
        # Write file via base64 to avoid heredoc escaping issues
        import base64
        b64 = base64.b64encode(content.encode('utf-8')).decode('ascii')
        rc, _, err = docker_exec(
            container_id,
            f"echo '{b64}' | base64 -d > {full_path}",
            timeout=30
        )
        if rc == 0:
            applied += 1

    if applied == 0:
        return False, "No files successfully written to container"

    return True, ""


def _run_manual_patch(container_id, workdir):
    """Apply a pre-prepared patch for testing the harness."""
    code, _, stderr = docker_exec(
        container_id,
        f"cd {workdir} && git apply /opt/task/fix.patch 2>/dev/null || true"
    )
    return True, ""


def run_tests_in_container(container_id: str, mode: str = "all") -> tuple:
    """Run test suite in container. Returns (exit_code, stdout)."""
    code, stdout, stderr = docker_exec(
        container_id,
        f"/opt/test-suite/run_tests.sh {mode}",
        timeout=300
    )
    return code, stdout + stderr


def evaluate_instance(
    instance: dict,
    agent: str,
    model: str,
    timeout: int = 1800,
    run_id: int = 1
) -> EvalResult:
    """Evaluate a single agent on a single task instance."""

    instance_id = instance["instance_id"]
    image_name = instance.get("image_url", f"swebench5g/free5gc:{instance_id}")
    workdir = instance.get("workdir", f"/opt/free5gc-pcf")

    print(f"\n{'='*60}")
    print(f"Evaluating: {instance_id} | Agent: {agent} | Model: {model} | Run: {run_id}")
    print(f"{'='*60}")

    container_id = None
    start_time = time.time()

    try:
        # 1. Start container
        print("[1/5] Starting container...")
        container_id = start_container(image_name)
        print(f"  Container: {container_id[:12]}")

        # 2. Verify environment (existing tests should pass)
        print("[2/5] Verifying environment...")
        code, output = run_tests_in_container(container_id, "existing")
        if code != 0:
            return EvalResult(
                instance_id=instance_id, agent=agent, model=model,
                resolved=False, existing_tests_pass=False, fail_tests_pass=False,
                patch="", time_seconds=time.time()-start_time,
                error="Environment verification failed",
                timestamp=datetime.now().isoformat()
            )

        # 3. Run agent
        print(f"[3/5] Running {agent} ({model})...")
        agent_ok, agent_error = run_agent_in_container(
            container_id, agent, model, workdir, timeout
        )

        # 4. Extract patch
        print("[4/5] Extracting patch...")
        patch = extract_patch(container_id, workdir)
        if not patch.strip():
            print("  WARNING: Agent produced no changes")

        # 5. Run all tests
        print("[5/5] Running test suite...")
        # Run existing tests
        e_code, e_out = run_tests_in_container(container_id, "existing")
        existing_pass = (e_code == 0)

        # Run fail-to-pass tests
        f_code, f_out = run_tests_in_container(container_id, "fail")
        fail_pass = (f_code == 0)

        resolved = existing_pass and fail_pass
        elapsed = time.time() - start_time

        status = "RESOLVED" if resolved else "NOT RESOLVED"
        print(f"\n  Result: {status} ({elapsed:.1f}s)")
        print(f"  Existing: {'PASS' if existing_pass else 'FAIL'}")
        print(f"  Fail-to-pass: {'PASS' if fail_pass else 'FAIL'}")

        return EvalResult(
            instance_id=instance_id, agent=agent, model=model,
            resolved=resolved, existing_tests_pass=existing_pass,
            fail_tests_pass=fail_pass, patch=patch,
            time_seconds=elapsed, error=agent_error,
            timestamp=datetime.now().isoformat()
        )

    except Exception as e:
        return EvalResult(
            instance_id=instance_id, agent=agent, model=model,
            resolved=False, existing_tests_pass=False, fail_tests_pass=False,
            patch="", time_seconds=time.time()-start_time,
            error=str(e), timestamp=datetime.now().isoformat()
        )

    finally:
        if container_id:
            print("  Cleaning up container...")
            stop_container(container_id)


def generate_report(results: list, output_dir: str):
    """Generate evaluation report."""
    os.makedirs(output_dir, exist_ok=True)

    # Save raw results
    results_path = os.path.join(output_dir, "results.json")
    with open(results_path, "w") as f:
        json.dump([asdict(r) for r in results], f, indent=2)

    # Generate summary
    total = len(results)
    resolved = sum(1 for r in results if r.resolved)
    resolve_rate = resolved / total * 100 if total > 0 else 0

    # By agent
    agent_stats = {}
    for r in results:
        key = f"{r.agent}+{r.model}"
        if key not in agent_stats:
            agent_stats[key] = {"total": 0, "resolved": 0, "times": []}
        agent_stats[key]["total"] += 1
        if r.resolved:
            agent_stats[key]["resolved"] += 1
        agent_stats[key]["times"].append(r.time_seconds)

    # Print summary
    print(f"\n{'='*60}")
    print(f" EVALUATION REPORT")
    print(f"{'='*60}")
    print(f" Total evaluations:  {total}")
    print(f" Resolved:           {resolved}/{total} ({resolve_rate:.1f}%)")
    print(f"\n By agent+model:")
    print(f" {'Configuration':<30} {'Resolved':>10} {'Rate':>8} {'Avg Time':>10}")
    print(f" {'-'*30} {'-'*10} {'-'*8} {'-'*10}")
    for key, stats in sorted(agent_stats.items()):
        rate = stats["resolved"]/stats["total"]*100
        avg_time = sum(stats["times"])/len(stats["times"])
        print(f" {key:<30} {stats['resolved']}/{stats['total']:>7} {rate:>7.1f}% {avg_time:>9.1f}s")

    # By instance
    print(f"\n By instance:")
    print(f" {'Instance':<30} {'Agent+Model':<25} {'Result':>10} {'Time':>8}")
    print(f" {'-'*30} {'-'*25} {'-'*10} {'-'*8}")
    for r in results:
        status = "RESOLVED" if r.resolved else "FAILED"
        print(f" {r.instance_id:<30} {r.agent}+{r.model:<14} {status:>10} {r.time_seconds:>7.1f}s")

    summary_path = os.path.join(output_dir, "summary.txt")
    # Also save as file
    with open(summary_path, "w") as f:
        f.write(f"SWE-Bench 5G Evaluation Report\n")
        f.write(f"Generated: {datetime.now().isoformat()}\n")
        f.write(f"Total: {total} | Resolved: {resolved} ({resolve_rate:.1f}%)\n")
        for key, stats in sorted(agent_stats.items()):
            rate = stats["resolved"]/stats["total"]*100
            f.write(f"  {key}: {stats['resolved']}/{stats['total']} ({rate:.1f}%)\n")

    print(f"\n Results saved to: {output_dir}/")
    print(f"   results.json  — raw results")
    print(f"   summary.txt   — summary report")


def main():
    parser = argparse.ArgumentParser(description="SWE-Bench 5G Evaluation Harness")
    parser.add_argument("--dataset", default="dataset/swebench5g.jsonl",
                        help="Path to dataset JSONL file")
    parser.add_argument("--agent", required=True,
                        choices=["claude-code", "aider", "codex-cli", "qwen", "manual-patch", "all"],
                        help="Agent to evaluate")
    parser.add_argument("--model", default="opus-4.6",
                        help="Model to use with the agent")
    parser.add_argument("--instances", nargs="*", default=None,
                        help="Specific instance IDs to evaluate (default: all)")
    parser.add_argument("--runs", type=int, default=1,
                        help="Number of independent runs per instance (default: 1)")
    parser.add_argument("--timeout", type=int, default=1800,
                        help="Timeout per instance in seconds (default: 1800)")
    parser.add_argument("--output", default="eval/results",
                        help="Output directory for results")
    args = parser.parse_args()

    # Load dataset
    instances = []
    with open(args.dataset, "r") as f:
        for line in f:
            instances.append(json.loads(line.strip()))

    # Filter instances if specified
    if args.instances:
        instances = [i for i in instances if i["instance_id"] in args.instances]

    if not instances:
        print("No instances to evaluate.")
        sys.exit(1)

    print(f"Loaded {len(instances)} instance(s)")

    # Determine agents to run
    if args.agent == "all":
        agents = [
            ("claude-code", "opus-4.6"),
            ("claude-code", "sonnet-4.6"),
            ("qwen", "qwen3.5-flash"),
            ("codex-cli", "gpt-4.1"),
        ]
    else:
        agents = [(args.agent, args.model)]

    # Run evaluations
    all_results = []
    for agent, model in agents:
        for instance in instances:
            for run_id in range(1, args.runs + 1):
                result = evaluate_instance(
                    instance, agent, model,
                    timeout=args.timeout, run_id=run_id
                )
                all_results.append(result)

    # Generate report
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    output_dir = os.path.join(args.output, timestamp)
    generate_report(all_results, output_dir)


if __name__ == "__main__":
    main()
