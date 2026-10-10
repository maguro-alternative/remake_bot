#!/usr/bin/env python3
"""voicevox_core のリリース zip から voicevox_core.h だけを取り出す。

linter(declscope / molint)が cgo パッケージを型チェックするにはヘッダが必要だが、
リリース zip は約 700MB あるため、HTTP Range で必要な部分だけを読む。

使い方: python3 scripts/fetch-voicevox-header.py <出力ディレクトリ>
"""

import io
import os
import sys
import urllib.request
import zipfile

VERSION = "0.14.1"
URL = (
    "https://github.com/VOICEVOX/voicevox_core/releases/download/"
    f"{VERSION}/voicevox_core-linux-x64-cpu-{VERSION}.zip"
)


class HTTPRangeReader(io.RawIOBase):
    def __init__(self, url):
        head = urllib.request.urlopen(urllib.request.Request(url, method="HEAD"))
        self.url = head.url
        self.size = int(head.headers["Content-Length"])
        self.pos = 0

    def readable(self):
        return True

    def seekable(self):
        return True

    def tell(self):
        return self.pos

    def seek(self, offset, whence=io.SEEK_SET):
        base = {io.SEEK_SET: 0, io.SEEK_CUR: self.pos, io.SEEK_END: self.size}[whence]
        self.pos = base + offset
        return self.pos

    def readinto(self, buf):
        if self.pos >= self.size:
            return 0
        end = min(self.pos + len(buf), self.size) - 1
        req = urllib.request.Request(self.url, headers={"Range": f"bytes={self.pos}-{end}"})
        data = urllib.request.urlopen(req).read()
        buf[: len(data)] = data
        self.pos += len(data)
        return len(data)


def main():
    out_dir = sys.argv[1]
    os.makedirs(out_dir, exist_ok=True)
    archive = zipfile.ZipFile(io.BufferedReader(HTTPRangeReader(URL), 1 << 20))
    name = next(n for n in archive.namelist() if n.endswith("/voicevox_core.h"))
    with open(os.path.join(out_dir, "voicevox_core.h"), "wb") as f:
        f.write(archive.read(name))


if __name__ == "__main__":
    main()
