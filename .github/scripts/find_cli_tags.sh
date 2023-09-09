#!/bin/bash

# Check if an argument was provided
if [ $# -eq 0 ]; then
  echo "Please provide a Git tag as an argument."
  exit 1
fi

# Get the input Git tag
input_tag="$1"

# Find the commit associated with the input tag
commit_hash=$(git rev-list -n 1 "$input_tag" 2>/dev/null)

# Check if the tag exists and get tags on the same commit
if [ -z "$commit_hash" ]; then
  echo "Tag '$input_tag' not found or is not associated with any commit."
  exit 1
else
  # Get all tags that are on the commit
  matching_tags=$(git tag --contains "$commit_hash" 2>/dev/null)

  # Filter tags that follow the format "cli/*"
  cli_tags=()
  for tag in "${matching_tags[@]}"; do
    if [[ "$tag" == "cli/"* ]]; then
      cli_tags+=("$tag")
    fi
  done

  # Sort the filtered tags in descending numerical order
  sorted_tags=($(printf "%s\n" "${cli_tags[@]}" | sed -n 's/cli\/\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)-\(.*\)/\4 \1 \2 \3/p' | sort -t' ' -k2,2nr -k3,3nr -k4,4nr | awk '{print "cli/"$2"."$3"."$4"-"$1}'))

  # Display the top (first) sorted tag
  if [ ${#sorted_tags[@]} -eq 0 ]; then
    echo "NO_RELEASE_TAG_FOUND"
    exit 1
  else
    echo "${sorted_tags[0]}"
  fi
fi
