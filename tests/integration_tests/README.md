## Preparations

### Run integration tests locally

1. Run `make prepare_test_binaries` to download TiCDC related binaries for integration test.
If you need to specify a version, os or arch, you can use, for example: `make prepare_test_binaries community=true ver=v7.0.0 os=linux arch=amd64`. (The `community=true` option is necessary when you need to specify a version, arch or whatever.)

   > Note: Please be aware that if you manually specify the versions of the TiDB-related components during the download process, there is a potential risk of incompatibility if the specified versions do not match the current TiCDC version. For example, see issue #11507 where version mismatches caused incompatibility problems during integration tests.

   You should find these executables in the `ticdc/bin` directory after downloading successfully.
   * `tidb-server` # version >= 6.0.0-rc.1
   * `tikv-server` # version >= 6.0.0-rc.1
   * `pd-server`   # version >= 6.0.0-rc.1
   * `pd-ctl`      # version >= 6.0.0-rc.1
   * `tiflash`     # tiflash binary
   * `libc++.so, libc++abi.so, libgmssl.so, libtiflash_proxy.so` # some necessary so files related to tiflash
   * `sync_diff_inspector`
   * [go-ycsb](https://github.com/pingcap/go-ycsb)
   * [etcdctl](https://github.com/etcd-io/etcd/tree/master/etcdctl)
   * [jq](https://stedolan.github.io/jq/)
   * [minio](https://github.com/minio/minio)

   > You could download these binaries by yourself from [tidb-community-toolkit](https://download.pingcap.org/tidb-community-toolkit-v6.0.0-linux-amd64.tar.gz) and [tidb-community-server](https://download.pingcap.org/tidb-community-server-v6.0.0-linux-amd64.tar.gz).

2. These are programs/packages need be installed manually.
   * [mysql](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/) (the MySQL cli client,
     currently [mysql client 8.0 is not supported](https://github.com/pingcap/tidb/issues/14021))
   * [s3cmd](https://s3tools.org/download)
   * unzip
   * psmisc

   > You can install `unzip` and `psmisc` using `apt-get` (Ubuntu / Debian) or `yum` (RHEL).

   > Since the integration test cases will use port 3306 on localhost, please make sure in advance that port 3306 is
   > not occupied. (You’d like to stop the local MySQL service on port 3306, if there is one)

3. The user used to execute the tests must have permission to create the folder /tmp/tidb_cdc_test. All test artifacts
   will be written into this folder.

### Run integration tests in docker

The following programs must be installed:

* [docker](https://docs.docker.com/get-docker/)
* [docker-compose](https://docs.docker.com/compose/install/)

We recommend that you provide docker with at least 6+ cores and 8G+ memory. Of course, the more resources, the better.

## Running

### Unit Test

1. Unit test does not need any dependencies, just running `make unit_test` in root dir of source code, or cd into
   directory of a test case and run single case via `GO111MODULE=on go test -check.f TestXXX`.

### Integration Test

#### Run integration tests in docker

> **Warning:**
> These scripts and files may not work under the arm architecture,
> and we have not tested against it. We will try to resolve it as soon as possible.
>
> The script is designed to download necessary binaries from the PingCAP
> intranet by default, requiring access to the PingCAP intranet. However,
> if you want to download the community version, you can specify it through
> the `COMMUNITY` environment variable. For instance, you can use the following
> command as an example:
> `BRANCH=master COMMUNITY=true VERSION=v7.0.0 START_AT="clustered_index" make kafka_docker_integration_test_with_build`

1. If you want to run kafka tests,
   run `START_AT="clustered_index" make kafka_docker_integration_test_with_build`

2. If you want to run MySQL tests,
   run `CASE="clustered_index" make mysql_docker_integration_test_with_build`

3. Use the command `make clean_integration_test_images`
   to clean up the corresponding environment.

Some useful tips:

1. The log files for the test are mounted in the `./deployments/ticdc/docker-compose/logs` directory.

2. You can specify multiple tests to run in CASE, for example: `CASE="clustered_index kafka_messages"`. You can even
   use `CASE="*"` to indicate that you are running all tests。

3. You can specify in the [integration-test.Dockerfile](../../deployments/ticdc/docker/integration-test.Dockerfile)
   the version of other dependencies that you want to download, such as tidb, tikv, pd, etc.
   > For example, you can change `RUN ./download-integration-test-binaries.sh master`
   to `RUN ./download-integration-test-binaries.sh release-5.2`
   > to use the release-5.2 dependency.
   > Then rebuild the image with `make build_mysql_integration_test_images`.

#### Run integration tests locally

1. Run `make integration_test_build` to generate TiCDC related binaries for integration test

2. Run `make integration_test` to execute the integration tests. This command will

   1. Check that all required executables exist.
   2. Execute `tests/integration_tests/run.sh`

   > If want to run one integration test case only, just pass the CASE parameter, such as `make integration_test CASE=simple`.

   > If want to run integration test cases from the specified one, just pass the START_AT parameter, such as `make integration_test START_AT=simple` .

   > There exists some environment variables that you can set by yourself, variable details can be found in [test_prepare](_utils/test_prepare).

   > `MySQL sink` will be used by default, if you want to test `Kafka sink`, please run with `make integration_test_kafka CASE=simple`.

3. After executing the tests, run `make coverage` to get a coverage report at `/tmp/tidb_cdc_test/all_cov.html`.

## Writing new tests

1. New integration tests can be written as shell scripts in `tests/integration_tests/TEST_NAME/run.sh`. The script should
exit with a nonzero error code on failure.

2. Add TEST_NAME to existing group in [run_group.sh](./run_group.sh), or add a new group for it.

3. If you add a new group, the name of the new group must be added to CI.
   * [cdc-kafka-integration-light](https://github.com/PingCAP-QE/ci/blob/main/pipelines/pingcap/ticdc/latest/pod-pull_cdc_kafka_integration_light.yaml)
   * [cdc-kafka-integration-heavy](https://github.com/PingCAP-QE/ci/blob/main/pipelines/pingcap/ticdc/latest/pod-pull_cdc_kafka_integration_heavy.yaml)
   * [cdc-mysql-integration-light](https://github.com/PingCAP-QE/ci/blob/main/pipelines/pingcap/ticdc/latest/pod-pull_cdc_mysql_integration_light.yaml)
   * [cdc-mysql-integration-heavy](https://github.com/PingCAP-QE/ci/blob/main/pipelines/pingcap/ticdc/latest/pod-pull_cdc_mysql_integration_heavy.yaml)
