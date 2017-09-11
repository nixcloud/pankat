[
  rec {
    goPackagePath = "gopkg.in/gomail.v2";
    fetch = {
      type = "git";
      url = "http://github.com/go-gomail/gomail";
      rev = "5ceb8e65415e45e1262fb385212b8193b55c0f99";
      sha256 = "03qqcfi50lp54dvfk5fxxfnwsvas6n2029ikgp1i7nd59jv7i1bm";
    };
  }
  rec {
    goPackagePath = "github.com/dchest/captcha";
    fetch = {
      type = "git";
      url = "http://github.com/dchest/captcha";
      rev = "9e952142169c3cd6268c6482a3a61c121536aca2";
      sha256 = "061lad5ynxjxq4ym656dyla5bym7vg4cqvrrcq0yp1h3lfkjxvbz";
    };
  }
  rec {
    goPackagePath = "golang.org/x/crypto";
    fetch = {
      type = "git";
      url = "http://github.com/golang/crypto";
      rev = "7b85b097bf7527677d54d3220065e966a0e3b613";
      sha256 = "0k21nnf0nszgbvml74sn68wc6p77pxbnfpi04dgarg6byd5rvxii";
    };
  }
  rec {
    goPackagePath = "github.com/gocraft/web";
    fetch = {
      type = "git";
      url = "http://github.com/gocraft/web";
      rev = "054fd625232f8041af2c88320633047eb0574dc0";
      sha256 = "1lfsjdm4hpgs11mvc5qminjyl1gsq9w0i0ddzggfmak75bwqgdsm";
    };
  }
  rec {
    goPackagePath = "github.com/githubnemo/CompileDaemon";
    fetch = {
      type = "git";
      url = "http://github.com/githubnemo/CompileDaemon";
      rev = "051a9ad079bf636e3db7fab6cfab1c4629b22519";
      sha256 = "1w4srfbyddw977q3329ww6czwhsb7lbkvf4px9dknjavzvxskycg";
    };
  }
  rec {
    goPackagePath = "github.com/fatih/color";
    fetch = {
      type = "git";
      url = "http://github.com/fatih/color";
      rev = "f773d4c806cc8e4a5749d6a35e2a4bbcd71443d6";
      sha256 = "1bd69gm6nig0g8zcsav68xs0h8sfdjj87fdrly9gf2k2r366bp9b";
    };
  }
  rec {
    goPackagePath = "github.com/mattn/go-isatty";
    fetch = {
      type = "git";
      url = "http://github.com/mattn/go-isatty";
      rev = "ae0b1f8f8004be68d791a576e3d8e7648ab41449";
      sha256 = "0qrcsh7j9mxcaspw8lfxh9hhflz55vj4aq1xy00v78301czq6jlj";
    };
  }
  rec {
    goPackagePath = "github.com/shiena/ansicolor";
    fetch = {
      type = "git";
      url = "http://github.com/shiena/ansicolor";
      rev = "a5e2b567a4dd6cc74545b8a4f27c9d63b9e7735b";
      sha256 = "0gwplb1b4fvav1vjf4b2dypy5rcp2w41vrbxkd1dsmac870cy75p";
    };
  }
  rec {
    goPackagePath = "github.com/howeyc/fsnotify";
    fetch = {
      type = "git";
      url = "http://github.com/howeyc/fsnotify";
      rev = "4894fe7efedeeef21891033e1cce3b23b9af7ad2";
      sha256 = "09r3h200nbw8a4d3rn9wxxmgma2a8i6ssaplf3zbdc2ykizsq7mn";
    };
  }
  rec {
    goPackagePath = "github.com/lib/pq";
    fetch = {
      type = "git";
      url = "http://github.com/lib/pq";
      rev = "11fc39a580a008f1f39bb3d11d984fb34ed778d9";
      sha256 = "02484mvy0c8ddhhhdwsjwhvzybsvzr2dwid8bws8zkvd6jlh0xdv";
    };
  }
  rec {
    goPackagePath = "github.com/jinzhu/inflection";
    fetch = {
      type = "git";
      url = "http://github.com/jinzhu/inflection";
      rev = "3272df6c21d04180007eb3349844c89a3856bc25";
      sha256 = "80997485f7dd5df5c9503f0cd3088768dcd180b88a15cff79a11995369b03d89";
    };
  }
  rec {
    goPackagePath = "github.com/jinzhu/gorm";
    fetch = {
      type = "git";
      url = "http://github.com/jinzhu/gorm";
      rev = "84c6b46011b5b146782affd77dcf5ff95e255c50";
      sha256 = "9f59809b3b3c2692df8b6a73814074f54dab0e71fe9e8ccd9a78c9516dd8b7a4";
    };
  }

  {
    goPackagePath  = "github.com/amir/raidman";
    fetch = {
      type = "git";
      url = "https://github.com/amir/raidman";
      rev =  "91c20f3f475cab75bb40ad7951d9bbdde357ade7";
      sha256 = "0pkqy5hzjkk04wj1ljq8jsyla358ilxi4lkmvkk73b3dh2wcqvpp";
    };
  }
  {
    goPackagePath = "github.com/elazarl/go-bindata-assetfs";
    fetch = {
      type = "git";
      url = "https://github.com/elazarl/go-bindata-assetfs";
      rev = "57eb5e1fc594ad4b0b1dbea7b286d299e0cb43c2";
      sha256 = "1za29pa15y2xsa1lza97jlkax9qj93ks4a2j58xzmay6rczfkb9i";
    };
  }
  {
   goPackagePath =  "github.com/garyburd/redigo";
    fetch =  {
      type =  "git";
       url =  "https://github.com/garyburd/redigo";
       rev =  "8873b2f1995f59d4bcdd2b0dc9858e2cb9bf0c13";
       sha256 =  "1lzhb99pcwwf5ddcs0bw00fwf9m1d0k7b92fqz2a01jlij4pm5l2";
    };
  }
  {
   goPackagePath =  "github.com/go-sql-driver/mysql";
    fetch =  {
      type =  "git";
       url =  "https://github.com/go-sql-driver/mysql";
       rev =  "7ebe0a500653eeb1859664bed5e48dec1e164e73";
       sha256 =  "1gyan3lyn2j00di9haq7zm3zcwckn922iigx3fvml6s2bsp6ljas";
    };
  }
  {
   goPackagePath =  "github.com/golang/protobuf";
    fetch =  {
      type =  "git";
       url =  "https://github.com/golang/protobuf";
       rev =  "bf531ff1a004f24ee53329dfd5ce0b41bfdc17df";
       sha256 =  "10lnvmq28jp2wk1xc32mdk4745lal2bmdvbjirckb9wlv07zzzf0";
    };
  }
  {
   goPackagePath =  "github.com/jeffail/gabs";
    fetch =  {
      type =  "git";
       url =  "https://github.com/jeffail/gabs";
       rev =  "ee1575a53249b51d636e62464ca43a13030afdb5";
       sha256 =  "0svv57193n8m86r7v7n0y9lny0p6nzr7xvz98va87h00mg146351";
    };
  }
  {
   goPackagePath =  "github.com/jeffail/util";
    fetch =  {
      type =  "git";
       url =  "https://github.com/jeffail/util";
       rev =  "48ada8ff9fcae546b5986f066720daa9033ad523";
       sha256 =  "0k8zz7gdv4hb691fdyb5mhlixppcq8x4ny84fanflypnv258a3i0";
    };
  }
  #{
  # goPackagePath =  "github.com/lib/pq";
  #  fetch =  {
  #    type =  "git";
  #     url =  "https://github.com/lib/pq";
  #     rev =  "3cd0097429be7d611bb644ef85b42bfb102ceea4";
  #     sha256 =  "1q7qfzyfgjk6rvid548r43fi4jhvsh4dhfvfjbp2pz4xqsvpsm7a";
  #  };
  #}
  {
   goPackagePath =  "github.com/satori/go.uuid";
    fetch =  {
      type =  "git";
       url =  "https://github.com/satori/go.uuid";
       rev =  "f9ab0dce87d815821e221626b772e3475a0d2749";
       sha256 =  "0z18j6zxq9kw4lgcpmhh3k7jrb9gy1lx252xz5qhs4ywi9w77xwi";
    };
  }

  {
   goPackagePath =  "golang.org/x/net";
    fetch =  {
      type =  "git";
       url =  "https://go.googlesource.com/net";
       rev =  "07b51741c1d6423d4a6abab1c49940ec09cb1aaf";
       sha256 =  "12lvdj0k2gww4hw5f79qb9yswqpy4i3bgv1likmf3mllgdxfm20w";
    };
  }
]

