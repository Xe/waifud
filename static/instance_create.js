import { u, g, x, r, h, t } from "./xeact.js";
import { div, br } from "./xeact-html.js";

const tr = (attrs = {}, children = []) => h("tr", attrs, children);
const td = (attrs = {}, children = []) => h("td", attrs, children);
const th = (attrs = {}, children = []) => h("th", attrs, children);

const getDistros = async () => {
    let resp = await fetch(u("/api/v1/distros"));
    if (resp.status !== 200) {
        let body = await resp.text();
        throw new Error("wrong status code: " + resp.status + "\n\n" + body);
    }

    resp = await resp.json();

    return resp;
}

const user_data_template = `#cloud-config
#vim:syntax=yaml

users:
  - name: xe
    groups: [ wheel ]
    sudo: [ "ALL=(ALL) NOPASSWD:ALL" ]
    shell: /bin/bash
`

r(async () => {
    // get a list of distros
    // form to build the create request
    // submit form
    // redirect to instance detail page

    /*
    pub struct NewInstance {
      pub name: Option<String>,
      pub memory_mb: Option<i32>,
      pub cpus: Option<i32>,
      pub host: String,
      pub disk_size_gb: Option<i32>,
      pub zvol_prefix: Option<String>,
      pub distro: String,
      pub user_data: Option<String>,
      pub join_tailnet: bool,
    }
    */

    let distros = await getDistros();

    let nameBox = h("input", {type: "text", placeholder: "crobat"});
    let memoryBox = h("input", {type: "text", placeholder: "512"});
    let cpuBox = h("input", {type: "text", placeholder: "2"});
    let hosts = [ "kos-mos", "logos", "ontos", "pneuma" ].map(host => h("option", {value: host}, t(host)));
    let host = h("select", {}, hosts);
    let disk_size_gb = h("input", {type: "text", placeholder: "25"});
    let zvol_prefix = h("input", {type: "text", value: "rpool/local/vms"});
    let distro_options = distros.map((d) => h("option", {value: d.name}, t(d.name)));
    let distro = h("select", {id: "selectBox"}, distro_options);
    distro.onchange = () => {
        let selectedDistro = null;
        distros.forEach(d => {
            if (d.name == distro.value) {
                selectedDistro = d;
            }
        });

        if (selectedDistro == null) {
            console.log("this shouldn't happen, selected distro doesn't exist in our list??");
            return;
        }

        let disk_size = parseInt(disk_size_gb.value, 10);
        if (disk_size < selectedDistro.minSize) {
            disk_size_gb.value = `${selectedDistro.minSize}`;
        }
    };

    let user_data = h("textarea", {rows: 10, cols: 40}, t(user_data_template));
    let join_tailnet = h("input", {type: "checkbox"});

    let submit = h("button", {}, t("Create that sucker"));
    submit.onclick = async () => {
        let req = {
            name: nameBox.value != "" ? nameBox.value : null,
            memory_mb: memoryBox.value != "" ? parseInt(memoryBox.value, 10) : null,
            cpus: cpuBox.value != "" ? parseInt(cpuBox.value, 10) : null,
            host: host.value,
            disk_size_gb: disk_size_gb.value != "" ? parseInt(disk_size_gb.value, 10) : null,
            zvol_prefix: zvol_prefix.value != "" ? zvol_prefix.value : null,
            distro: distro.value,
            user_data: user_data.value,
            join_tailnet: join_tailnet.checked,
        };
        console.log(req);

        let resp = await fetch(u("/api/v1/instances"), {
            method: "POST",
            body: JSON.stringify(req),
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }
        });

        if (resp.status !== 200) {
            let body = await resp.text();
            throw new Error("wrong status code: " + resp.status + "\n\n" + body);
        }

        resp = await resp.json();
        console.log(resp);
        // resp is waifud::models::Instance
        window.location.href = u(`/admin/instances/${resp.uuid}`);
    };

    let root = g("root");
    x(root);

    root.appendChild(div({}, [
        h("table", {}, [
            tr({}, [
                th({}, t("Name")),
                td({}, nameBox),
            ]),
            tr({}, [
                th({}, t("Memory (MB)")),
                td({}, memoryBox)
            ]),
            tr({}, [
                th({}, t("CPU cores")),
                td({}, cpuBox),
            ]),
            tr({}, [
                th({}, t("Host")),
                td({}, host),
            ]),
            tr({}, [
                th({}, t("Disk size (GB)")),
                td({}, disk_size_gb),
            ]),
            tr({}, [
                th({}, t("ZVol prefix")),
                td({}, zvol_prefix),
            ]),
            tr({}, [
                th({}, t("Distro")),
                td({}, distro),
            ]),
            tr({}, [
                th({}, t("Userdata")),
                td({}, user_data),
            ]),
            tr({}, [
                th({}, t("Join tailnet + SSH?")),
                td({}, join_tailnet),
            ]),
            tr({}, t("")),
            h("tr", {}, [
                h("td", {}, submit)
            ]),
        ]),
    ]));
})
