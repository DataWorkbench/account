module github.com/DataWorkbench/account

go 1.15

require (
	github.com/DataWorkbench/common v0.0.0-20220220155217-0f1831b67a5d
	github.com/DataWorkbench/glog v0.0.0-20210809050640-4960fd6de6ab
	github.com/DataWorkbench/gproto v0.0.0-20220216144250-797e691c176b
	github.com/DataWorkbench/loader v0.0.0-20210808061751-6e1ba7df5cdd
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-redis/redis/v8 v8.11.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/yunify/qingcloud-sdk-go v2.0.0-alpha.38+incompatible
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/gorm v1.20.6
)

replace (
	github.com/DataWorkbench/common => ../common
	github.com/DataWorkbench/gproto => ../gproto
)
