module go-hep.org/x/hep

go 1.16

require (
	gioui.org v0.0.0-20210309172710-4b377aa89637
	github.com/apache/arrow/go/arrow v0.0.0-20201119084055-60ea0dcac5a8
	github.com/astrogo/fitsio v0.2.1
	github.com/campoy/embedmd v1.0.0
	github.com/go-mmap/mmap v0.6.0
	github.com/gonuts/binary v0.2.0
	github.com/gonuts/commander v0.3.1
	github.com/google/go-cmp v0.5.6
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/hashicorp/go-uuid v1.0.2
	github.com/jcmturner/gokrb5/v8 v8.4.2
	github.com/klauspost/compress v1.14.2
	github.com/peterh/liner v1.2.2
	github.com/pierrec/lz4/v4 v4.1.8
	github.com/pierrec/xxHash v0.1.5
	github.com/sbinet/npyio v0.6.0
	github.com/ulikunitz/xz v0.5.10
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
	golang.org/x/exp v0.0.0-20210220032938-85be41e4509f
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/tools v0.1.8
	gonum.org/v1/gonum v0.9.3
	gonum.org/v1/plot v0.10.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	modernc.org/ql v1.4.1
)

require (
	github.com/mattn/go-runewidth v0.0.13 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	modernc.org/db v1.0.4 // indirect
	modernc.org/lldb v1.0.4 // indirect
)

replace github.com/apache/arrow/go/arrow => git.sr.ht/~sbinet/go-arrow v0.1.1
