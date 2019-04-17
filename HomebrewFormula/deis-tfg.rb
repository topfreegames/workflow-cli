class DeisTfg < Formula
  desc "The CLI for Deis Workflow"
  homepage "https://www.github.com/topfreegames/workflow-cli"
  url "https://github.com/topfreegames/workflow-cli/archive/v2.20.0.tar.gz"
  sha256 "ed5a6335c5ed292782deddebcc5e98a0fbc22f03f4619e3a4737ba7a85ae061d"
  version "2.20.0"

  depends_on "glide" => :build
  depends_on "go" => :build

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
