/** @jsxImportSource xeact */

import { g, r, t, u, x } from "xeact";

type InstanceButtonProps = {
  text: string;
  instance_id: string;
  action: string;
  message: string;
  confirm?: boolean;
};

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

function Page(instance_id: string) {
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
    </div>
  );
}

r(() => {
  let root = g("actions");
  let instance_id = g("instance_id").innerText;

  let app = Page(instance_id);
  x(root);
  root.appendChild(app);
});
