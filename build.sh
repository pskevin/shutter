#!/usr/bin/env bash

package_root="github.com/pskevin/shutter"
package_names=("cluster" "master")
platforms=("windows/amd64" "darwin/386" "linux/amd64")

for package_name in "${package_names[@]}"
do
    for platform in "${platforms[@]}"
    do
        platform_split=(${platform//\// })
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        output_name=$package_name'@'$GOOS'-'$GOARCH
        if [ $GOOS = "windows" ]; then
            output_name+='.exe'
        fi
        package=$package_root'/'$package_name
        printf "\nBuilding $package as $output_name\n"
        env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
        fi
    done
done