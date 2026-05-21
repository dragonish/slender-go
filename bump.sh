#!/bin/bash

echo "Version bump script"

get_current_version() {
  ver=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
  ver=${ver:-"0.0.0"}
  echo "$ver"
}

current=$(get_current_version)

IFS='.' read -r -a version_parts <<< "$current"
major="${version_parts[0]}"
minor="${version_parts[1]}"
patch="${version_parts[2]}"

has_breaking_change=false
has_feature_change=false
has_patch_change=false

commits=$(git log --pretty=format:"%s %b" $(git describe --tags --abbrev=0 HEAD)..HEAD)

echo "Checking commit history:"
echo "$commits"
echo "----------------------------------"

while IFS= read -r commit; do
    if echo "$commit" | grep -q "BREAKING CHANGE:"; then
        has_breaking_change=true
    fi
    if echo "$commit" | grep -qE '^(feat|deprecate):' || echo "$commit" | grep -qE '^(feat|deprecate)\([^)]*\):'; then
        has_feature_change=true
    fi

    if [[ "$has_breaking_change" == false && "$has_feature_change" == false ]]; then
        if [[ -n "$commit" ]]; then
            has_patch_change=true
        fi
    fi
done <<< "$commits"

if [[ "$has_breaking_change" == true ]]; then
    echo "Found breaking change"
    new_version="$((major + 1)).0.0"
    update_type="MAJOR"
elif [[ "$has_feature_change" == true ]]; then
    echo "Found feature/deprecation commit"
    new_version="$major.$((minor + 1)).0"
    update_type="MINOR"
else
    new_version="$major.$minor.$((patch + 1))"
    update_type="PATCH"
fi

echo "----------------------------------"
echo "Current version: v$current"
echo "Detected $update_type version change"
echo "New version: v$new_version"

echo "Creating Git tag: $new_version ..."
git tag -a "$new_version" -m "Release v$new_version"

current_version=$(get_current_version)
echo "Current project version: $current_version"

if [ "$current_version" == "$new_version" ]; then
  echo "Done, exiting script"
  exit 0
else
  echo "Current version does not match expected version, please check manually!"
  exit 1
fi
