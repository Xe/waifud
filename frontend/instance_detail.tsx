/** @jsxImportSource xeact */

import { g, r, t, u, x } from "xeact";

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
        g("actions").appendChild(t("Confirmation failed."));
        return;
      }
    }
    await fetch(u(`/api/v1/instances/${instance_id}`), {
      method: "DELETE",
    });
    g("actions").appendChild(t(message));
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
        g("actions").appendChild(t("Confirmation failed."));
        return;
      }
    }
    await fetch(u(`/api/v1/instances/${instance_id}/${action}`), {
      method: "POST",
    });
    g("actions").appendChild(t(message));
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
        action = "delete"
        message="Instance deleted, redirecting you to instances page."
      />
    </div>
  );
}

/* r(async () => {
 *   const root = g("app");
 *   const app = await Page();
 *   x(root);
 *   root.appendChild(app);
 * }); */
