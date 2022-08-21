import { u } from "xeact";

export type Config = {   
    base_url: string,
    hosts: string[],
    bind_host: string,
    port: number,
    rpool_base: string,
    qemu_path: string,
};

export const getConfig = async (): Promise<Config> => {
    const resp = await fetch(u("/admin/api/config"));
    if (resp.status !== 200) {
        const body = await resp.text();
        throw new Error("wrong status code: " + resp.status + "\n\n" + body);
    }

    const result: Config = await resp.json();
    return result;
}

export type Distro = {
    name: string;
    downloadURL: string;
    sha256Sum: string;
    minSize: string;
    format: string;
};

export const getDistros = async (): Promise<Distro[]> => {
    const resp = await fetch(u("/api/v1/distros"));
    if (resp.status !== 200) {
        const body = await resp.text();
        throw new Error("wrong status code: " + resp.status + "\n\n" + body);
    }

    const result: Distro[] = await resp.json();

    return result;
};

export type NewInstance = {
    name?: string;
    memory_mb?: number;
    cpus?: number;
    host: string;
    disk_size_gb?: number;
    zvol_prefix?: string;
    distro: string;
    user_data?: string;
    join_tailnet: boolean;
};

export type Instance = {
    uuid: string;
    name: string;
    host: string;
    mac_address: string;
    memory: number;
    disk_size: number;
    zvol_name: string;
    status: string;
    distro: string;
    join_tailnet: boolean;
};

export const makeInstance = async (ni: NewInstance): Promise<Instance> => {
    const resp = await fetch(u("/api/v1/instances"), {
        method: "POST",
        body: JSON.stringify(ni),
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json",
        },
    });

    if (resp.status !== 200) {
        const body = await resp.text();
        throw new Error("wrong status code: " + resp.status + "\n\n" + body);
    }

    const instance: Instance = await resp.json();
    return instance;
};

export const deleteInstance = async (id: string): Promise<void> => {
    await fetch(u(`/api/v1/instances/${id}`), {
        method: "DELETE",
    });
}

const doThingToInstance = (action: string): (id: string) => Promise<void> => {
    return (async (id: string): Promise<void> => {       
        const resp = await fetch(u(`/api/v1/instances/${id}/${action}`), {
            method: "POST",
        });

        if (resp.status !== 200) {
            const body = await resp.text();
            throw new Error("wrong status code: " + resp.status + "\n\n" + body);
        }
    });
}

export const rebootInstance = doThingToInstance("reboot");
export const hardRebootInstance = doThingToInstance("hardreboot");
export const reinitInstance = doThingToInstance("reinit");
export const shutdownInstance = doThingToInstance("shutdown");
export const startInstance = doThingToInstance("start");
