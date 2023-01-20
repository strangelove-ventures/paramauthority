# paramauthority

This is mainly intended for ICS consumer chains that need authoritative governance.

Override modules for `params` and `upgrades` Cosmos SDK v0.45.x modules, with an added `authority` in the params. The `authority` is configurable for both modules in the genesis.json file at `app_state.params.params.authority` and `app_state.upgrade.params.authority`, respectively.

`params` has an additional message, `MsgUpdateParams`, which allows the authority to update params.

`upgrade` has additional messages `MsgSoftwareUpgrade` and `MsgCancelUpgrade` which allow the authority to schedule and cancel scheduled upgrades.
