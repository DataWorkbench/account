module github.com/DataWorkbench/account

go 1.15

require (
	github.com/DataWorkbench/common v0.0.0-20220414084739-c0ad6d190e9d
	github.com/DataWorkbench/glog v0.0.0-20220302035436-25a1ae256704
	github.com/DataWorkbench/gproto v0.0.0-20220420085920-2f499b139536
	github.com/DataWorkbench/loader v0.0.0-20210808061751-6e1ba7df5cdd
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-redis/redis/v8 v8.11.0
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/yunify/qingcloud-sdk-go v2.0.0-alpha.38+incompatible
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/gorm v1.23.1
)

replace (
	github.com/DataWorkbench/common => ../common
	github.com/DataWorkbench/gproto => ../gproto
)
