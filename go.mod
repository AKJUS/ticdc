module github.com/pingcap/ticdc

go 1.23.8

toolchain go1.23.10

require (
	cloud.google.com/go/storage v1.39.1
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0
	github.com/BurntSushi/toml v1.5.0
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/IBM/sarama v1.41.2
	github.com/KimMachineGun/automemlimit v0.2.4
	github.com/agiledragon/gomonkey/v2 v2.11.0
	github.com/apache/pulsar-client-go v0.11.0
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2 v1.19.1
	github.com/aws/aws-sdk-go-v2/config v1.18.30
	github.com/aws/aws-sdk-go-v2/credentials v1.13.29
	github.com/aws/aws-sdk-go-v2/service/glue v1.58.1
	github.com/benbjohnson/clock v1.3.5
	github.com/bradleyjkemp/grpc-tools v0.2.5
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/cockroachdb/pebble v1.1.4-0.20250120151818-5dd133a1e6fb
	github.com/confluentinc/confluent-kafka-go/v2 v2.3.0
	github.com/coreos/go-semver v0.3.1
	github.com/fatih/color v1.18.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/go-sql-driver/mysql v1.7.1
	github.com/goccy/go-json v0.10.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.6.0
	github.com/google/btree v1.1.2
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.1
	github.com/imdario/mergo v0.3.16
	github.com/integralist/go-findroot v0.0.0-20160518114804-ac90681525dc
	github.com/jarcoal/httpmock v1.2.0
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.18.0
	github.com/linkedin/goavro/v2 v2.11.1
	github.com/mailru/easyjson v0.7.7
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pierrec/lz4/v4 v4.1.18
	github.com/pingcap/errors v0.11.5-0.20240318064555-6bd07397691f
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86
	github.com/pingcap/kvproto v0.0.0-20250224053625-b6a98c6bf02d
	github.com/pingcap/log v1.1.1-0.20241212030209-7e3ff8601a2a
	github.com/pingcap/tidb v1.1.0-beta.0.20250415053443-ea52376d91c8
	github.com/pingcap/tidb-dashboard v0.0.0-20240326110213-9768844ff5d7
	github.com/pingcap/tidb/pkg/parser v0.0.0-20250415053443-ea52376d91c8
	github.com/pingcap/tiflow v0.0.0-20250707050724-ac2fcf466cf3
	github.com/prometheus/client_golang v1.22.0
	github.com/r3labs/diff v1.1.0
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	github.com/robfig/cron v1.2.0
	github.com/shirou/gopsutil/v3 v3.24.5
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6
	github.com/stretchr/testify v1.10.0
	github.com/thanhpk/randstr v1.0.6
	github.com/tikv/client-go/v2 v2.0.8-0.20250327030528-1906b778d741
	github.com/tikv/pd v1.1.0-beta.0.20240407022249-7179657d129b
	github.com/tikv/pd/client v0.0.0-20250327162546-07c19b1f2f9f
	github.com/tinylib/msgp v1.1.6
	github.com/uber-go/atomic v1.4.0
	github.com/xdg/scram v1.0.5
	github.com/zeebo/assert v1.3.0
	go.etcd.io/etcd/api/v3 v3.5.12
	go.etcd.io/etcd/client/pkg/v3 v3.5.12
	go.etcd.io/etcd/client/v3 v3.5.12
	go.etcd.io/etcd/pkg/v3 v3.5.12
	go.etcd.io/etcd/server/v3 v3.5.12
	go.uber.org/atomic v1.11.0
	go.uber.org/goleak v1.3.0
	go.uber.org/mock v0.5.1
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.39.0
	golang.org/x/oauth2 v0.29.0
	golang.org/x/sync v0.13.0
	golang.org/x/sys v0.32.0
	golang.org/x/term v0.31.0
	golang.org/x/text v0.24.0
	golang.org/x/time v0.11.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.36.6
)

