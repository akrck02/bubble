#!/bin/sh
# Captures Bubble's material gallery, screen by screen, light and dark, and compares each with
# reference/: the same screens as Go UI for Android's gallery drew them (at commit de05770), the
# look Bubble's material style keeps. It prints how many pixels differ and leaves the pairs side
# by side. Needs a display (or xvfb-run) and ImageMagick.
#
#   cmd/gallery/compare.sh pictures/      run from bubble/goui
set -e
out=${1:-pictures}
mkdir -p "$out"
here=$(cd "$(dirname "$0")" && pwd)
go build -o "$out/gallery" "$here"
run=""
[ -z "$DISPLAY" ] && command -v xvfb-run >/dev/null && run="xvfb-run -a"
for screen in 0 1 2 3 4; do
	for mode in light dark; do
		flag=""
		[ "$mode" = dark ] && flag="-dark"
		name="$screen-$mode"
		HOME="$out/home" $run "$out/gallery" -width 412 -height 892 -screen "$screen" $flag -capture "$out/$name.png" </dev/null >/dev/null 2>&1
		differing=$(compare -metric AE -fuzz 2% "$here/reference/$name.png" "$out/$name.png" null: 2>&1 || true)
		total=$(identify -format "%[fx:w*h]" "$out/$name.png")
		echo "screen $name: $(awk "BEGIN { printf \"%.1f\", 100 * $differing / $total }")% of pixels differ"
		montage "$here/reference/$name.png" "$out/$name.png" -tile 2x -geometry +4+4 "$out/side-$name.png"
	done
done
