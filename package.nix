{
  writeShellScriptBin,
  buildGoModule,
  antlr4,
  gitMinimal,
}:
buildGoModule rec {
  pname = "scadformat";
  version = "v0.9";

  src = ./.;
  vendorHash = "sha256-HOjfKFDG4otwu5TGXNtQCBQ7PURtPoeN8M8+uVHn5+4=";

  nativeBuildInputs = [
    antlr4
    (
      writeShellScriptBin "git" ''
        echo "${version}"
      ''
    )
  ];
  preBuild = ''
    go generate ./...
  '';
  subPackages = [
    "cmd/scadformat.go"
  ];
}