require (
	cloud.google.com/go v0.112.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.7 // indirect
	cloud.google.com/go/kms v1.15.8 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1581 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.30 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.20.1 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.14.3 // indirect
	github.com/blacktear23/go-proxyprotocol v1.0.6 // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/cakturk/go-netstat v0.0.0-20200220111822-e5b49efee7a5 // indirect
	github.com/carlmjohnson/flagext v0.21.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cheggaaa/pb/v3 v3.0.8 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cilium/ebpf v0.4.0 // indirect
	github.com/cloudfoundry/gosigar v1.3.6 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240606204812-0bbfbd93a7ce // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/coocood/bbloom v0.0.0-20190830030839-58deb6228d64 // indirect
	github.com/coocood/freecache v1.2.1 // indirect
	github.com/coocood/rtutil v0.0.0-20190304133409-c84515f646f2 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/danjacques/gofslock v0.0.0-20240212154529-d899e02bfe22 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20240924180020-3414d57e47da // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dolthub/maphash v0.1.0 // indirect
	github.com/dolthub/swiss v0.2.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/eapache/go-resiliency v1.4.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/edwingeng/deque v0.0.0-20191220032131-8596380dee17 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gavv/monotime v0.0.0-20190418164738-30dba4353424 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-ldap/ldap/v3 v3.4.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-mysql-org/go-mysql v1.12.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/go-resty/resty/v2 v2.11.0 // indirect
	github.com/goccy/go-reflect v1.2.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang/glog v1.2.4 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20241001023024-f4c0cfd0cf1d // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.3 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/improbable-eng/grpc-web v0.12.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/influxdata/tdigest v0.0.1 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jedib0t/go-pretty/v6 v6.2.2 // indirect
	github.com/jellydator/ttlcache/v3 v3.0.1 // indirect
	github.com/jfcg/sixb v1.3.8 // indirect
	github.com/jfcg/sorty/v2 v2.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/sqltocsv v0.0.0-20210428211105-a6d6801d59df // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/joomcode/errorx v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/ks3sdklib/aws-sdk-go v1.2.9 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.0.21 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncw/directio v1.0.5 // indirect
	github.com/ngaut/pools v0.0.0-20180318154953-b7bc8c42aac7 // indirect
	github.com/ngaut/sync2 v0.0.0-20141008032647-7a24ed77b2ef // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/otiai10/copy v1.14.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pingcap/badger v1.5.1-0.20241015064302-38533b6cbf8d // indirect
	github.com/pingcap/check v0.0.0-20211026125417-57bd13f7b5f0 // indirect
	github.com/pingcap/fn v1.0.0 // indirect
	github.com/pingcap/goleveldb v0.0.0-20191226122134-f82aafb29989 // indirect
	github.com/pingcap/sysutil v1.0.1-0.20240311050922-ae81ee01f3a5 // indirect
	github.com/pingcap/tipb v0.0.0-20250331100511-d2c561dad347 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.0 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/qri-io/jsonschema v0.2.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/segmentio/kafka-go v0.4.41-0.20230526171612-f057b1d369cd // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.3.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spkg/bom v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2 // indirect
	github.com/swaggo/gin-swagger v1.2.0 // indirect
	github.com/swaggo/swag v1.16.3 // indirect
	github.com/tiancaiamao/appdash v0.0.0-20181126055449-889f96f722a2 // indirect
	github.com/tiancaiamao/gp v0.0.0-20221230034425-4025bc8a4d4a // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tidwall/buntdb v1.3.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/vbauerster/mpb/v7 v7.5.3 // indirect
	github.com/wangjohn/quickselect v0.0.0-20161129230411-ed8402a42d5f // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	github.com/xiang90/probing v0.0.0-20221125231312-a49e3df8f510 // indirect
	github.com/xitongsys/parquet-go v1.6.3-0.20240520233950-75e935fc3e17 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	github.com/zyedidia/generic v1.2.1 // indirect
	go.etcd.io/bbolt v1.3.9 // indirect
	go.etcd.io/etcd/client/v2 v2.305.12 // indirect
	go.etcd.io/etcd/raft/v3 v3.5.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.22.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/api v0.170.0 // indirect
	google.golang.org/genproto v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.25.11 // indirect
	k8s.io/api v0.29.11 // indirect
	k8s.io/apimachinery v0.29.11 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0 // indirect
	sourcegraph.com/sourcegraph/appdash-data v0.0.0-20151005221446-73f23eafcf67 // indirect
)

// copy from TiDB
replace github.com/go-ldap/ldap/v3 v3.4.4 => github.com/yangkeao/ldap/v3 v3.4.5-0.20230421065457-369a3bab1117

replace github.com/IBM/sarama v1.41.2 => github.com/pingcap/sarama v1.41.2-pingcap-20250416
