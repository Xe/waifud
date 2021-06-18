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
        owner -> Text,
    }
}

joinable!(instances -> connections (owner));

allow_tables_to_appear_in_same_query!(
    connections,
    instances,
);
