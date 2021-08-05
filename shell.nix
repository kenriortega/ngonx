let
 pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  name = "goproxy";
  buildInputs = with pkgs; [
    # Build
    gnumake
    go
    gosec
    golangci-lint
  ];

  shellHook = ''
    unset GOPATH GOROOT
    export GO111MODULE=on
  '';
}
