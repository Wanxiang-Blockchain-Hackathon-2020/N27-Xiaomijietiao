module github.com/despard/ioubackend

// go: no requirements found in vendor.conf

require (
	github.com/despard/log v0.0.0-20191214084335-596a71ad6b25
	github.com/julienschmidt/httprouter v1.3.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
	github.com/LyricTian/go.uuid v1.0.0 // indirect
	github.com/LyricTian/inject v0.0.0-20160612111808-5ff84550205f // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/bitly/go-simplejson v0.5.0
	github.com/cweill/gotests v1.5.3 // indirect
	github.com/despard/iouproto v0.0.0
	github.com/despard/stat v0.0.0-20191214084335-596a71ad6b25
	github.com/gavv/httpexpect v2.0.0+incompatible
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/satori/go.uuid v1.1.0
	golang.org/x/tools v0.0.0-20191227053925-7b8e75db28f4 // indirect
	gopkg.in/LyricTian/lib.v2 v2.2.7 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/oauth2.v2 v2.1.0 // indirect
	gopkg.in/oauth2.v3 v3.12.0
	gopkg.in/redis.v4 v4.2.4 // indirect
)

replace (
	github.com/despard/iouproto v0.0.0 => ../iouproto
	github.com/despard/log v0.0.0-20191214084335-596a71ad6b25 => ../log
	github.com/despard/stat v0.0.0-20191214084335-596a71ad6b25 => ../stat
)
