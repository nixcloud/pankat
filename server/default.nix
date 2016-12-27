with import <nixpkgs> { };
#  nix-shell -I nixpkgs=/home/joachim/Desktop/projects/nixcloud/nixcloud-backend/nixpkgs

let 
  myDeps =  (import ./myDeps.nix);
in

buildGoPackage rec {
  n = "pankat-server";
  name = "${n}-${version}";
  version = "0.0.1";

  goPackagePath = "github.com/nixcloud/${n}";
  goDeps = ./deps.json;

  shellHook = ''
    export HISTFILE=".zsh_history"
    kill $(pidof 'gocode')
    gocode
    gocode set propose-builtins true
    gocode set lib-path $GOPATH:`pwd`
    #kate server.go 2>/dev/null 1>/dev/null &
    #CompileDaemon -build 'go build -o pankat-server' -color &
  '';

  buildInputs = with myDeps; [ go crypto captcha gomailv2 gocraft-web CompileDaemon.bin inflection pq  ];
}

