import { u, g, x, r, h, t } from "./xeact.js";
import { div, br } from "./xeact-html.js";

r(() => {
    let root = g("actions");
    let instance_id = g("instance_id").innerText;

    x(root);

    let reboot = h("button", {}, t("Reboot"));
    reboot.onclick = async () => {
        await fetch(u(`/api/v1/instances/${instance_id}/reboot`), {
            method: "POST",
        });
        root.appendChild(t("rebooted."));
    };

    let hard_reboot = h("button", {}, t("Hard Reboot"));
    hard_reboot.onclick = async () => {
        await fetch(u(`/api/v1/instances/${instance_id}/hardreboot`), {
            method: "POST",
        });
        root.appendChild(t("forcibly rebooted."));
    };

    let reinit = h("button", {}, t("Recreate VM"));
    reinit.onclick = async () => {
        let response = prompt("Type 'I don't care about the data' to continue.");b
        if (response !== "I don't care about the data") {
            root.appendChild(t("Recreate VM call failed."));
            return;
        }
        await fetch(u(`/api/v1/instances/${instance_id}/reinit`), {
            method: "POST",
        });
    }

    let shutdown = h("button", {}, t("Shutdown"));
    shutdown.onclick = async () => {
        await fetch(u(`/api/v1/instances/${instance_id}/shutdown`), {
            method: "POST",
        });
        root.appendChild(t("VM shut down."));
        
    };

    let start = h("button", {}, t("Start"));
    start.onclick = async () => {
        await fetch(u(`/api/v1/instances/${instance_id}/start`), {
            method: "POST",
        });
        root.appendChild(t("VM started."));
        
    };

    root.appendChild(div({}, [
        reboot,
        br(),
        hard_reboot,
        br(),
        reinit,
        br(),
        shutdown,
        br(),
        start,
        br(),
    ]));
});

