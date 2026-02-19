# Maintainer: Levi Renato <levirenato at gmail dot com>
pkgname=youtui
pkgver=1.1.0
pkgrel=1
pkgdesc="YouTube TUI player with playlist, thumbnails and Catppuccin themes"
arch=('x86_64' 'aarch64')
url="https://github.com/levirenato/YouTui"
license=('MIT')
depends=('mpv' 'yt-dlp' 'socat')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/refs/tags/v$pkgver.tar.gz")
b2sums=('SKIP')

prepare() {
    cd "YouTui-$pkgver"
    go mod download
}

build() {
    cd "YouTui-$pkgver"
    export CGO_ENABLED=0
    go build \
        -trimpath \
        -ldflags "-X main.Version=$pkgver -s -w" \
        -o "$pkgname" .
}

package() {
    cd "YouTui-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
