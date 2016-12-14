with import ../nixpkgs { };

let 
  myDeps =  (import ./myDeps.nix);
in

buildGoPackage rec {
  n = "pankat-server";
  name = "${n}-${version}";
  version = "0.0.1";

  goPackagePath = "github.com/nixcloud/${n}";

  shellHook = ''
    export HISTFILE=".zsh_history"
    kill $(pidof 'gocode')
    gocode
    gocode set propose-builtins true
    gocode set lib-path $GOPATH:`pwd`
    #kate server.go 2>/dev/null 1>/dev/null &
    #CompileDaemon -build 'go build -o pankat-server' -color &
  '';

  buildInputs = with myDeps; [ go net crypto captcha gomailv2 gocraft-web CompileDaemon.bin 
    inflection pq  ];
}

