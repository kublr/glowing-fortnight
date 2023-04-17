load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_layer", "container_push")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("@rules_pkg//:mappings.bzl", "strip_prefix")

# Gazelle task to refresh bazel go dependencies
load("@bazel_gazelle//:def.bzl", "gazelle")

# This needs to be here for gazelle to work correctly
#
# gazelle:prefix github.com/kublr/snowflake-poc

# This is excluded because gazelle does not understand cgo pkg-config directive in the mongocrypt.go
# Excluding it is ok because it is only used when compiled with "ccalloc" tag, which we don't do.
# "ccalloc" tag is necessary to use C++ memory allocator in Apache Arrow library.
#
# gazelle:exclude vendor/github.com/apache/arrow/go/v10/arrow/memory/internal/cgoalloc/allocator.go

# Updating vendored dependencies - run in the workspace root:
#   go get -u <package>
#   go mod tidy
#   go mod vendor
#   bazel run :gazelle-update-repos
#   find vendor '(' -name BUILD.bazel -o -name BUILD ')' -delete
#   bazel run :gazelle
#   find . '(' -name BUILD.bazel -o -name BUILD ')' -size 0c -delete

# run "bazel run :gazelle" to update dependencies in BUILD files (or generate new) in go packages
gazelle(name = "gazelle")

# run "bazel run :gazelle-fix" to update dependencies in BUILD files (or generate new) in go packages
gazelle(
    name = "gazelle-fix",
    command = "fix",
)

# run "bazel run :gazelle-update-repos" or "bazel run :gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies -prune"
# to update external dependencies using go.mod in go_dependencies macros in deps.bzl file
gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

# # Docker image
pkg_tar(
    name = "image-tar-binary",
    files = {":snowflake-poc": "snowflake-poc"},
    mode = "0555",
    package_dir = "/opt/snowflake-poc",
    strip_prefix = strip_prefix.from_pkg(),
)

container_layer(
    name = "image-layer-binary",
    tars = [":image-tar-binary"],
)

container_image(
    name = "image",
    base = "@base-distroless-image//image",
    creation_time = "0",
    entrypoint = ["/opt/snowflake-poc/snowflake-poc"],
    layers = [":image-layer-binary"],
    ports = ["4080"],
    workdir = "/opt/snowflake-poc",
)

container_push(
    name = "release",
    format = "Docker",
    image = ":image",
    registry = "docker-registry.kcp.kublr-demo.com",
    repository = "snowflake-poc/ui",
    skip_unchanged_digest = False,
    tag = "{STABLE_BUILD_GIT_COMMIT}",
)

# gazelle generated go targets
go_library(
    name = "snowflake-poc_lib",
    srcs = ["main.go"],
    importpath = "github.com/kublr/snowflake-poc",
    visibility = ["//visibility:private"],
    deps = ["//cmd"],
)

go_binary(
    name = "snowflake-poc",
    embed = [":snowflake-poc_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "snowflake-poc_test",
    srcs = ["main_test.go"],
    embed = [":snowflake-poc_lib"],
)
