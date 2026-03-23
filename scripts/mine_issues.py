#!/usr/bin/env python3
"""
mine_issues.py — Scan free5GC sub-repositories for candidate task instances.

Uses the GitHub REST API to find closed issues with linked merged PRs,
then extracts metadata needed for SWE-Bench 5G task construction.

Requirements:
    pip install requests
    export GITHUB_TOKEN=ghp_xxx  (optional, raises rate limit from 60 to 5000/hr)

Usage:
    python mine_issues.py                    # scan all repos
    python mine_issues.py --repo amf smf     # scan specific repos
    python mine_issues.py --output candidates.json
"""

import argparse
import json
import os
import sys
import time
from datetime import datetime
from typing import Optional

import requests

# --- Configuration ---

GITHUB_API = "https://api.github.com"
ORG = "free5gc"

# All free5GC sub-repositories (NF name -> repo name)
REPOS = {
    "AMF":     "amf",
    "SMF":     "smf",
    "PCF":     "pcf",
    "UPF":     "upf",
    "UDM":     "udm",
    "UDR":     "udr",
    "AUSF":    "ausf",
    "NRF":     "nrf",
    "NSSF":    "nssf",
    "N3IWF":   "n3iwf",
    "NAS":     "nas",
    "NGAP":    "ngap",
    "PFCP":    "pfcp",
    "OpenAPI": "openapi",
    "Util":    "util",
    "Main":    "free5gc",
}

# Exclusion patterns for PR titles (skip non-bug-fix PRs)
EXCLUDE_PATTERNS = [
    "chore:", "chore(", "docs:", "doc:", "ci:", "ci(",
    "refactor:", "style:", "build:", "dependabot",
    "bump", "upgrade go", "upgrade version",
    "readme", "README", "CHANGELOG",
]


def get_session() -> requests.Session:
    """Create a requests session with optional GitHub token."""
    session = requests.Session()
    token = os.environ.get("GITHUB_TOKEN")
    if token:
        session.headers["Authorization"] = f"token {token}"
        print("[INFO] Using GITHUB_TOKEN for authentication")
    else:
        print("[WARN] No GITHUB_TOKEN set. Rate limit: 60 requests/hour.")
        print("       Set: export GITHUB_TOKEN=ghp_xxx")
    session.headers["Accept"] = "application/vnd.github.v3+json"
    return session


def rate_limit_wait(session: requests.Session):
    """Check rate limit and wait if needed."""
    resp = session.get(f"{GITHUB_API}/rate_limit")
    if resp.ok:
        data = resp.json()
        remaining = data["resources"]["core"]["remaining"]
        reset_time = data["resources"]["core"]["reset"]
        if remaining < 5:
            wait = reset_time - time.time() + 5
            if wait > 0:
                print(f"[WAIT] Rate limit nearly exhausted. Waiting {wait:.0f}s...")
                time.sleep(wait)


def fetch_all_pages(session: requests.Session, url: str, params: dict) -> list:
    """Fetch all pages of a paginated GitHub API endpoint."""
    results = []
    page = 1
    while True:
        params["page"] = page
        params["per_page"] = 100
        resp = session.get(url, params=params)
        if resp.status_code == 403:
            rate_limit_wait(session)
            continue
        if not resp.ok:
            print(f"[ERROR] {resp.status_code} for {url}: {resp.text[:200]}")
            break
        data = resp.json()
        if not data:
            break
        results.extend(data)
        page += 1
        if len(data) < 100:
            break
    return results


def get_merged_prs(session: requests.Session, repo: str) -> list:
    """Get all merged PRs for a repository."""
    url = f"{GITHUB_API}/repos/{ORG}/{repo}/pulls"
    params = {"state": "closed", "sort": "updated", "direction": "desc"}
    prs = fetch_all_pages(session, url, params)
    return [pr for pr in prs if pr.get("merged_at")]


def get_pr_details(session: requests.Session, repo: str, pr_number: int) -> dict:
    """Get detailed PR info including files changed."""
    url = f"{GITHUB_API}/repos/{ORG}/{repo}/pulls/{pr_number}"
    resp = session.get(url)
    if not resp.ok:
        return {}
    pr_data = resp.json()

    # Get files changed
    files_url = f"{url}/files"
    files_resp = session.get(files_url)
    files = files_resp.json() if files_resp.ok else []

    pr_data["changed_files_detail"] = files
    return pr_data


