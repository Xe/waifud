/** @jsxImportSource xeact */

import { g, t, u } from "xeact";
import {
  deleteInstance,
  hardRebootInstance,
  rebootInstance,
  reinitInstance,
  shutdownInstance,
  startInstance,
} from "./waifud/mod.ts";

type InstanceButtonProps = {
  text: string;
  instance_id: string;
  action: string;
  message: string;
  confirm?: boolean;
};

function DeleteInstanceButton(
  { text, instance_id, message, confirm = true }: InstanceButtonProps,
) {
  const onclick = async () => {
    if (confirm) {
      const response = prompt(
        "Type 'I don't care about the data' to continue.",
      );
      if (response !== "I don't care about the data") {
        g("messages").appendChild(t("Confirmation failed."));
        return;
      }
    }
    await deleteInstance(instance_id);
    g("messages").appendChild(t(message));
    alert(message);
    window.location.href = u("/admin/instances");
  };
  return (
    <div>
      <button onclick={() => onclick()}>{text}</button>
      <br />
    </div>
  );
}

function InstanceButton(
  { text, instance_id, action, message, confirm = false }: InstanceButtonProps,
) {
  const onclick = async () => {
    if (confirm) {
      const response = prompt(
        "Type 'I don't care about the data' to continue.",
      );
      if (response !== "I don't care about the data") {
        g("messages").appendChild(t("Confirmation failed."));
        return;
      }
    }
    switch (action) {
      case "reboot":
        await rebootInstance(instance_id);
        break;
      case "hardreboot":
        await hardRebootInstance(instance_id);
        break;
      case "reinit":
        await reinitInstance(instance_id);
        break;
      case "shutdown":
        await shutdownInstance(instance_id);
        break;
      case "start":
        await startInstance(instance_id);
        break;
    }
    g("messages").appendChild(t(message));
  };
  return (
    <div>
      <button onclick={() => onclick()}>{text}</button>
      <br />
    </div>
  );
}

export async function Page() {
  const instance_id = g("instance_id").innerText;
  return (
    <div>
      <InstanceButton
        text="Reboot"
        instance_id={instance_id}
        action="reboot"
        message="VM Rebooted."
      />
      <InstanceButton
        text="Hard Reboot"
        instance_id={instance_id}
        action="hardreboot"
        message="VM hard-rebooted."
      />
      <InstanceButton
        text="Recreate VM"
        instance_id={instance_id}
        action="reinit"
        message="Recreating VM from scratch."
        confirm={true}
      />
      <InstanceButton
        text="Shutdown"
        instance_id={instance_id}
        action="shutdown"
        message="VM shut down."
      />
      <InstanceButton
        text="Start"
        instance_id={instance_id}
        action="start"
        message="VM Started."
      />
      <DeleteInstanceButton
        text="Delete instance"
        instance_id={instance_id}
        action="delete"
        message="Instance deleted, redirecting you to instances page."
      />
      <div id="messages">
        <h3>Messages</h3>
      </div>
    </div>
  );
}
