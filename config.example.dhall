let Config =
      { Type =
          { baseURL : Text, hosts : List Text, bindHost : Text, port : Natural }
      , default =
        { baseURL = "http://192.168.122.1:23818"
        , hosts = [ "vmhost1", "vmhost2" ]
        , bindHost = "::"
        , port = 23818
        }
      }

let defaultPort = env:PORT ? 23818

in  Config::{ port = defaultPort }
