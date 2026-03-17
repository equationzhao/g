# get the latest git tag like v0.14.0 and remove the prefix v
latest := `git describe --abbrev=0 --tags | sed 's/v//'`
ldflags := "-ldflags='-s -w'"

COLOR_GREEN := "[0;32m"
COLOR_RED := "[0;31m"

# build lite version (minimal dependencies, smaller binary)
build-lite:
    mkdir -p build
    CGO_ENABLED=0 go build {{ldflags}} -o build/g-lite

# build full version (all features)  
build-full:
    mkdir -p build
    CGO_ENABLED=0 go build {{ldflags}} -tags="fuzzy mounts" -o build/g-full

# build with only fuzzy search support
build-fuzzy:
    mkdir -p build
    CGO_ENABLED=0 go build {{ldflags}} -tags="fuzzy" -o build/g-fuzzy

# build with only mounts support
build-mounts:
    mkdir -p build
    CGO_ENABLED=0 go build {{ldflags}} -tags="mounts" -o build/g-mounts

# build binaries for all platforms
build: 
    # build the binary in build/
    # Linux macOS Windows
    # 386 amd64 arm arm64
    mkdir -p build
    CGO_ENABLED=0 GOOS=linux GOARCH=386       go build {{ldflags}} -tags="fuzzy mounts" -o build/g-linux-386
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64     go build {{ldflags}} -tags="fuzzy mounts" -o build/g-linux-amd64
    CGO_ENABLED=0 GOOS=linux GOARCH=arm       go build {{ldflags}} -tags="fuzzy mounts" -o build/g-linux-arm
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64     go build {{ldflags}} -tags="fuzzy mounts" -o build/g-linux-arm64
    CGO_ENABLED=0 GOOS=linux GOARCH=loong64   go build {{ldflags}} -tags="fuzzy mounts" -o build/g-linux-loong64

    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64  go build {{ldflags}} -tags="fuzzy mounts" -o build/g-darwin-amd64
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64  go build {{ldflags}} -tags="fuzzy mounts" -o build/g-darwin-arm64

    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build {{ldflags}} -tags="fuzzy mounts" -o build/g-windows-amd64.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=386   go build {{ldflags}} -tags="fuzzy mounts" -o build/g-windows-386.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build {{ldflags}} -tags="fuzzy mounts" -o build/g-windows-arm64.exe

    upx --best build/g-linux-386
    upx --best build/g-linux-amd64
    upx --best build/g-linux-arm
    upx --best build/g-linux-arm64
    upx --best build/g-windows-amd64.exe
    upx --best build/g-windows-386.exe
#    upx doesn't support darwin-amd64, darwin-arm64, linux-loong64, windows-arm64
#    upx build/g-darwin-amd64
#    upx build/g-darwin-arm64
#    upx build/g-linux-loong64
#    upx build/g-windows-arm64.exe

# compress the binaries for all platforms
compress: 
    # g_OS_ARCH.tar.gz or g_Windows_ARCH.zip
    tar -zcvf build/g-Linux-386.tar.gz build/g-linux-386
    tar -zcvf build/g-Linux-amd64.tar.gz build/g-linux-amd64
    tar -zcvf build/g-Linux-arm.tar.gz build/g-linux-arm
    tar -zcvf build/g-Linux-arm64.tar.gz build/g-linux-arm64
    tar -zcvf build/g-Linux-loong64.tar.gz build/g-linux-loong64

    tar -zcvf build/g-Darwin-amd64.tar.gz build/g-darwin-amd64
    tar -zcvf build/g-Darwin-arm64.tar.gz build/g-darwin-arm64

    zip -r build/g-Windows-amd64.zip build/g-windows-amd64.exe
    zip -r build/g-Windows-386.zip build/g-windows-386.exe
    zip -r build/g-Windows-arm64.zip build/g-windows-arm64.exe

