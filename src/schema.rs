table! {
    connections (hostname) {
        hostname -> Text,
        connection_uri -> Text,
    }
}

table! {
    instances (id) {
        id -> Text,
        name -> Text,
        ram -> Integer,
        cores -> Integer,
        zvol -> Text,
        zvol_size -> Integer,
        use_sata -> Nullable<Bool>,
        mac_address -> Text,
        owner -> Text,
    }
}

table! {
    templates (uuid) {
        uuid -> Text,
        name -> Text,
        distro -> Text,
        version -> Text,
        download_url -> Text,
        sha256sum -> Text,
        local_url -> Text,
    }
}

joinable!(instances -> connections (owner));

allow_tables_to_appear_in_same_query!(
    connections,
    instances,
    templates,
);
