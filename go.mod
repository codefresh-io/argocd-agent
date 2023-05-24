module github.com/codefresh-io/argocd-listener

go 1.17

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/bradleyfalzon/ghinstallation/v2 v2.4.0
	github.com/codefresh-io/argocd-sdk v0.4.3
	github.com/codefresh-io/go-sdk v1.0.7-agent
	github.com/elliotchance/orderedmap v1.3.0
	github.com/fatih/structs v1.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-github v17.0.0+incompatible
	github.com/gookit/color v1.2.7
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/newrelic/go-agent v3.12.0+incompatible
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/afero v1.3.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.8.2
	github.com/thoas/go-funk v0.8.0
	github.com/ulule/deepcopier v0.0.0-20200430083143-45decc6639b6
	github.com/whilp/git-urls v1.0.0
	github.com/xanzy/go-gitlab v0.50.0
	golang.org/x/oauth2 v0.7.0
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0

)

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/Azure/go-autorest/autorest v0.9.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.5.0 // indirect
	github.com/Azure/go-autorest/autorest/date v0.1.0 // indirect
	github.com/Azure/go-autorest/logger v0.1.0 // indirect
	github.com/Azure/go-autorest/tracing v0.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/cloudflare/circl v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-github/v52 v52.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.8 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	golang.org/x/net => golang.org/x/net v0.7.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
