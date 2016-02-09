# Haru
comic crawler

## Install (for develop)

``
go get github.com/Perlmint/goautoenv
go get github.com/tools/godep

git clone git@github.com:if1live/haru.git
cd haru

goautoenv init haru
source .goenv/bin/activate
godep restore
goautoenv link

go test ./...
``
