#!/bin/bash

# Define the list of OS and architectures
declare -a platforms=(
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

# Loop through each platform and execute the build command
for platform in "${platforms[@]}"; do
    # Split the platform string into OS and architecture
    IFS='/' read -r target_os target_arch <<< "$platform"
    
    # Execute the build command
    echo "Building for $target_os/$target_arch..."
    GOOS=$target_os GOARCH=$target_arch go build -o build/$target_os/$target_arch/s3_proxy
done

echo "Build process completed."
