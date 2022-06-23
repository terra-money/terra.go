package msg

func NewMsgStoreCode(sender AccAddress, wasmByteCode []byte, instantiatePermission *AccessConfig) *StoreCode {
	return &StoreCode{Sender: sender.String(), WASMByteCode: wasmByteCode, InstantiatePermission: instantiatePermission}
}

func NewMsgInstantiateContract(sender, admin AccAddress, codeId uint64, label string, msg []byte) *InstantiateContract {
	return &InstantiateContract{Sender: sender.String(), Admin: admin.String(), CodeID: codeId, Label: label, Msg: msg}
}

func NewMsgExecuteContract(sender, contract AccAddress, msg []byte, funds Coins) *ExecuteContract {
	return &ExecuteContract{Sender: sender.String(), Contract: contract.String(), Msg: msg, Funds: funds}
}

func NewMsgMigrateContract(sender, contract AccAddress, codeId uint64, msg []byte) *MigrateContract {
	return &MigrateContract{Sender: sender.String(), Contract: contract.String(), CodeID: codeId, Msg: msg}
}