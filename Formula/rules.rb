# typed: false
# frozen_string_literal: true

# This is a reference formula for the `rules` CLI. In practice it is generated
# and kept up to date by GoReleaser (see .goreleaser.yaml) and published to the
# DeBros Homebrew tap, so `brew install debros/tap/rules` works. The url/sha256
# values below are placeholders filled in per release.
class Rules < Formula
  desc "Install DeBros engineering rules and skills into any project for Claude Code"
  homepage "https://github.com/DeBrosDAO/rules"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/DeBrosDAO/rules/releases/download/v0.1.0/rules_0.1.0_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_RELEASE_SHA256"
    end
    on_intel do
      url "https://github.com/DeBrosDAO/rules/releases/download/v0.1.0/rules_0.1.0_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_RELEASE_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/DeBrosDAO/rules/releases/download/v0.1.0/rules_0.1.0_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_RELEASE_SHA256"
    end
    on_intel do
      url "https://github.com/DeBrosDAO/rules/releases/download/v0.1.0/rules_0.1.0_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_RELEASE_SHA256"
    end
  end

  def install
    bin.install "rules"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/rules version")
  end
end
