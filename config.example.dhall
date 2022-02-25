let Tailscale =
      { Type = { apiKey : Text, tailnet : Text }
      , default =
        { apiKey = env:TAILSCALE_API_KEY ? ""
        , tailnet = env:TAILSCALE_TAILNET ? "cetacean.org.github"
        }
      }

let Yubikey =
      { Type = { clientID : Text, key : Text }
      , default =
        { clientID = env:YUBIKEY_CLIENT_ID ? "", key = env:YUBIKEY_KEY ? "" }
      }

let Config =
      { Type =
          { baseURL : Text
          , hosts : List Text
          , bindHost : Text
          , port : Natural
          , tailscale : Tailscale.Type
          , yubikey : Yubikey.Type
          }
      , default =
        { baseURL = "http://192.168.122.1:23818"
        , hosts = [ "vmhost1", "vmhost2" ]
        , bindHost = "::"
        , port = 23818
        , tailscale = Tailscale::{=}
        , yubikey = Yubikey::{=}
        }
      }

let defaultPort = env:PORT ? 23818

in  Config::{ port = defaultPort }
