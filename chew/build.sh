version="v0.0.1"
curdate=`date '+%F'`

go build -ldflags \
"-X bitbucket.org/lovromazgon/chew.Version=${version}
 -X bitbucket.org/lovromazgon/chew.VersionDate=${curdate}"
