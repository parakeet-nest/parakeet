# Synchronous usage
from gitingest import ingest
import os

GITHUB_REPOSITORY = os.environ.get('GITHUB_REPOSITORY', 'https://github.com/cyclotruc/gitingest')

env_value = os.getenv("INCLUDE_PATTERNS", "")
include_patterns = {pattern.strip() for pattern in env_value.split(",") if pattern.strip()}


summary, tree, content = ingest(
    source=GITHUB_REPOSITORY,
    include_patterns=include_patterns
)


#result = ingest(
#    source="https://github.com/user/repo",
#    max_file_size=5 * 1024 * 1024,  # 5MB
#    include_patterns={"*.py", "*.md"},
#    exclude_patterns={"test/*"},
#    branch="main",
#    output="output_dir"
#)


# Create data directory if it doesn't exist
os.makedirs('/app/data', exist_ok=True)

# Generate file from tree string
with open('/app/data/tree.txt', 'w') as f:
    f.write(tree)

# Generate file from summary string
with open('/app/data/summary.txt', 'w') as f:
    f.write(summary)

# Generate file from content string
with open('/app/data/content.txt', 'w') as f:
    f.write(content)

print(summary)
