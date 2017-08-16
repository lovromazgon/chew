version="v0.0.1"
curdate=`date '+%F'`

go build -o chew-${version} -ldflags \
"-X bitbucket.org/lovromazgon/chew.Version=${version}
 -X bitbucket.org/lovromazgon/chew.ReleaseDate=${curdate}" chew/main.go
