let Config =
      { Type =
          { baseURL : Text, hosts : List Text, bindHost : Text, port : Natural }
      , default =
        { baseURL = "http://100.87.242.16:23818"
        , hosts = [ "logos", "ontos", "kos-mos", "pneuma" ]
        , bindHost = "::"
        , port = 23818
        }
      }

let defaultPort = env:PORT ? 23818

in  Config::{ port = defaultPort }
