#!/bin/bash

# Exit on error
set -e

# Function to print usage
usage() {
  echo "Usage: $0 [--dry-run]"
  exit 1
}

# Parse command-line arguments
DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
  DRY_RUN=true
elif [[ "$#" -gt 0 ]]; then
  usage
fi

# Fetch latest annotated tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Extract version components
IFS='.' read -r MAJOR MINOR PATCH <<< "${LATEST_TAG#v}"

# Determine increment type based on commit messages
if [ "$LATEST_TAG" == "v0.0.0" ]; then
  # No existing tags, get all commits
  COMMITS=$(git log --oneline --pretty=format:"- %s (%h)" --no-merges HEAD)
else
  # Get commits since the latest tag
  COMMITS=$(git log --oneline --pretty=format:"- %s (%h)" --no-merges "${LATEST_TAG}"..HEAD)
fi

if echo "$COMMITS" | grep -qE "^- .+(\\(.+\\))?!:"; then
  MAJOR=$((MAJOR+1))
  MINOR=0
  PATCH=0
  INCREMENT="Major"
elif echo "$COMMITS" | grep -qE 'feat'; then
  MINOR=$((MINOR+1))
  PATCH=0
  INCREMENT="Minor"
elif echo "$COMMITS" | grep -qE 'fix'; then
  PATCH=$((PATCH+1))
  INCREMENT="Patch"
else
  echo "No increment-worthy commits found. Exiting."
  exit 0
fi

# Create new tag
NEW_TAG="v$MAJOR.$MINOR.$PATCH"
echo "Tag: $NEW_TAG"

ANNOTATED_MESSAGE="$INCREMENT $NEW_TAG

$(echo "$COMMITS" | grep -E "^- (feat|fix|.+(\\(.+\\))?!:)")"

echo "Message: \"$ANNOTATED_MESSAGE\""

if $DRY_RUN; then
  echo "[Dry Run] Tag creation and push skipped."
  exit 0
fi

git tag -a "$NEW_TAG" -m "$ANNOTATED_MESSAGE"

# Push the new tag to the remote
echo "Pushing tag $NEW_TAG"
git push origin "$NEW_TAG"

# # For github action
# echo "new_tag=$(echo $NEW_TAG)" >> $GITHUB_OUTPUT
