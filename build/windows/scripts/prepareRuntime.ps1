function prepareRuntime() {
    Copy-Item -Path ./build/assets/* -Destination ./output/runtime -Force -Recurse
    Copy-Item -Path ./build/windows/config.toml -Destination ./output/runtime -Force
}

prepareRuntime