def is_bug_fix_pr(pr: dict) -> bool:
    """Heuristic to determine if a PR is a bug fix."""
    title = pr.get("title", "").lower()

    # Exclude non-bug-fix PRs
    for pattern in EXCLUDE_PATTERNS:
        if pattern.lower() in title:
            return False

    # Include if title suggests a fix
    fix_indicators = ["fix", "bug", "crash", "panic", "nil pointer",
                      "error", "issue", "resolve", "patch", "hotfix"]
    for indicator in fix_indicators:
        if indicator in title:
            return True

    # Include if PR body references an issue
    body = pr.get("body", "") or ""
    if "fixes #" in body.lower() or "closes #" in body.lower() or "resolves #" in body.lower():
        return True

    return False


def has_code_changes(files: list) -> bool:
    """Check if PR modifies actual source code (not just config/docs)."""
    code_extensions = {".go", ".c", ".h"}
    for f in files:
        filename = f.get("filename", "")
        if any(filename.endswith(ext) for ext in code_extensions):
            # Exclude test files from "code changes" check
            if "_test.go" not in filename:
                return True
    return False


def classify_difficulty(files: list, additions: int, deletions: int) -> str:
    """Classify task difficulty based on change scope."""
    code_files = [f for f in files
                  if f["filename"].endswith(".go") and "_test.go" not in f["filename"]]
    n_files = len(code_files)
    total_lines = additions + deletions

    if n_files <= 1 and total_lines < 50:
        return "easy"
    elif n_files <= 3 and total_lines < 200:
        return "medium"
    else:
        return "hard"


def classify_bug_type(title: str, body: str) -> str:
    """Classify bug type from PR title and body."""
    text = (title + " " + (body or "")).lower()
    if "nil pointer" in text or "nil dereference" in text or "panic" in text:
        return "nil-pointer-crash"
    elif "index out of range" in text or "out of bound" in text:
        return "index-out-of-range"
    elif "crash" in text or "segfault" in text:
        return "crash"
    elif "race" in text or "concurrent" in text or "goroutine" in text:
        return "concurrency"
    elif "auth" in text or "security" in text or "permission" in text:
        return "security"
    elif "missing" in text or "validation" in text or "check" in text:
        return "missing-validation"
    elif "logic" in text or "incorrect" in text or "wrong" in text:
        return "logic-error"
    else:
        return "other"


def extract_linked_issue(pr: dict, session: requests.Session, repo: str) -> Optional[dict]:
    """Try to find the linked issue for a PR."""
    body = pr.get("body", "") or ""

    # Parse "fixes #123", "closes #123", "resolves #123"
    import re
    patterns = [
        r"(?:fix|fixes|fixed|close|closes|closed|resolve|resolves|resolved)\s+#(\d+)",
        r"(?:fix|fixes|fixed|close|closes|closed|resolve|resolves|resolved)\s+(?:free5gc/\w+)?#(\d+)",
        r"https://github\.com/free5gc/\w+/issues/(\d+)",
    ]

    for pattern in patterns:
        match = re.search(pattern, body, re.IGNORECASE)
        if match:
            issue_num = int(match.group(1))
            # Try to fetch the issue
            issue_url = f"{GITHUB_API}/repos/{ORG}/{repo}/issues/{issue_num}"
            resp = session.get(issue_url)
            if resp.ok:
                return resp.json()
            # Try main repo
            issue_url = f"{GITHUB_API}/repos/{ORG}/free5gc/issues/{issue_num}"
            resp = session.get(issue_url)
            if resp.ok:
                return resp.json()

    return None


