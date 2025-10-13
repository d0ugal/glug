class Glug < Formula
  desc "JSON log parser and colorizer"
  homepage "https://github.com/d0ugal/glug"
  url "https://github.com/d0ugal/glug/archive/v1.0.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", "glug", "."
    bin.install "glug"
  end

  test do
    # Test that the binary works
    assert_match "glug", shell_output("#{bin}/glug --version")
  end
end
