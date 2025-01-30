{pkgs ? import <nixpkgs> {}}:

pkgs.mkShell {
  name = "markslide";
  buildInputs = [
    pkgs.go
    pkgs.git
  ];
}
