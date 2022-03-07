. test_harness.sh

# search for an asset in a snapshot checksums file
test_search_for_asset() {
  fixture=./test-fixtures/docker-sbom-cli-plugin_0.0.0-SNAPSHOT-ac27dcf_checksums.txt

  # search_for_asset [checksums-file-path] [name] [os] [arch] [format]

  # positive case
  actual=$(search_for_asset "${fixture}" "docker-sbom-cli-plugin" "linux" "amd64" "tar.gz")
  assertEquals "docker-sbom-cli-plugin_0.0.0-SNAPSHOT-ac27dcf_linux_amd64.tar.gz" "${actual}" "unable to find snapshot asset"

  # negative case
  actual=$(search_for_asset "${fixture}" "docker-sbom-cli-plugin" "linux" "amd64" "zip")
  assertEquals "" "${actual}" "found a snapshot asset but did not expect to (format)"
}

run_test_case test_search_for_asset