# build deb package for all ARCH
deb: 
    # g_VERSION_ARCH.deb
    # |__DEBIAN
    # |  |__control
    # |__usr
    #    |__local
    #       |__bin
    #       |  |__g
    #       |__share
    #           |__man
    #              |__man1
    #                 |__g.1.gz

    # amd64
    mkdir -vp build/g_{{latest}}_amd64/DEBIAN
    mkdir -vp build/g_{{latest}}_amd64/usr/local/bin
    mkdir -vp build/g_{{latest}}_amd64/usr/local/share/man/man1
    cp build/g-linux-amd64 build/g_{{latest}}_amd64/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_amd64/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_amd64/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Architecture: amd64" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_amd64/DEBIAN/control

    dpkg-deb -b --root-owner-group build/g_{{latest}}_amd64
    rm -rf build/g_{{latest}}_amd64


    # arm64
    mkdir -vp build/g_{{latest}}_arm64/DEBIAN
    mkdir -vp build/g_{{latest}}_arm64/usr/local/bin
    mkdir -vp build/g_{{latest}}_arm64/usr/local/share/man/man1
    cp build/g-linux-arm64 build/g_{{latest}}_arm64/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_arm64/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_arm64/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Architecture: arm64" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_arm64/DEBIAN/control

    dpkg-deb -b --root-owner-group build/g_{{latest}}_arm64
    rm -rf build/g_{{latest}}_arm64


    # 386
    mkdir -vp build/g_{{latest}}_386/DEBIAN
    mkdir -vp build/g_{{latest}}_386/usr/local/bin
    mkdir -vp build/g_{{latest}}_386/usr/local/share/man/man1
    cp build/g-linux-386 build/g_{{latest}}_386/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_386/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_386/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Architecture: 386" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_386/DEBIAN/control

    dpkg-deb -b --root-owner-group build/g_{{latest}}_386
    rm -rf build/g_{{latest}}_386

    # arm
    mkdir -vp build/g_{{latest}}_arm/DEBIAN
    mkdir -vp build/g_{{latest}}_arm/usr/local/bin
    mkdir -vp build/g_{{latest}}_arm/usr/local/share/man/man1
    cp build/g-linux-arm build/g_{{latest}}_arm/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_arm/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_arm/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Architecture: arm" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_arm/DEBIAN/control

    dpkg-deb -b --root-owner-group build/g_{{latest}}_arm
    rm -rf build/g_{{latest}}_arm

    # loong64
    mkdir -vp build/g_{{latest}}_loong64/DEBIAN
    mkdir -vp build/g_{{latest}}_loong64/usr/local/bin
    mkdir -vp build/g_{{latest}}_loong64/usr/local/share/man/man1
    cp build/g-linux-loong64 build/g_{{latest}}_loong64/usr
    cp man/g.1.gz build/g_{{latest}}_loong64/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_loong64/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_loong64/DEBIAN/control
    echo "Architecture: loong64" >> build/g_{{latest}}_loong64/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_loong64/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_loong64/DEBIAN/control

    dpkg-deb -b --root-owner-group build/g_{{latest}}_loong64
    rm -rf build/g_{{latest}}_loong64

# generate the checksum file for all build files
checksum: 
    shasum -a 256 build/* > build/checksum.txt

# build executables and compress and deb and checksum
genrelease: build compress deb checksum

# release to github
release: 
    gh release create v{{latest}} build/*

url := "https://github.com/Equationzhao/g"

# update aur 
aur: 
    # download the latest source code from {{url}}/"archive/refs/tags/v{{latest}}.tar.gz" if v{{latest}}.tar.gz not exists
    #!/usr/bin/env bash
    if [ ! -f v{{latest}}.tar.gz ]; then \
        wget -c {{url}}/archive/refs/tags/v{{latest}}.tar.gz -O v{{latest}}.tar.gz; \
    fi \
    # update PKGBUILD
    sed -i bak "s/sha256sums=.*/sha256sums=('$(shasum -a 256 v{{latest}}.tar.gz | choose 0)')/g" ../g-ls/PKGBUILD
    sed -i bak "s/pkgver=.*/pkgver={{latest}}/g" ../g-ls/PKGBUILD

    # update .SRCINFO
    sed -i bak "s/pkgver = .*/pkgver = {{latest}}/g" ../g-ls/.SRCINFO
    sed -i.bak "s/sha256sums = .*/sha256sums = '$(shasum -a 256 v{{latest}}.tar.gz | awk '{print $1}')'/g" ../g-ls/.SRCINFO
    sed -i.bak "s|source = g-.*::.*|source = g-{{latest}}.tar.gz::{{url}}/archive/refs/tags/v{{latest}}.tar.gz|g" ../g-ls/.SRCINFO

    cd ../g-ls
    # input git commit message
    cd ../g-ls && git add -u
    cd ../g-ls && git commit
    cd ../g-ls && git push

