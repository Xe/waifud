let Tailscale =
      { Type = { apiKey : Text, tailnet : Text }
      , default =
        { apiKey = env:TAILSCALE_API_KEY ? ""
        , tailnet = env:TAILSCALE_TAILNET ? "cetacean.org.github"
        }
      }

let Config =
      { Type =
          { baseURL : Text
          , hosts : List Text
          , bindHost : Text
          , port : Natural
          , tailscale : Tailscale.Type
          }
      , default =
        { baseURL = "http://192.168.122.1:23818"
        , hosts = [ "vmhost1", "vmhost2" ]
        , bindHost = "::"
        , port = 23818
        , tailscale = Tailscale::{=}
        }
      }

let defaultPort = env:PORT ? 23818

in  Config::{ port = defaultPort }
