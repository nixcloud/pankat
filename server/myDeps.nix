with import <nixpkgs> { };
with goPackages;
with pkgs;

rec {
  

    

#   #golang.org/x/exp/inotify
#   inotify = buildGoPackage rec {
#     rev = "a16ab4a434a95ebed3f8c4b58b88bcc2b1fa64b2";
#     name = "inotify-${stdenv.lib.strings.substring 0 7 rev}";
#     goPackagePath = "golang.org/x/exp/inotify";
#     buildInputs = [  ];
# 
# #     golang.org/x/exp/io/i2c
# #     golang.org/x/exp/shiny/screen
# #     golang.org/x/mobile/event/mouse
# #     golang.org/x/image/math/fixed
# #     github.com/BurntSushi/xgb
# #     github.com/BurntSushi/xgb/xproto
# #     github.com/BurntSushi/xgb
# #     github.com/BurntSushi/xgb/xproto
# #     github.com/BurntSushi/xgb
# #     golang.org/x/exp/shiny/widget/node
# #     golang.org/x/exp/shiny/widget/theme
# #     golang.org/x/exp/shiny/driver/gldriver
# #     golang.org/x/exp/shiny/widget/node
# #     golang.org/x/exp/shiny/widget/theme
# #     golang.org/x/mobile/gl
# #     golang.org/x/exp/shiny/unit
# #     golang.org/x/image/font
# #     golang.org/x/image/font/inconsolata
# #     golang.org/x/image/math/fixed
# 
#     src = fetchFromGitHub {
#       inherit rev;
#       owner = "golang";
#       repo = "exp";
#       sha256 = "01p7nzysfcd7myxvy4ffr0dl0wq6r8zxclwank6d0xbgz52b6n6h";
#     };
#   };



  gomailv2 = buildGoPackage rec {
    rev = "5ceb8e65415e45e1262fb385212b8193b55c0f99";
    name = "gomailv2-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "gopkg.in/gomail.v2";
  
    src =  fetchFromGitHub {
      inherit rev;
      owner = "go-gomail";
      repo = "gomail";
      sha256 = "03qqcfi50lp54dvfk5fxxfnwsvas6n2029ikgp1i7nd59jv7i1bm";
    };
  };

  captcha = buildGoPackage rec {
    rev = "9e952142169c3cd6268c6482a3a61c121536aca2";
    name = "captcha-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/dchest/captcha";
  
    src =  fetchFromGitHub {
      inherit rev;
      owner = "dchest";
      repo = "captcha";
      sha256 = "061lad5ynxjxq4ym656dyla5bym7vg4cqvrrcq0yp1h3lfkjxvbz";
    };
  };

  crypto = buildGoPackage rec {
    rev = "7b85b097bf7527677d54d3220065e966a0e3b613";
    name = "crypto-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "golang.org/x/crypto";
    #doCheck = true;

    src = fetchFromGitHub {
      inherit rev;
      owner = "golang";
      repo = "crypto";
      sha256 = "0k21nnf0nszgbvml74sn68wc6p77pxbnfpi04dgarg6byd5rvxii";
    };
  };

#  net = buildGoPackage rec {
#    rev = "I41c1adfdebd7dff0f75c2ba2b669f13b163984a2";
#    rev = "04b9de9b512f58addf28c9853d50ebef61c3953e";
#    name = "net-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "golang.org/x/net";
    #doCheck = true;

#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "golang";
#      repo = "net";
#      sha256 = "";
#    };
#  };

  gocraft-web = buildGoPackage rec {
    rev = "054fd625232f8041af2c88320633047eb0574dc0";
    name = "gocraft-web-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/gocraft/web";
    #doCheck = true;

    src = fetchFromGitHub {
      inherit rev;
      owner = "gocraft";
      repo = "web";
      sha256 = "1lfsjdm4hpgs11mvc5qminjyl1gsq9w0i0ddzggfmak75bwqgdsm";
    };
  };

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

  #gorp = buildGoPackage rec {
  #  rev = "5d19ebd22fdcf3b9fb1c8cbfcf1d4e78e102f1e0";
  #  name = "gorp-${stdenv.lib.strings.substring 0 7 rev}";
  #  goPackagePath = "github.com/go-gorp/gorp";

  #  src = fetchFromGitHub {
  #    inherit rev;
  #    owner = "go-gorp";
  #    repo = "gorp";
  #    sha256 = "0wnkpdb5d2msv6bgicsaiv88a3ll1db9xxms4nvhmrm26dil7zza";
  #  };
  #};


  pq = buildGoPackage rec {
    rev = "11fc39a580a008f1f39bb3d11d984fb34ed778d9";
    name = "pq-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/lib/pq";

    src = fetchFromGitHub {
      inherit rev;
      owner = "lib";
      repo = "pq";
      sha256 = "02484mvy0c8ddhhhdwsjwhvzybsvzr2dwid8bws8zkvd6jlh0xdv";
    };
  };


  inflection = buildGoPackage rec {
    rev = "3272df6c21d04180007eb3349844c89a3856bc25";
    name = "inflection-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/jinzhu/inflection";

    src = fetchFromGitHub {
      inherit rev;
      owner = "jinzhu";
      repo = "inflection";
      sha256 = "80997485f7dd5df5c9503f0cd3088768dcd180b88a15cff79a11995369b03d89";
    };
  };

  gorm = buildGoPackage rec {
    rev = "84c6b46011b5b146782affd77dcf5ff95e255c50";
    name = "gorm-${stdenv.lib.strings.substring 0 7 rev}";
    goPackagePath = "github.com/jinzhu/gorm";

    buildInputs = [inflection pq];

    src = fetchFromGitHub {
      inherit rev;
      owner = "jinzhu";
      repo = "gorm";
      sha256 = "9f59809b3b3c2692df8b6a73814074f54dab0e71fe9e8ccd9a78c9516dd8b7a4";
    };
  };

#  gorilla-context = buildGoPackage rec {
#    rev = "1c83b3eabd45b6d76072b66b746c20815fb2872d";
#    name = "gorilla-context-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "github.com/gorilla/context";
#
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "gorilla";
#      repo = "context";
#      sha256 = "8bcdac0d9573e8476e19fe356d0324c154792f182c994d2b09dbf1649786b6ae";
#    };
#  };
#
#  gorilla-securecookie = buildGoPackage rec {
#    rev = "e95799a481bbcc3d01c2ad5178524cb8bec9f370";
#    name = "gorilla-securecookie-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "github.com/gorilla/securecookie";
#
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "gorilla";
#      repo = "securecookie";
#      sha256 = "dec0805eabad2389ebd496759c6999a8f1efb74440ebf63eacb3ec309a690d67";
#    };
#  };
#
#  gorilla-sessions = buildGoPackage rec {
#    rev = "f7261893ca3ea922c30eabe742c036d2c1de6e0a";
#    name = "gorilla-sessions-${stdenv.lib.strings.substring 0 7 rev}";
#    goPackagePath = "github.com/gorilla/sessions";
#
#    buildInputs = [ gorilla-context gorilla-securecookie ];
#
#    src = fetchFromGitHub {
#      inherit rev;
#      owner = "gorilla";
#      repo = "sessions";
#      sha256 = "959a28e74b2b2e93f4d3e7ef74d252a0af25b354550412042dc2cf15b3fd7d05";
#    };
#  };

}