# update homebrew-tap
brew-tap:
    if [ ! -f v{{latest}}.tar.gz ]; then \
        wget -c {{url}}/archive/refs/tags/v{{latest}}.tar.gz -O v{{latest}}.tar.gz; \
    fi
    sed -i bak "s#url .*#url \"{{url}}/archive/refs/tags/v{{latest}}.tar.gz\", tag: \"v{{latest}}\"#g" ../homebrew-g/g-ls.rb
    sed -i bak "s/sha256 .*/sha256 \"$(shasum -a 256 v{{latest}}.tar.gz | choose 0)\"/g" ../homebrew-g/g-ls.rb
    sed -i bak '/assert_match/s/"[0-9.]*"/"{{latest}}"/' ../homebrew-g/g-ls.rb
    cd ../homebrew-g
    cd ../homebrew-g && git add -u
    cd ../homebrew-g && git commit
    cd ../homebrew-g && git push

# update homebrew-core
brew: 
    brew bump-formula-pr --strict --online g-ls

# update scoop
scoop:
    cd scoop && sh scoop.sh

# clean the build directory
clean: 
    rm -rf build

# golangci-lint
lint:
    golangci-lint run ./...

# format the code
format:
    gofumpt --extra -w -l .

# precheck the code(format and lint)
precheck: format lint

# generate the documentation
doc: 
    go build -tags 'doc'
    ./g 
    rm g

testcustomtheme:
    @sh ./script/theme_test.sh

# generate the theme
theme:
    go build -tags 'theme'
    ./g
    rm g

# generate the docs(doc and theme)
gendocs: doc theme

check-git-status:
    @if [ "$(git rev-parse HEAD)" == "$(git rev-parse v{{latest}})" ]; then \
      if [ -z "$(git status --porcelain)" ]; then \
        if [ "$(grep 'Version' internal/cli/version.go | awk '{print $4}' | sed 's/"//g')" == {{latest}} ]; then \
          echo "{{COLOR_GREEN}}latest tag v{{latest}} is on the current HEAD and the git status is clean. And version matches."; \
        else \
          echo "{{COLOR_RED}}latest tag v{{latest}} is on the current HEAD and the git status is clean."; \
          echo "{{COLOR_RED}}But version doesn't match. Please update the version in internal/cli/version.go."; \
          exit 1; \
        fi; \
      else \
        echo "{{COLOR_RED}}latest tag v{{latest}} is on the current HEAD but git status is dirty." && exit 1; \
      fi; \
    else \
      echo "{{COLOR_RED}}latest tag v{{latest}} isn't on the current HEAD." && exit 1; \
    fi;


check-install-script:
    @echo "git tag: v{{latest}}"
    @if [ "$(sh ./script/install.sh -v)" == {{latest}} ]; then \
      echo "{{COLOR_GREEN}}install.sh -v matches {{latest}} "; \
    else \
      echo "{{COLOR_RED}}install.sh -v "$(sh ./script/install.sh -v)""; \
      echo "{{COLOR_RED}}script version doesn't match {{latest}}"; \
    fi;

# check git tag and git status
check: check-install-script check-git-status

newtest:
    @sh ./script/new_test.sh

reproducetest:
    @sh ./script/reproduce_test_result.sh

test: testcustomtheme
    go test -cover -gcflags=all=-l -v ./...
    @echo "-------- start --------"
    go build
    @sh ./script/run_test.sh
    @rm g

newpatch:
    git add -u && git commit -m ":bookmark: new patch version"

newminor:
    git add -u && git commit -m ":bookmark: new minor version"

newmajor:
    git add -u && git commit -m ":bookmark: new major version"

bootstrap:
    sh script/install_dev_requirement.sh

all: precheck gendocs test check clean genrelease release
