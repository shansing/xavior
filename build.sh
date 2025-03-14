#!/usr/bin/bash

archs=(386 amd64 arm arm64 riscv64 ppc64 ppc64le loong64 mips mips64 mips64le mipsle s390x)
oss=(aix darwin dragonfly freebsd linux netbsd openbsd plan9 solaris windows)

mkdir -p bin

for arch in "${archs[@]}"; do
    for os in "${oss[@]}"; do
        suffix=""
        if [ "${os}" == "windows" ]; then
            suffix=".exe"
        fi

        output_file="bin/xavior-${os}-${arch}${suffix}"

        env GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -ldflags "-s -w" -o "${output_file}" xavior.go || echo "bad: ${os}-${arch}"
    done
done