def process_repo(session: requests.Session, nf_name: str, repo: str) -> list:
    """Process a single repository and return candidate task instances."""
    print(f"\n{'='*60}")
    print(f"Scanning {ORG}/{repo} ({nf_name})...")
    print(f"{'='*60}")

    candidates = []
    prs = get_merged_prs(session, repo)
    print(f"  Found {len(prs)} merged PRs")

    bug_fix_prs = [pr for pr in prs if is_bug_fix_pr(pr)]
    print(f"  Filtered to {len(bug_fix_prs)} bug-fix PRs")

    for pr in bug_fix_prs:
        pr_number = pr["number"]
        print(f"  Processing PR #{pr_number}: {pr['title'][:60]}...")

        # Get detailed PR info
        details = get_pr_details(session, repo, pr_number)
        if not details:
            continue

        files = details.get("changed_files_detail", [])
        if not has_code_changes(files):
            print(f"    Skipped: no source code changes")
            continue

        # Extract commit info
        merge_commit = details.get("merge_commit_sha", "")
        base_sha = details.get("base", {}).get("sha", "")

        # Get parent commit (commit before merge)
        if merge_commit:
            commit_url = f"{GITHUB_API}/repos/{ORG}/{repo}/commits/{merge_commit}"
            commit_resp = session.get(commit_url)
            if commit_resp.ok:
                parents = commit_resp.json().get("parents", [])
                if parents:
                    parent_commit = parents[0]["sha"]
                else:
                    parent_commit = base_sha
            else:
                parent_commit = base_sha
        else:
            parent_commit = base_sha

        # Try to find linked issue
        issue = extract_linked_issue(pr, session, repo)
        issue_url = issue["html_url"] if issue else ""
        issue_body = issue.get("body", "") if issue else ""
        problem_statement = issue_body if issue_body else (pr.get("body", "") or "")

        # Classify
        additions = details.get("additions", 0)
        deletions = details.get("deletions", 0)
        difficulty = classify_difficulty(files, additions, deletions)
        bug_type = classify_bug_type(pr["title"], pr.get("body", ""))

        # Affected files
        affected_files = [f["filename"] for f in files
                         if f["filename"].endswith(".go") and "_test.go" not in f["filename"]]

        candidate = {
            "instance_id": f"{repo}_pr{pr_number}",
            "nf_type": nf_name,
            "repo": f"{ORG}/{repo}",
            "pr_number": pr_number,
            "pr_title": pr["title"],
            "pr_url": pr["html_url"],
            "pr_merged_at": pr.get("merged_at", ""),
            "issue_url": issue_url,
            "problem_statement": problem_statement[:2000],  # truncate
            "commit_id": merge_commit[:12] if merge_commit else "",
            "parent_commit": parent_commit[:12] if parent_commit else "",
            "full_commit_id": merge_commit,
            "full_parent_commit": parent_commit,
            "difficulty": difficulty,
            "bug_type": bug_type,
            "affected_files": affected_files,
            "files_changed": len(affected_files),
            "lines_added": additions,
            "lines_deleted": deletions,
            "language": "go",
        }

        candidates.append(candidate)
        print(f"    Added: {difficulty} | {bug_type} | {len(affected_files)} files | +{additions}/-{deletions}")

    return candidates


def print_summary(all_candidates: list):
    """Print summary statistics."""
    print(f"\n{'='*60}")
    print(f"SUMMARY")
    print(f"{'='*60}")
    print(f"Total candidates: {len(all_candidates)}")

    # By NF
    print(f"\nBy NF type:")
    nf_counts = {}
    for c in all_candidates:
        nf = c["nf_type"]
        nf_counts[nf] = nf_counts.get(nf, 0) + 1
    for nf, count in sorted(nf_counts.items(), key=lambda x: -x[1]):
        print(f"  {nf:8s}: {count}")

    # By difficulty
    print(f"\nBy difficulty:")
    diff_counts = {}
    for c in all_candidates:
        d = c["difficulty"]
        diff_counts[d] = diff_counts.get(d, 0) + 1
    for d in ["easy", "medium", "hard"]:
        print(f"  {d:8s}: {diff_counts.get(d, 0)}")

    # By bug type
    print(f"\nBy bug type:")
    bt_counts = {}
    for c in all_candidates:
        bt = c["bug_type"]
        bt_counts[bt] = bt_counts.get(bt, 0) + 1
    for bt, count in sorted(bt_counts.items(), key=lambda x: -x[1]):
        print(f"  {bt:25s}: {count}")


def main():
    parser = argparse.ArgumentParser(description="Mine free5GC issues for SWE-Bench 5G")
    parser.add_argument("--repo", nargs="*", default=None,
                        help="Specific repo names to scan (e.g., amf smf pcf)")
    parser.add_argument("--output", default="candidates.json",
                        help="Output file path (default: candidates.json)")
    args = parser.parse_args()

    session = get_session()

    # Determine which repos to scan
    if args.repo:
        repos_to_scan = {nf: r for nf, r in REPOS.items() if r in args.repo}
    else:
        repos_to_scan = REPOS

    print(f"Scanning {len(repos_to_scan)} repositories...")

    all_candidates = []
    for nf_name, repo in repos_to_scan.items():
        try:
            candidates = process_repo(session, nf_name, repo)
            all_candidates.extend(candidates)
        except Exception as e:
            print(f"[ERROR] Failed to process {repo}: {e}")
            continue

    # Sort by merge date (newest first)
    all_candidates.sort(key=lambda x: x.get("pr_merged_at", ""), reverse=True)

    # Print summary
    print_summary(all_candidates)

    # Save results
    output_path = args.output
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(all_candidates, f, indent=2, ensure_ascii=False)
    print(f"\nResults saved to: {output_path}")
    print(f"Next step: review candidates and select 30+ for task construction")


if __name__ == "__main__":
    main()
