# Synchronous usage
from gitingest import ingest
import os

GITHUB_REPOSITORY = os.environ.get('GITHUB_REPOSITORY', 'https://github.com/cyclotruc/gitingest')

summary, tree, content = ingest(GITHUB_REPOSITORY)

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
