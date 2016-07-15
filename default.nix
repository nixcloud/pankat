with import <nixpkgs> { };
with goPackages;

let

   CompileDaemon = buildGoPackage rec {
    rev = "051a9ad079bf636e3db7fab6cfab1c4629b22519";
    name = "CompileDaemon-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/githubnemo/CompileDaemon";
    #doCheck = true;
    buildInputs = [ fatih-color fsnotify go-isatty ansicolor ];

    src = fetchFromGitHub {
      inherit rev;
      owner = "githubnemo";
      repo = "CompileDaemon";
      sha256 = "1w4srfbyddw977q3329ww6czwhsb7lbkvf4px9dknjavzvxskycg";
    };
  };

   fatih-color = buildGoPackage rec {
    rev = "f773d4c806cc8e4a5749d6a35e2a4bbcd71443d6";
    name = "fatih-color-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/fatih/color";
    #doCheck = true;
    buildInputs = [ go-isatty ansicolor ];

    src = fetchFromGitHub {
      inherit rev;
      owner = "fatih";
      repo = "color";
      sha256 = "1bd69gm6nig0g8zcsav68xs0h8sfdjj87fdrly9gf2k2r366bp9b";
    };
  };

   go-isatty = buildGoPackage rec {
    rev = "ae0b1f8f8004be68d791a576e3d8e7648ab41449";
    name = "go-isatty-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/mattn/go-isatty";

    src = fetchFromGitHub {
      inherit rev;
      owner = "mattn";
      repo = "go-isatty";
      sha256 = "0qrcsh7j9mxcaspw8lfxh9hhflz55vj4aq1xy00v78301czq6jlj";
    };
  };

   ansicolor = buildGoPackage rec {
    rev = "a5e2b567a4dd6cc74545b8a4f27c9d63b9e7735b";
    name = "ansicolor-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/shiena/ansicolor";

    src = fetchFromGitHub {
      inherit rev;
      owner = "shiena";
      repo = "ansicolor";
      sha256 = "0gwplb1b4fvav1vjf4b2dypy5rcp2w41vrbxkd1dsmac870cy75p";
    };
  };
  fsnotify = buildGoPackage rec {
    rev = "4894fe7efedeeef21891033e1cce3b23b9af7ad2";
    name = "fsnotify-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/howeyc/fsnotify";

    src = fetchFromGitHub {
      inherit rev;
      owner = "howeyc";
      repo = "fsnotify";
      sha256 = "09r3h200nbw8a4d3rn9wxxmgma2a8i6ssaplf3zbdc2ykizsq7mn";
    };
  };

###### kingpin #######

#  kingpin = buildGoPackage rec {
#    rev = "8852570bd3865e9c4d4cb7cf5001c4295b07cad5";
#    name = "kingpin-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "gopkg.in/alecthomas/kingpin.v2";
#
#    buildInputs = [ alecthomas-template alecthomas-units ];
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "alecthomas";
#      repo = "kingpin";
#      sha256 = "16ldz1axbkl6w44s4f57jf90bz9idhlj3c9qf25lvrwq8jljhkha";
#    };
#  };
#
#  alecthomas-template = buildGoPackage rec {
#    rev = "14fd436dd20c3cc65242a9f396b61bfc8a3926fc";
#    name = "alecthomas-template-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "github.com/alecthomas/template";
#
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "alecthomas";
#      repo = "template";
#      sha256 = "19rzvvcgvr1z2wz9xpqsmlm8syizbpxjp5zbzgakvrqlajpbjvx2";
#    };
#  };
#
#  alecthomas-units = buildGoPackage rec {
#    rev = "2efee857e7cfd4f3d0138cc3cbb1b4966962b93a";
#    name = "alecthomas-units-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "github.com/alecthomas/units";
#
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "alecthomas";
#      repo = "units";
#      sha256 = "1j65b91qb9sbrml9cpabfrcf07wmgzzghrl7809hjjhrmbzri5bl";
#    };
#  };





in

buildGoPackage rec {
  n = "pankat";
  name = "${n}-${version}";
  version = "0.0.1";

  goPackagePath = "github.com/nixcloud/${n}";

  shellHook = ''
    export HISTFILE=".zsh_history"
    kill $(pidof 'gocode')
    gocode
    gocode set propose-builtins true
    gocode set lib-path $GOPATH:`pwd`
    kate pankat.go 2>/dev/null 1>/dev/null &
    CompileDaemon -build 'go build -o pankat' -color &
  '';

  buildInputs = [ go net pandoc CompileDaemon.bin rsync git ];
}

