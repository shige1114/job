#!/bin/bash
# Update package declarations and imports in all .kt files under backend/src/main/kotlin/backend/

files=$(find backend/src/main/kotlin/backend/ -name "*.kt" -type f)

for file in $files; do
    # Determine the expected package name based on directory structure
    # Path format: backend/src/main/kotlin/backend/dir1/dir2/.../file.kt
    # We want to remove 'backend/src/main/kotlin/' and the file name, then replace '/' with '.'
    rel_path=${file#backend/src/main/kotlin/}
    dir_path=$(dirname "$rel_path")
    expected_package=${dir_path//\//.}

    # Update package declaration
    # Only if it currently doesn't match
    current_package=$(grep "^package " "$file" | awk '{print $2}')
    if [ "$current_package" != "$expected_package" ]; then
        echo "Updating package in $file: $current_package -> $expected_package"
        sed -i '' "s/^package .*/package $expected_package/" "$file"
    fi

    # Update imports
    # Simple heuristic: find imports that might be internal. 
    # This is tricky without full AST parsing, but we can look for imports not starting with common libraries
    # and not starting with 'backend.'
    # For now, manually targeting known cases based on the file content.
done
