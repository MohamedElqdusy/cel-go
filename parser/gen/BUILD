load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(
    default_visibility = ["//parser:__subpackages__"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cel_base_listener.go",
        "cel_base_visitor.go",
        "cel_lexer.go",
        "cel_listener.go",
        "cel_parser.go",
        "cel_visitor.go",
    ],
    deps = [
        "@com_github_antlr//runtime/Go/antlr:go_default_library",
    ],
    data = [
        "CEL.tokens",
        "CELLexer.tokens",
    ],
)
