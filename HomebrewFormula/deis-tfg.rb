class DeisTfg < Formula
  desc "The CLI for Deis Workflow"
  homepage "https://www.github.com/topfreegames/workflow-cli"
  url "https://github.com/topfreegames/workflow-cli/archive/v2.20.1.tar.gz"
  sha256 "2763a3bcf00f53d4f15e5f65115098b5036d94f3e36a87a66d8e59b69c748b69"
  version "2.20.1"

  depends_on "glide" => :build
  depends_on "go" => :build

  bottle do
    root_url "https://github.com/topfreegames/workflow-cli/releases/download/v2.20.1"
    cellar :any_skip_relocation
    sha256 "ea4b6c0b37b22f2171ba8230af88acd287df8a32fea16e402006717dea27dfcb" => :mojave
  end

  def install
    ENV["GOPATH"] = buildpath
    ENV["GLIDE_HOME"] = HOMEBREW_CACHE/"glide_home/#{name}"
    deispath = buildpath/"src/github.com/deis/workflow-cli"
    deispath.install Dir["{*,.git}"]

    cd deispath do
      system "glide", "install"
      system "go", "build", "-o", "build/deis",
        "-ldflags",
        "'-X=github.com/deis/workflow-cli/version.Version=v#{version}'"
      bin.install "build/deis"
    end
  end

  test do
    system bin/"deis", "logout"
  end
end
