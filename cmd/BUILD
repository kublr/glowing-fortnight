load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "root.go",
        "version.go",
    ],
    embedsrcs = [
        "templates/page.html",
    ],
    importpath = "github.com/kublr/snowflake-poc/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/snowflakedb/gosnowflake",
        "//vendor/github.com/spf13/cobra",
    ],
)
